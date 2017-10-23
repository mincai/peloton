package tracked

import (
	"context"
	"fmt"

	"code.uber.internal/infra/peloton/.gen/peloton/api/peloton"

	jobmgr_task "code.uber.internal/infra/peloton/jobmgr/task"
	"code.uber.internal/infra/peloton/util"

	"github.com/pkg/errors"
)

// Initialize action does the following:
// 1. Sets the current state to TaskState_INITIALIZED
// 2. Sets the goal state depending on the JobType
// 3. Regenerates a new mesos task ID
func (t *task) initialize(ctx context.Context) error {
	m := t.job.m

	taskID := &peloton.TaskID{
		Value: fmt.Sprintf("%s-%d", t.job.ID().GetValue(), t.ID()),
	}

	taskInfo, err := m.taskStore.GetTaskByID(ctx, taskID.GetValue())
	if err != nil || taskInfo == nil {
		return fmt.Errorf("task info not found for %v", taskID.GetValue())
	}

	jobConfig, err := m.jobStore.GetJobConfig(ctx, t.job.id)
	if err != nil {
		return fmt.Errorf("job config not found for %v", t.job.id)
	}

	util.RegenerateMesosTaskID(taskInfo)

	// update task runtime
	taskInfo.Runtime.GoalState = jobmgr_task.GetDefaultGoalState(jobConfig.GetType())
	taskInfo.Runtime.StartTime = ""
	taskInfo.Runtime.CompletionTime = ""

	err = m.UpdateTask(ctx, t.job.ID(), taskInfo.GetInstanceId(), taskInfo)
	if err != nil {
		return errors.Wrapf(err, "unable to update task")
	}
	return nil
}