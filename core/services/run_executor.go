package services

import (
	"fmt"
	"time"

	"nulink/core/adapters"
	"nulink/core/logger"
	"nulink/core/services/synchronization"
	"nulink/core/store"
	"nulink/core/store/models"
	"nulink/core/store/orm"

	"github.com/pkg/errors"
)

//go:generate mockery -name RunExecutor -output ../internal/mocks/ -case=underscore

// RunExecutor handles the actual running of the job tasks
type RunExecutor interface {
	Execute(*models.ID) error
}

type runExecutor struct {
	store       *store.Store
	statsPusher synchronization.StatsPusher
}

// NewRunExecutor initializes a RunExecutor.
func NewRunExecutor(store *store.Store, statsPusher synchronization.StatsPusher) RunExecutor {
	return &runExecutor{
		store:       store,
		statsPusher: statsPusher,
	}
}

// Execute performs the work associate with a job run
func (re *runExecutor) Execute(runID *models.ID) error {
	run, err := re.store.Unscoped().FindJobRun(runID)
	if err != nil {
		return errors.Wrapf(err, "error finding run %s", runID)
	}

	for taskIndex := range run.TaskRuns {
		taskRun := &run.TaskRuns[taskIndex]
		if !run.Status.Runnable() {
			logger.Debugw("Run execution blocked", run.ForLogger("task", taskRun.ID.String())...)
			break
		}

		if taskRun.Status.Completed() {
			continue
		}

		if meetsMinimumConfirmations(&run, taskRun, run.ObservedHeight) {
			start := time.Now()

			result := re.executeTask(&run, taskRun)

			taskRun.ApplyOutput(result)
			run.ApplyOutput(result)

			elapsed := time.Since(start).Seconds()

			logger.Debugw(fmt.Sprintf("Executed task %s", taskRun.TaskSpec.Type), run.ForLogger("task", taskRun.ID.String(), "elapsed", elapsed)...)

		} else {
			logger.Debugw("Pausing run pending confirmations",
				run.ForLogger("required_height", taskRun.MinimumConfirmations)...,
			)
			taskRun.Status = models.RunStatusPendingConfirmations
			run.Status = models.RunStatusPendingConfirmations

		}

		if err := re.store.ORM.SaveJobRun(&run); errors.Cause(err) == orm.OptimisticUpdateConflictError {
			logger.Debugw("Optimistic update conflict while updating run", run.ForLogger()...)
			return nil
		} else if err != nil {
			return err
		}

		re.statsPusher.PushNow()
	}

	if run.Status.Finished() {
		if run.Status.Errored() {
			logger.Warnw("Task failed", run.ForLogger()...)
		} else {
			logger.Debugw("All tasks complete for run", run.ForLogger()...)
		}
	}
	return nil
}

func (re *runExecutor) executeTask(run *models.JobRun, taskRun *models.TaskRun) models.RunOutput {
	taskCopy := taskRun.TaskSpec // deliberately copied to keep mutations local

	params, err := models.Merge(run.RunRequest.RequestParams, taskCopy.Params)
	if err != nil {
		return models.NewRunOutputError(err)
	}
	taskCopy.Params = params

	adapter, err := adapters.For(taskCopy, re.store.Config, re.store.ORM)
	if err != nil {
		return models.NewRunOutputError(err)
	}

	previousTaskRun := run.PreviousTaskRun()

	previousTaskInput := models.JSON{}
	if previousTaskRun != nil {
		previousTaskInput = previousTaskRun.Result.Data
	}

	data, err := models.Merge(run.RunRequest.RequestParams, previousTaskInput, taskRun.Result.Data)
	if err != nil {
		return models.NewRunOutputError(err)
	}

	input := *models.NewRunInput(run.ID, data, taskRun.Status)
	result := adapter.Perform(input, re.store)
	return result
}
