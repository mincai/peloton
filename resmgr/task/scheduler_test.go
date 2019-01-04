package task

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/uber-go/tally"

	"github.com/uber/peloton/.gen/peloton/api/v0/peloton"
	pb_respool "github.com/uber/peloton/.gen/peloton/api/v0/respool"
	"github.com/uber/peloton/.gen/peloton/api/v0/task"
	"github.com/uber/peloton/.gen/peloton/private/hostmgr/hostsvc/mocks"
	"github.com/uber/peloton/.gen/peloton/private/resmgr"
	"github.com/uber/peloton/.gen/peloton/private/resmgrsvc"

	"github.com/uber/peloton/common"
	"github.com/uber/peloton/common/eventstream"
	res_common "github.com/uber/peloton/resmgr/common"
	"github.com/uber/peloton/resmgr/queue"
	"github.com/uber/peloton/resmgr/respool"
	"github.com/uber/peloton/resmgr/scalar"
	store_mocks "github.com/uber/peloton/storage/mocks"
)

type SchedulerTestSuite struct {
	suite.Suite
	resTree            respool.Tree
	readyQueue         queue.MultiLevelList
	taskSched          *scheduler
	mockCtrl           *gomock.Controller
	rmTaskTracker      Tracker
	eventStreamHandler *eventstream.Handler
	mockHostmgr        *mocks.MockInternalHostServiceYARPCClient
}

func (suite *SchedulerTestSuite) SetupSuite() {
	suite.mockCtrl = gomock.NewController(suite.T())
	mockResPoolStore := store_mocks.NewMockResourcePoolStore(suite.mockCtrl)
	gomock.InOrder(
		mockResPoolStore.EXPECT().
			GetAllResourcePools(context.Background()).Return(suite.getResPools(), nil).AnyTimes(),
	)
	mockJobStore := store_mocks.NewMockJobStore(suite.mockCtrl)
	mockTaskStore := store_mocks.NewMockTaskStore(suite.mockCtrl)
	gomock.InOrder(
		mockJobStore.EXPECT().GetJobsByStates(context.Background(),
			gomock.Any()).Return(nil, nil).AnyTimes(),
	)
	respool.InitTree(tally.NoopScope, mockResPoolStore, mockJobStore,
		mockTaskStore, res_common.PreemptionConfig{Enabled: false})

	suite.resTree = respool.GetTree()
	suite.readyQueue = queue.NewMultiLevelList("ready-queue", maxReadyQueueSize)

	// Initializing the resmgr state machine
	InitTaskTracker(tally.NoopScope, &Config{})
	suite.rmTaskTracker = GetTracker()
	suite.eventStreamHandler = eventstream.NewEventStreamHandler(
		1000,
		[]string{
			common.PelotonJobManager,
			common.PelotonResourceManager,
		},
		nil,
		tally.Scope(tally.NoopScope))
	suite.taskSched = &scheduler{
		condition:        sync.NewCond(&sync.Mutex{}),
		resPoolTree:      suite.resTree,
		runningState:     res_common.RunningStateNotStarted,
		schedulingPeriod: time.Duration(1) * time.Second,
		stopChan:         make(chan struct{}, 1),
		queue:            suite.readyQueue,
		rmTaskTracker:    suite.rmTaskTracker,
		metrics:          NewMetrics(tally.NoopScope),
	}
	InitScheduler(
		tally.NoopScope,
		time.Duration(1)*time.Second,
		suite.rmTaskTracker)
}

func (suite *SchedulerTestSuite) TearDownSuite() {
	suite.mockCtrl.Finish()
	suite.rmTaskTracker.Clear()
}

func (suite *SchedulerTestSuite) SetupTest() {
	fmt.Println("setting up")
	suite.resTree.Start()
	suite.taskSched.runningState = 0
	suite.Nil(suite.taskSched.Start())

	suite.taskSched.runningState = 1
	suite.taskSched.Start()
	suite.AddTasks()
}

func (suite *SchedulerTestSuite) TearDownTest() {
	fmt.Println("tearing down")
	err := suite.resTree.Stop()
	suite.NoError(err)
	err = suite.taskSched.Stop()
	suite.NoError(err)
}

