// Task state machine

package task

import (
	"context"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	uatomic "github.com/uber-go/atomic"
	"github.com/uber-go/tally"
	"go.uber.org/yarpc"

	mesos "code.uber.internal/infra/peloton/.gen/mesos/v1"
	sched "code.uber.internal/infra/peloton/.gen/mesos/v1/scheduler"
	pb_eventstream "code.uber.internal/infra/peloton/.gen/peloton/private/eventstream"
	"code.uber.internal/infra/peloton/.gen/peloton/private/resmgrsvc"

	"code.uber.internal/infra/peloton/common"
	"code.uber.internal/infra/peloton/common/cirbuf"
	"code.uber.internal/infra/peloton/common/eventstream"
	hostmgr_mesos "code.uber.internal/infra/peloton/hostmgr/mesos"
	"code.uber.internal/infra/peloton/yarpc/encoding/mpb"
)

const (
	_errorWaitInterval = 10 * time.Second
)

// StateManager is the interface for mesos task status updates stream.
type StateManager interface {
	// Update is the mesos callback to framework to notify mesos task status update change.
	Update(ctx context.Context, body *sched.Event) error

	// UpdateCounters manages counters of task status update & ack counts.
	UpdateCounters(_ *uatomic.Bool)

	// EventPurged is for implementing PurgedEventsProcessor interface.
	EventPurged(events []*cirbuf.CircularBufferItem)
}

type stateManager struct {
	schedulerclient mpb.SchedulerClient

	updateAckConcurrency int
	ackChannel           chan *mesos.TaskStatus // Buffers the mesos task status updates to be acknowledged

	eventStreamHandler *eventstream.Handler
	metrics            *Metrics
}

// eventForwarder is the struct to forward status update events to
// resource manager. It implements eventstream.EventHandler and it
// forwards the events to remote in the OnEvents function.
type eventForwarder struct {
	// client to send NotifyTaskUpdatesRequest
	client resmgrsvc.ResourceManagerServiceYARPCClient
	// Tracking the progress returned from remote side
	progress *uint64
}

// initResMgrEventForwarder, creates an event stream client to push
// mesos task status update events to Resource Manager from Host Manager.
func initResMgrEventForwarder(
	eventStreamHandler *eventstream.Handler,
	client resmgrsvc.ResourceManagerServiceYARPCClient,
	scope tally.Scope) {
	var progress uint64
	eventForwarder := &eventForwarder{
		client:   client,
		progress: &progress,
	}
	eventstream.NewLocalEventStreamClient(
		common.PelotonResourceManager,
		eventStreamHandler,
		eventForwarder,
		scope,
	)
}

// initEventStreamHandler initializes two event streams for communicating
// task status updates with Job Manager & Resource Manager.
// Job Manager: pulls task status update events from event stream.
// Resource Manager: Host Manager call event stream client
// to push task status update events.
func initEventStreamHandler(
	d *yarpc.Dispatcher,
	purgedEventProcessor eventstream.PurgedEventsProcessor,
	bufferSize int,
	scope tally.Scope) *eventstream.Handler {
	eventStreamHandler := eventstream.NewEventStreamHandler(
		bufferSize,
		[]string{common.PelotonJobManager, common.PelotonResourceManager},
		purgedEventProcessor,
		scope,
	)

	d.Register(pb_eventstream.BuildEventStreamServiceYARPCProcedures(eventStreamHandler))

	return eventStreamHandler
}

// NewStateManager init the task state manager by setting up input stream
// to receive mesos task status update, and outgoing event stream
// for Job Manager & Resource Manager for consumption of these task status updates.
func NewStateManager(
	d *yarpc.Dispatcher,
	schedulerClient mpb.SchedulerClient,
	updateBufferSize int,
	updateAckConcurrency int,
	resmgrClient resmgrsvc.ResourceManagerServiceYARPCClient,
	parentScope tally.Scope) StateManager {

	stateManagerScope := parentScope.SubScope("taskStateManager")
	handler := &stateManager{
		schedulerclient:      schedulerClient,
		updateAckConcurrency: updateAckConcurrency,
		ackChannel:           make(chan *mesos.TaskStatus, updateBufferSize),
		metrics:              NewMetrics(stateManagerScope),
	}
	mpb.Register(
		d,
		hostmgr_mesos.ServiceName,
		mpb.Procedure(sched.Event_UPDATE.String(), handler.Update))
	handler.startAsyncProcessTaskUpdates()
	handler.eventStreamHandler = initEventStreamHandler(
		d,
		handler,
		updateBufferSize,
		stateManagerScope.SubScope("EventStreamHandler"))
	initResMgrEventForwarder(
		handler.eventStreamHandler,
		resmgrClient,
		stateManagerScope.SubScope("ResourceManagerClient"))
	NewMetrics(parentScope)
	return handler
}

