package goalstate

import (
	"context"
	"time"

	pbjob "code.uber.internal/infra/peloton/.gen/peloton/api/v0/job"
	pbupdate "code.uber.internal/infra/peloton/.gen/peloton/api/v0/update"

	"code.uber.internal/infra/peloton/common/goalstate"
	"code.uber.internal/infra/peloton/jobmgr/cached"
	"code.uber.internal/infra/peloton/util"

	"go.uber.org/yarpc/yarpcerrors"
)

// UpdateAbortIfNeeded checks if the update identifier in the goal
// state engine is the same as the one in the job runtime updater (tracking
// the current job update). If not, then it aborts the update in the goal
// state engine and enqueue the current update.
func UpdateAbortIfNeeded(ctx context.Context, entity goalstate.Entity) error {
	updateEnt := entity.(*updateEntity)
	goalStateDriver := updateEnt.driver
	jobID := updateEnt.jobID

	cachedJob := goalStateDriver.jobFactory.AddJob(jobID)
	runtime, err := cachedJob.GetRuntime(ctx)
	if err != nil {
		return err
	}

	if runtime.GetUpdateID().GetValue() == updateEnt.id.GetValue() {
		// update not been aborted, keep going
		return nil
	}

	cachedUpdate := goalStateDriver.updateFactory.GetUpdate(updateEnt.id)
	if cachedUpdate == nil {
		// no update in cache, recovery will be run anyways
		return nil
	}

	if err := cached.AbortJobUpdate(
		ctx,
		updateEnt.id,
		goalStateDriver.updateStore,
		goalStateDriver.updateFactory,
	); err != nil {
		return err
	}

	// return an error to ensure other update actions are not run and to
	// enqueue the same update back to the queue again for untracking
	return yarpcerrors.AbortedErrorf("update aborted")
}

// UpdateReload reloads the update from the DB.
func UpdateReload(ctx context.Context, entity goalstate.Entity) error {
	updateEnt := entity.(*updateEntity)
	goalStateDriver := updateEnt.driver
	cachedUpdate := goalStateDriver.updateFactory.AddUpdate(updateEnt.id)
	goalStateDriver.mtx.updateMetrics.UpdateReload.Inc(1)
	if err := cachedUpdate.Recover(ctx); err != nil {
		if !yarpcerrors.IsNotFound(err) {
			return err
		}
		// update not found in DB, just clean up from cache and goal state
		return UpdateUntrack(ctx, entity)
	}
	goalStateDriver.EnqueueUpdate(updateEnt.jobID, updateEnt.id, time.Now())
	return nil
}

// UpdateComplete indicates that all instances have been updated,
// and the update state should be marked complete.
func UpdateComplete(ctx context.Context, entity goalstate.Entity) error {
	updateEnt := entity.(*updateEntity)
	goalStateDriver := updateEnt.driver
	cachedUpdate := goalStateDriver.updateFactory.GetUpdate(updateEnt.id)
	if cachedUpdate == nil {
		goalStateDriver.mtx.updateMetrics.UpdateCompleteFail.Inc(1)
		return nil
	}

	completeState := pbupdate.State_SUCCEEDED
	if cachedUpdate.GetState().State == pbupdate.State_ROLLING_BACKWARD {
		completeState = pbupdate.State_ROLLED_BACK
	}

	if err := cachedUpdate.WriteProgress(
		ctx,
		completeState,
		cachedUpdate.GetInstancesDone(),
		cachedUpdate.GetInstancesFailed(),
		[]uint32{},
	); err != nil {
		goalStateDriver.mtx.updateMetrics.UpdateCompleteFail.Inc(1)
		return err
	}

	// enqueue to the goal state engine to untrack the update
	goalStateDriver.EnqueueUpdate(updateEnt.jobID, updateEnt.id, time.Now())
	goalStateDriver.mtx.updateMetrics.UpdateComplete.Inc(1)
	return nil
}