func TestTaskScheduler(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}

func (suite *SchedulerTestSuite) getResourceConfig() []*pb_respool.ResourceConfig {

	resConfigs := []*pb_respool.ResourceConfig{
		{
			Share:       1,
			Kind:        "cpu",
			Reservation: 100,
			Limit:       1000,
		},
		{
			Share:       1,
			Kind:        "memory",
			Reservation: 100,
			Limit:       1000,
		},
		{
			Share:       1,
			Kind:        "disk",
			Reservation: 100,
			Limit:       1000,
		},
		{
			Share:       1,
			Kind:        "gpu",
			Reservation: 2,
			Limit:       4,
		},
	}
	return resConfigs
}

// Returns resource pools
func (suite *SchedulerTestSuite) getResPools() map[string]*pb_respool.ResourcePoolConfig {
	rootID := peloton.ResourcePoolID{Value: "root"}
	policy := pb_respool.SchedulingPolicy_PriorityFIFO

	return map[string]*pb_respool.ResourcePoolConfig{
		"root": {
			Name:      "root",
			Parent:    nil,
			Resources: suite.getResourceConfig(),
			Policy:    policy,
		},
		"respool1": {
			Name:      "respool1",
			Parent:    &rootID,
			Resources: suite.getResourceConfig(),
			Policy:    policy,
		},
		"respool2": {
			Name:      "respool2",
			Parent:    &rootID,
			Resources: suite.getResourceConfig(),
			Policy:    policy,
		},
		"respool3": {
			Name:      "respool3",
			Parent:    &rootID,
			Resources: suite.getResourceConfig(),
			Policy:    policy,
		},
		"respool11": {
			Name:      "respool11",
			Parent:    &peloton.ResourcePoolID{Value: "respool1"},
			Resources: suite.getResourceConfig(),
			Policy:    policy,
		},
		"respool12": {
			Name:      "respool12",
			Parent:    &peloton.ResourcePoolID{Value: "respool1"},
			Resources: suite.getResourceConfig(),
			Policy:    policy,
		},
		"respool21": {
			Name:      "respool21",
			Parent:    &peloton.ResourcePoolID{Value: "respool2"},
			Resources: suite.getResourceConfig(),
			Policy:    policy,
		},
		"respool22": {
			Name:      "respool22",
			Parent:    &peloton.ResourcePoolID{Value: "respool2"},
			Resources: suite.getResourceConfig(),
			Policy:    policy,
		},
	}
}

func (suite *SchedulerTestSuite) AddTasks() {
	tasks := []*resmgr.Task{
		{
			Name:     "job1-1",
			Priority: 0,
			JobId:    &peloton.JobID{Value: "job1"},
			Id:       &peloton.TaskID{Value: "job1-1"},
			Resource: &task.ResourceConfig{
				CpuLimit:    1,
				DiskLimitMb: 10,
				GpuLimit:    0,
				MemLimitMb:  100,
			},
			Preemptible: true,
		},
		{
			Name:     "job1-1",
			Priority: 1,
			JobId:    &peloton.JobID{Value: "job1"},
			Id:       &peloton.TaskID{Value: "job1-2"},
			Resource: &task.ResourceConfig{
				CpuLimit:    1,
				DiskLimitMb: 10,
				GpuLimit:    0,
				MemLimitMb:  100,
			},
			Preemptible: true,
		},
		{
			Name:     "job2-1",
			Priority: 2,
			JobId:    &peloton.JobID{Value: "job2"},
			Id:       &peloton.TaskID{Value: "job2-1"},
			Resource: &task.ResourceConfig{
				CpuLimit:    1,
				DiskLimitMb: 10,
				GpuLimit:    0,
				MemLimitMb:  100,
			},
			Preemptible: true,
		},
		{
			Name:     "job2-2",
			Priority: 2,
			JobId:    &peloton.JobID{Value: "job2"},
			Id:       &peloton.TaskID{Value: "job2-2"},
			Resource: &task.ResourceConfig{
				CpuLimit:    1,
				DiskLimitMb: 10,
				GpuLimit:    0,
				MemLimitMb:  100,
			},
			Preemptible: true,
		},
	}
	resPool, err := suite.resTree.Get(&peloton.ResourcePoolID{
		Value: "respool11",
	})
	resPool.SetNonSlackEntitlement(suite.getEntitlement())

	for _, t := range tasks {
		suite.NoError(err)
		suite.addTasktotracker(t)
		rmTask := suite.rmTaskTracker.GetTask(t.Id)
		err = rmTask.TransitTo(task.TaskState_PENDING.String())
		suite.NoError(err)
		resPool.EnqueueGang(createGang(t))
	}
}

