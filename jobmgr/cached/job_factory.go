package cached

import (
	"sync"
	"time"

	pbjob "github.com/uber/peloton/.gen/peloton/api/v0/job"
	"github.com/uber/peloton/.gen/peloton/api/v0/peloton"
	pbtask "github.com/uber/peloton/.gen/peloton/api/v0/task"
	"github.com/uber/peloton/storage"

	log "github.com/sirupsen/logrus"
	"github.com/uber-go/tally"
)

// JobFactory is the entrypoint object into the cache which stores job and tasks.
// This only runs in the job manager leader.
type JobFactory interface {
	// AddJob will create a Job if not present in cache,
	// else returns the current cached Job.
	AddJob(id *peloton.JobID) Job

	// ClearJob cleans up the job from the cache.
	ClearJob(jobID *peloton.JobID)

	// GetJob will return the current cached Job,
	// and nil if currently not in cache.
	GetJob(id *peloton.JobID) Job

	// GetAllJobs returns the list of all jobs in cache.
	GetAllJobs() map[string]Job

	// Start emitting metrics.
	Start()

	// Stop clears the current jobs and tasks in cache, stops metrics.
	Stop()
}

type jobFactory struct {
	sync.RWMutex //  Mutex to acquire before accessing any variables in the job factory object

	// map of active jobs (job identifier -> cache job object) in the system
	jobs        map[string]*job
	running     bool                          // whether job factory is running
	jobStore    storage.JobStore              // storage job store object
	taskStore   storage.TaskStore             // storage task store object
	updateStore storage.UpdateStore           // storage update store object
	volumeStore storage.PersistentVolumeStore // storage volume store object
	mtx         *Metrics                      // cache metrics
	// Tob/task listeners. This list is immutable after object is created.
	// So it can read without a lock.
	listeners []JobTaskListener
	// channel to indicate that the job factory needs to stop
	stopChan chan struct{}
}

// InitJobFactory initializes the job factory object.
func InitJobFactory(
	jobStore storage.JobStore,
	taskStore storage.TaskStore,
	updateStore storage.UpdateStore,
	volumeStore storage.PersistentVolumeStore,
	parentScope tally.Scope,
	listeners []JobTaskListener) JobFactory {
	return &jobFactory{
		jobs:        map[string]*job{},
		jobStore:    jobStore,
		taskStore:   taskStore,
		updateStore: updateStore,
		volumeStore: volumeStore,
		mtx:         NewMetrics(parentScope.SubScope("cache")),
		listeners:   listeners,
	}
}

func (f *jobFactory) AddJob(id *peloton.JobID) Job {
	if j := f.GetJob(id); j != nil {
		return j
	}

	f.Lock()
	defer f.Unlock()
	// check whether the job exists again, in case it
	// is created between RUnlock and Lock
	j, ok := f.jobs[id.GetValue()]
	if !ok {
		j = newJob(id, f)
		f.jobs[id.GetValue()] = j
	}

	return j
}

// ClearJob removes the job and all it tasks from inventory
func (f *jobFactory) ClearJob(id *peloton.JobID) {
	j := f.GetJob(id)
	if j == nil {
		return
	}

	f.Lock()
	defer f.Unlock()
	delete(f.jobs, j.ID().GetValue())
}

func (f *jobFactory) GetJob(id *peloton.JobID) Job {
	f.RLock()
	defer f.RUnlock()

	if j, ok := f.jobs[id.GetValue()]; ok {
		return j
	}

	return nil
}

func (f *jobFactory) GetAllJobs() map[string]Job {
	f.RLock()
	defer f.RUnlock()

	jobMap := make(map[string]Job)
	for k, v := range f.jobs {
		jobMap[k] = v
	}
	return jobMap
}

// Start the job factory, starts emitting metrics.
func (f *jobFactory) Start() {
	f.Lock()
	defer f.Unlock()

	if f.running {
		return
	}
	f.running = true

	f.stopChan = make(chan struct{})
	go f.runPublishMetrics(f.stopChan)
	log.Info("job factory started")
}

// Stop clears the current jobs and tasks in cache, stops emitting metrics.
func (f *jobFactory) Stop() {
	f.Lock()
	defer f.Unlock()

	// Do not do anything if not runnning
	if !f.running {
		log.Info("job factory stopped")
		return
	}

	f.running = false
	f.jobs = map[string]*job{}
	close(f.stopChan)
	log.Info("job factory stopped")
}

//TODO Refactor to remove the metrics loop into a separate component.
// JobFactory should only implement an interface like MetricsProvides
// to periodically publish metrics instead of having its own go routine.
// runPublishMetrics is the entrypoint to start and stop publishing cache metrics
func (f *jobFactory) runPublishMetrics(stopChan <-chan struct{}) {
	ticker := time.NewTicker(_defaultMetricsUpdateTick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			f.publishMetrics()
		case <-stopChan:
			return
		}
	}
}

// publishMetrics is the routine which publishes cache metrics to M3
func (f *jobFactory) publishMetrics() {
	// Initialise tasks count map for all possible pairs of (state, goal_state)
	tCount := map[pbtask.TaskState]map[pbtask.TaskState]float64{}
	for s := range pbtask.TaskState_name {
		tCount[pbtask.TaskState(s)] = map[pbtask.TaskState]float64{}
		for gs := range pbtask.TaskState_name {
			tCount[pbtask.TaskState(s)][pbtask.TaskState(gs)] = 0.0
		}
	}

	// Iterate through jobs, tasks and count
	jobsCopy := f.getJobs()
	for _, j := range jobsCopy {
		j.RLock()
		for _, t := range j.tasks {
			t.RLock()
			tCount[t.runtime.GetState()][t.runtime.GetGoalState()]++
			t.RUnlock()
		}
		j.RUnlock()
	}

	// Publish
	f.mtx.scope.Gauge("jobs_count").Update(float64(len(jobsCopy)))
	for s, sm := range tCount {
		for gs, tc := range sm {
			f.mtx.scope.Tagged(map[string]string{"state": s.String(), "goal_state": gs.String()}).Gauge("tasks_count").Update(tc)
		}
	}
}

// getJobs returns a copy of the internal map of job
func (f *jobFactory) getJobs() map[string]*job {
	result := make(map[string]*job)
	f.RLock()
	defer f.RUnlock()

	for k, v := range f.jobs {
		result[k] = &job{
			tasks: v.tasks,
		}
	}
	return result
}

func (f *jobFactory) notifyJobRuntimeChanged(
	jobID *peloton.JobID,
	jobType pbjob.JobType,
	runtime *pbjob.RuntimeInfo) {

	if runtime != nil {
		for _, l := range f.listeners {
			l.JobRuntimeChanged(jobID, jobType, runtime)
		}
		// TODO add metric for listener execution latency
	}
}

func (f *jobFactory) notifyTaskRuntimeChanged(
	jobID *peloton.JobID,
	instanceID uint32,
	jobType pbjob.JobType,
	runtime *pbtask.RuntimeInfo) {

	if runtime != nil {
		for _, l := range f.listeners {
			l.TaskRuntimeChanged(jobID, instanceID, jobType, runtime)
		}
		// TODO add metric for listener execution latency
	}
}