// GetEventProgress returns the event forward progress
func (f *eventForwarder) GetEventProgress() uint64 {
	return atomic.LoadUint64(f.progress)
}

// OnEvent callback
func (f *eventForwarder) OnEvent(event *pb_eventstream.Event) {
	//Not implemented
}

// OnEvents callback. In this callback, a batch of events are forwarded to
// resource manager.
func (f *eventForwarder) OnEvents(events []*pb_eventstream.Event) {
	if len(events) > 0 {
		// Forward events
		request := &resmgrsvc.NotifyTaskUpdatesRequest{
			Events: events,
		}
		var response *resmgrsvc.NotifyTaskUpdatesResponse
		for {
			var err error
			response, err = f.notifyResourceManager(request)
			if err == nil {
				break
			} else {
				log.WithError(err).WithField("progress", events[0].Offset).
					Error("Failed to call ResourceManager.NotifyTaskUpdate")
				time.Sleep(_errorWaitInterval)
			}
		}
		if response.PurgeOffset > 0 {
			atomic.StoreUint64(f.progress, response.PurgeOffset)
		}
		if response.Error != nil {
			log.WithField("notify_task_updates_response_error",
				response.Error).Error("NotifyTaskUpdatesRequest failed")
		}
	}
}

func (f *eventForwarder) notifyResourceManager(
	request *resmgrsvc.NotifyTaskUpdatesRequest) (
	*resmgrsvc.NotifyTaskUpdatesResponse, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	return f.client.NotifyTaskUpdates(ctx, request)
}

// Update is the Mesos callback on mesos state updates
func (m *stateManager) Update(ctx context.Context, body *sched.Event) error {
	var err error
	taskUpdate := body.GetUpdate()
	m.metrics.taskUpdateCounter.Inc(1)
	taskStateCounter := m.metrics.scope.Counter(
		"task_state_" + taskUpdate.GetStatus().GetState().String())
	taskStateCounter.Inc(1)

	event := &pb_eventstream.Event{
		MesosTaskStatus: taskUpdate.GetStatus(),
		Type:            pb_eventstream.Event_MESOS_TASK_STATUS,
	}
	err = m.eventStreamHandler.AddEvent(event)
	if err != nil {
		log.WithError(err).
			WithField("status_update", taskUpdate.GetStatus()).
			Error("Cannot add status update")
	}

	// If buffer is full, AddStatusUpdate would fail and peloton would not
	// ack the status update and mesos master would resend the status update.
	// Return nil otherwise the framework would disconnect with the mesos master
	return nil
}

// UpdateCounters tracks the count for task status update & ack count.
func (m *stateManager) UpdateCounters(_ *uatomic.Bool) {
	m.metrics.taskAckChannelSize.Update(float64(len(m.ackChannel)))
}

// startAsyncProcessTaskUpdates concurrently process task status update events
// ready to ACK iff uuid is not nil.
func (m *stateManager) startAsyncProcessTaskUpdates() {
	for i := 0; i < m.updateAckConcurrency; i++ {
		go func() {
			for taskStatus := range m.ackChannel {
				if len(taskStatus.GetUuid()) != 0 {
					// TODO (varung): Add retry for acknowledging status update
					err := m.acknowledgeTaskUpdate(
						context.Background(),
						taskStatus)
					if err != nil {
						log.WithField("task_status", *taskStatus).
							WithError(err).
							Error("Failed to acknowledgeTaskUpdate")
					}
				}
			}
		}()
	}
}

// acknowledgeTaskUpdate, ACK task status update events
// thru POST scheduler client call to Mesos Master.
func (m *stateManager) acknowledgeTaskUpdate(
	ctx context.Context,
	taskStatus *mesos.TaskStatus) error {
	m.metrics.taskUpdateAck.Inc(1)
	callType := sched.Call_ACKNOWLEDGE
	msg := &sched.Call{
		FrameworkId: hostmgr_mesos.GetSchedulerDriver().GetFrameworkID(ctx),
		Type:        &callType,
		Acknowledge: &sched.Call_Acknowledge{
			AgentId: taskStatus.AgentId,
			TaskId:  taskStatus.TaskId,
			Uuid:    taskStatus.Uuid,
		},
	}
	msid := hostmgr_mesos.GetSchedulerDriver().GetMesosStreamID(ctx)
	err := m.schedulerclient.Call(msid, msg)
	if err != nil {
		return err
	}
	log.WithField("task_status", *taskStatus).
		Debug("Acked task update")
	return nil
}

// EventPurged is for implementing PurgedEventsProcessor interface.
func (m *stateManager) EventPurged(events []*cirbuf.CircularBufferItem) {
	for _, e := range events {
		event, ok := e.Value.(*pb_eventstream.Event)
		if ok {
			if len(event.GetMesosTaskStatus().GetUuid()) > 0 {
				m.ackChannel <- event.GetMesosTaskStatus()
			}
		}
	}
}