// createGang creates a gang from a single task
func createGang(task *resmgr.Task) *resmgrsvc.Gang {
	var gang resmgrsvc.Gang
	gang.Tasks = append(gang.Tasks, task)
	return &gang
}

func (suite *SchedulerTestSuite) validateReadyQueue(expectedTaskIDs []string) {
	for i := 0; i < len(expectedTaskIDs); i++ {
		item, err := suite.readyQueue.Pop(suite.readyQueue.Levels()[0])
		suite.NoError(err)
		gang := item.(*resmgrsvc.Gang)
		task := gang.Tasks[0]
		suite.Equal(expectedTaskIDs[i], task.Id.Value)
	}
}

func (suite *SchedulerTestSuite) TestMovingToReadyQueue() {
	time.Sleep(2000 * time.Millisecond)
	expectedTaskIDs := []string{
		"job2-1",
		"job2-2",
		"job1-2",
		"job1-1",
	}
	suite.validateReadyQueue(expectedTaskIDs)
}

func (suite *SchedulerTestSuite) TestMovingTasks() {
	suite.taskSched.scheduleTasks()
	expectedTaskIDs := []string{
		"job2-1",
		"job2-2",
		"job1-2",
		"job1-1",
	}
	suite.validateReadyQueue(expectedTaskIDs)
}

func (suite *SchedulerTestSuite) TestTaskStates() {
	suite.taskSched.scheduleTasks()
	for i := 0; i < 4; i++ {
		item, err := suite.readyQueue.Pop(suite.readyQueue.Levels()[0])
		suite.NoError(err)
		gang := item.(*resmgrsvc.Gang)
		task := gang.Tasks[0]
		rmTask := suite.rmTaskTracker.GetTask(task.Id)
		suite.Equal(rmTask.GetCurrentState().String(), "READY")
	}
}

func (suite *SchedulerTestSuite) getEntitlement() *scalar.Resources {
	return &scalar.Resources{
		CPU:    100,
		MEMORY: 1000,
		DISK:   100,
		GPU:    2,
	}
}

func setupScheduler(limit int64) *scheduler {
	return &scheduler{
		condition: sync.NewCond(&sync.Mutex{}),
		queue:     queue.NewMultiLevelList("ready-queue", limit),
		random:    rand.New(rand.NewSource(time.Now().UnixNano())),
		metrics:   NewMetrics(tally.NoopScope),
	}
}

func createGangs(count int, taskType resmgr.TaskType) []*resmgrsvc.Gang {
	result := make([]*resmgrsvc.Gang, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, &resmgrsvc.Gang{
			Tasks: []*resmgr.Task{
				{
					Type: taskType,
				},
			},
		})
	}
	return result
}

func TestScheduler_EnqueueGang_enqueues_until_the_queue_is_full(t *testing.T) {
	scheduler := setupScheduler(2)
	gangs := createGangs(3, resmgr.TaskType_BATCH)
	for i, gang := range gangs {
		err := scheduler.EnqueueGang(gang)
		if i < 2 {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, "list size limit reached")
		}
	}
}

func TestScheduler_DequeueGang_blocks_until_a_gang_is_added(t *testing.T) {
	scheduler := setupScheduler(3)
	gangs := createGangs(1, resmgr.TaskType_BATCH)
	start := time.Now()
	go func() {
		time.Sleep(20 * time.Millisecond)
		scheduler.EnqueueGang(gangs[0])
	}()
	gang, err := scheduler.DequeueGang(50*time.Millisecond, resmgr.TaskType_BATCH)
	end := time.Now()
	timeSpent := end.Sub(start)
	assert.True(t, timeSpent < 50*time.Millisecond)
	assert.True(t, timeSpent > 20*time.Millisecond)
	assert.NoError(t, err)
	assert.NotNil(t, gang)
}