// UpdateUntrack deletes the update from the cache and the goal state engine.
func UpdateUntrack(ctx context.Context, entity goalstate.Entity) error {
	var runtime *pbjob.RuntimeInfo
	var err error

	updateEnt := entity.(*updateEntity)
	goalStateDriver := updateEnt.driver
	jobID := updateEnt.jobID
	cachedJob := goalStateDriver.jobFactory.AddJob(jobID)
	runtime, err = cachedJob.GetRuntime(ctx)
	if err != nil {
		return err
	}

	// clean up the update from cache and goal state
	goalStateDriver.DeleteUpdate(jobID, updateEnt.id)
	goalStateDriver.updateFactory.ClearUpdate(updateEnt.id)
	goalStateDriver.mtx.updateMetrics.UpdateUntrack.Inc(1)

	// check if we have another job update to run
	if len(runtime.GetUpdateID().GetValue()) > 0 &&
		runtime.GetUpdateID().GetValue() != updateEnt.id.GetValue() {
		goalStateDriver.EnqueueUpdate(jobID, runtime.GetUpdateID(), time.Now())
		return nil
	}

	// update can be applied to a terminated job,
	// need to remove job from cache upon completion
	if util.IsPelotonJobStateTerminal(runtime.GetState()) {
		goalStateDriver.EnqueueJob(jobID, time.Now())
	}

	// No more job update to run, so use the time to clean up any old
	// updates if they have not reached a terminal state yet
	updates, err := goalStateDriver.updateStore.GetUpdatesForJob(ctx, jobID)
	if err != nil {
		return nil
	}

	for _, prevUpdateID := range updates {
		updateModel, err :=
			goalStateDriver.updateStore.GetUpdateProgress(ctx, prevUpdateID)
		if err != nil {
			continue
		}
		if !cached.IsUpdateStateTerminal(updateModel.GetState()) {
			// just enqueue one and let it untrack first
			goalStateDriver.EnqueueUpdate(jobID, prevUpdateID, time.Now())
			return nil
		}
	}
	return nil
}

// UpdateWriteProgress write the current progress of update
func UpdateWriteProgress(ctx context.Context, entity goalstate.Entity) error {
	updateEnt := entity.(*updateEntity)
	goalStateDriver := updateEnt.driver
	cachedUpdate := goalStateDriver.updateFactory.GetUpdate(updateEnt.id)
	if cachedUpdate == nil {
		goalStateDriver.mtx.updateMetrics.UpdateWriteProgressFail.Inc(1)
		return nil
	}

	// all the instances being updated are finished, nothing new to update
	if len(cachedUpdate.GetInstancesCurrent()) == 0 {
		goalStateDriver.mtx.updateMetrics.UpdateWriteProgress.Inc(1)
		return nil
	}

	cachedJob := goalStateDriver.jobFactory.GetJob(cachedUpdate.JobID())
	if cachedJob == nil {
		goalStateDriver.mtx.updateMetrics.UpdateWriteProgressFail.Inc(1)
		return nil
	}

	instancesCurrent, instancesDone, instancesFailed, err := cached.GetUpdateProgress(
		ctx,
		cachedJob,
		cachedUpdate,
		cachedUpdate.GetGoalState().JobVersion,
		cachedUpdate.GetInstancesCurrent(),
	)
	if err != nil {
		goalStateDriver.mtx.updateMetrics.UpdateWriteProgressFail.Inc(1)
		return err
	}

	err = cachedUpdate.WriteProgress(
		ctx,
		cachedUpdate.GetState().State,
		append(cachedUpdate.GetInstancesDone(), instancesDone...),
		append(cachedUpdate.GetInstancesFailed(), instancesFailed...),
		instancesCurrent,
	)
	if err != nil {
		goalStateDriver.mtx.updateMetrics.UpdateWriteProgressFail.Inc(1)
		return err
	}

	goalStateDriver.mtx.updateMetrics.UpdateWriteProgress.Inc(1)
	return nil

}