func TestScheduler_DequeueGang_of_different_types_at_the_same_time(t *testing.T) {
	scheduler := setupScheduler(4)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(4)
	maxWaitTime := 10 * time.Millisecond
	dequeueType := resmgr.TaskType_STATELESS
	type state struct {
		dequeued bool
		used     time.Duration
	}
	states := make([]*state, 4)
	types := []resmgr.TaskType{
		resmgr.TaskType_BATCH, resmgr.TaskType_STATELESS, resmgr.TaskType_DAEMON, resmgr.TaskType_STATEFUL,
	}
	for i := range states {
		go func(index int) {
			start := time.Now()
			gang, err := scheduler.DequeueGang(maxWaitTime, types[index])
			end := time.Now()
			used := end.Sub(start)
			states[index] = &state{
				used: used,
			}
			if gang != nil && err == nil {
				states[index].dequeued = true
			}
			waitGroup.Done()
		}(i)
	}

	gangs := createGangs(1, dequeueType)
	for _, gang := range gangs {
		scheduler.EnqueueGang(gang)
	}
	waitGroup.Wait()
	for i, tt := range types {
		if tt == dequeueType {
			assert.True(t, states[i].dequeued)
			assert.True(t, states[i].used < maxWaitTime)
		} else {
			assert.False(t, states[i].dequeued)
			assert.True(t, states[i].used >= maxWaitTime)
		}
	}
}

func TestScheduler_EnqueueGang_and_DequeueGang_multiple_times(t *testing.T) {
	size := 100
	gangs := createGangs(size, resmgr.TaskType_BATCH)
	scheduler := setupScheduler(5)
	dequeued := uint32(0)
	join := sync.WaitGroup{}
	join.Add(2 * size)
	for _, gang := range gangs {
		go func() {
			for {
				result, err := scheduler.DequeueGang(time.Duration(50)*time.Millisecond, resmgr.TaskType_BATCH)
				if err == nil {
					assert.NotNil(t, result, "result should not be nil")
					if result != nil {
						atomic.AddUint32(&dequeued, 1)
					}
					break
				}
			}
			join.Done()
		}()
		go func(g *resmgrsvc.Gang) {
			for {
				err := scheduler.EnqueueGang(g)
				if err == nil {
					break
				}
			}
			join.Done()
		}(gang)
	}
	join.Wait()
	assert.Equal(t, size, int(dequeued))
}

func createMixedGangs(groupSizes map[resmgr.TaskType]int) []*resmgrsvc.Gang {
	total := 0
	for _, size := range groupSizes {
		total += size
	}
	result := make([]*resmgrsvc.Gang, 0, total)
	for taskType, size := range groupSizes {
		result = append(result, createGangs(size, taskType)...)
	}
	return result
}

func BenchmarkScheduler_EnqueueGang(b *testing.B) {
	share := b.N / 4
	gangs := createMixedGangs(map[resmgr.TaskType]int{
		resmgr.TaskType_BATCH:     share,
		resmgr.TaskType_STATELESS: share,
		resmgr.TaskType_DAEMON:    share,
		resmgr.TaskType_STATEFUL:  b.N - 3*share,
	})
	scheduler := setupScheduler(int64(b.N))

	b.ResetTimer()
	for _, gang := range gangs {
		scheduler.EnqueueGang(gang)
	}
}

func BenchmarkScheduler_DequeueGang(b *testing.B) {
	share := b.N / 4
	groupSizes := map[resmgr.TaskType]int{
		resmgr.TaskType_BATCH:     share,
		resmgr.TaskType_STATELESS: share,
		resmgr.TaskType_DAEMON:    share,
		resmgr.TaskType_STATEFUL:  b.N - 3*share,
	}
	gangs := createMixedGangs(groupSizes)
	scheduler := setupScheduler(int64(b.N))
	for _, gang := range gangs {
		scheduler.EnqueueGang(gang)
	}

	b.ResetTimer()
	for taskType, size := range groupSizes {
		for i := 0; i < size; i++ {
			scheduler.DequeueGang(10*time.Millisecond, taskType)
		}
	}
}

func (suite *SchedulerTestSuite) addTasktotracker(task *resmgr.Task) {
	resPool, _ := suite.resTree.Get(&peloton.ResourcePoolID{
		Value: "respool11",
	})
	suite.rmTaskTracker.AddTask(task, suite.eventStreamHandler, resPool, &Config{
		PolicyName: ExponentialBackOffPolicy,
	})
}

func (suite *SchedulerTestSuite) TestUntrackedTasks() {
	suite.rmTaskTracker.DeleteTask(&peloton.TaskID{Value: "job1-1"})
	suite.taskSched.scheduleTasks()
	expectedTaskIDs := []string{
		"job2-1",
		"job2-2",
		"job1-2",
	}
	suite.validateReadyQueue(expectedTaskIDs)
}

func (suite *SchedulerTestSuite) TestDeletePartialTasksTasks() {
	suite.rmTaskTracker.Clear()
	tasks := []*resmgr.Task{
		{
			Name:     "job1-1",
			Priority: 1,
			JobId:    &peloton.JobID{Value: "job1"},
			Id:       &peloton.TaskID{Value: "job1-1"},
			Resource: &task.ResourceConfig{
				CpuLimit:    1,
				DiskLimitMb: 10,
				GpuLimit:    0,
				MemLimitMb:  100,
			},
			Preemptible: true,
		},
		{
			Name:     "job1-2",
			Priority: 1,
			JobId:    &peloton.JobID{Value: "job1"},
			Id:       &peloton.TaskID{Value: "job1-2"},
			Resource: &task.ResourceConfig{
				CpuLimit:    1,
				DiskLimitMb: 10,
				GpuLimit:    0,
				MemLimitMb:  100,
			},
			Preemptible: true,
		},
	}

	resPool, err := suite.resTree.Get(&peloton.ResourcePoolID{
		Value: "respool11",
	})
	resPool.SetNonSlackEntitlement(suite.getEntitlement())

	for _, t := range tasks {
		suite.NoError(err)
		suite.addTasktotracker(t)
		rmTask := suite.rmTaskTracker.GetTask(t.Id)
		err = rmTask.TransitTo(task.TaskState_PENDING.String())
		suite.NoError(err)
	}
	gang := &resmgrsvc.Gang{
		Tasks: tasks,
	}
	resPool.EnqueueGang(gang)
	suite.rmTaskTracker.DeleteTask(&peloton.TaskID{Value: "job1-1"})

	suite.taskSched.scheduleTasks()

	// 1. Admission control passes for gang with 2 tasks.
	//    Thereby allocation is equal to resources to both tasks
	// 2. During transit tasks from pending -> ready state it fails
	//	  Thereby, it decreases the allocation for 1 task
	//	  and increases then demand
	suite.Equal(
		scalar.GetGangResources(gang).GetCPU()/2,
		resPool.GetTotalAllocatedResources().GetCPU())
	suite.Equal(
		scalar.GetGangResources(gang).GetMem()/2,
		resPool.GetTotalAllocatedResources().GetMem())
	suite.Equal(
		scalar.GetGangResources(gang).GetDisk()/2,
		resPool.GetTotalAllocatedResources().GetDisk())

	suite.Equal(
		scalar.GetGangResources(gang).GetCPU()/2,
		resPool.GetDemand().GetCPU())
	suite.Equal(
		scalar.GetGangResources(gang).GetMem()/2,
		resPool.GetDemand().GetMem())
	suite.Equal(
		scalar.GetGangResources(gang).GetDisk()/2,
		resPool.GetDemand().GetDisk())

	suite.rmTaskTracker.Clear()

	expectedTaskIDs := []string{
		"job1-2",
	}
	suite.validateReadyQueue(expectedTaskIDs)
}
