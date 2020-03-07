//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/internal/master/envs"
	"github.com/lastbackend/lastbackend/internal/master/state/cluster"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/util/generator"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"

	"github.com/lastbackend/lastbackend/internal/pkg/types"
)

func testTaskObserver(t *testing.T, name, werr string, wjs *JobState, js *JobState, task *types.Task) {
	var (
		ctx = context.Background()
		stg = envs.Get().GetStorage()
		err error
	)

	err = stg.Del(ctx, stg.Collection().Task(), "")
	if !assert.NoError(t, err) {
		return
	}

	err = stg.Del(ctx, stg.Collection().Pod(), "")
	if !assert.NoError(t, err) {
		return
	}

	t.Run(name, func(t *testing.T) {

		err := taskObserve(js, task)
		if werr != types.EmptyString {

			if assert.NoError(t, err, "error should be presented") {
				return
			}

			if !assert.Equal(t, werr, err.Error(), "err message different") {
				return
			}

			return
		}

		if wjs.job == nil {
			if !assert.Nil(t, js.job, "job should be nil") {
				return
			}
		}

		if err := compareJobStateProperties(wjs, js); assert.NoError(t, err) {
			return
		}
	})
}

func TestHandleTaskStateCreated(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
			task     *types.Task
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle without tasks should set task in queued state"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateCreated, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateQueued

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "failed state handle without tasks should set task in error state"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateCreated, types.EmptyString)
		volume := &types.SpecTemplateVolume{
			Volume: types.SpecTemplateVolumeClaim{Name: "demo"},
		}
		task.Spec.Template.Volumes = append(task.Spec.Template.Volumes, volume)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateError
		wt.Status.Error = true
		wt.Status.Message = fmt.Sprintf("%s: volume not found", strings.Title(volume.Volume.Name))

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt

		return s
	}())

	for _, tt := range tests {
		testTaskObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.task)
	}

}

func TestHandleTaskStateQueued(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
			task     *types.Task
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle without tasks should set task in provision state"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateQueued, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateProvision

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateRunning
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.queue[wt.SelfLink().String()] = wt

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with tasks in queue should be two task in provision state"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task1 := getTaskAsset(job, types.StateQueued, types.EmptyString)
		task2 := getTaskAsset(job, types.StateQueued, types.EmptyString)

		s.args.task = task2
		s.args.jobState = js
		s.args.jobState.task.list[task1.SelfLink().String()] = task1
		s.args.jobState.task.queue[task1.SelfLink().String()] = task1
		s.args.jobState.task.list[task2.SelfLink().String()] = task2

		wt1 := getTaskCopy(task1)
		wt1.Status.State = types.StateProvision
		wt2 := getTaskCopy(task2)
		wt2.Status.State = types.StateQueued

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateRunning
		s.want.jobState.task.list[wt1.SelfLink().String()] = wt1
		s.want.jobState.task.list[wt2.SelfLink().String()] = wt2
		s.want.jobState.task.queue[wt1.SelfLink().String()] = wt1
		s.want.jobState.task.queue[wt2.SelfLink().String()] = wt2

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with tasks in active should be two task in provisioned state because limits 2"}

		job := getJobAsset(types.StateRunning, types.EmptyString)
		job.Spec.Concurrency.Limit = 2
		js := getJobStateAsset(job)
		task1 := getTaskAsset(job, types.StateProvision, types.EmptyString)
		task2 := getTaskAsset(job, types.StateQueued, types.EmptyString)

		s.args.task = task2
		s.args.jobState = js
		s.args.jobState.task.list[task1.SelfLink().String()] = task1
		s.args.jobState.task.active[task1.SelfLink().String()] = task1
		s.args.jobState.task.list[task2.SelfLink().String()] = task2

		wt1 := getTaskCopy(task1)
		wt1.Status.State = types.StateProvision
		wt2 := getTaskCopy(task2)
		wt2.Status.State = types.StateProvision

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateRunning
		s.want.jobState.task.list[wt1.SelfLink().String()] = wt1
		s.want.jobState.task.list[wt2.SelfLink().String()] = wt2
		s.want.jobState.task.active[wt1.SelfLink().String()] = wt1
		s.want.jobState.task.queue[wt2.SelfLink().String()] = wt2

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with tasks in active should be task in queue state because limits 1"}

		job := getJobAsset(types.StateRunning, types.EmptyString)
		js := getJobStateAsset(job)
		task1 := getTaskAsset(job, types.StateProvision, types.EmptyString)
		task2 := getTaskAsset(job, types.StateQueued, types.EmptyString)

		s.args.task = task2
		s.args.jobState = js
		s.args.jobState.task.list[task1.SelfLink().String()] = task1
		s.args.jobState.task.active[task1.SelfLink().String()] = task1
		s.args.jobState.task.list[task2.SelfLink().String()] = task2

		wt1 := getTaskCopy(task1)
		wt1.Status.State = types.StateProvision
		wt2 := getTaskCopy(task2)
		wt2.Status.State = types.StateQueued

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateRunning
		s.want.jobState.task.list[wt1.SelfLink().String()] = wt1
		s.want.jobState.task.list[wt2.SelfLink().String()] = wt2
		s.want.jobState.task.active[wt1.SelfLink().String()] = wt1
		s.want.jobState.task.queue[wt2.SelfLink().String()] = wt2

		return s
	}())

	for _, tt := range tests {
		testTaskObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.task)
	}

}

func TestHandleTaskStateProvision(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
			job      *types.Job
			task     *types.Task
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle provision task without pod"}

		job := getJobAsset(types.StateRunning, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.job = job
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.queue[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateProvision
		pod, _ := podCreate(task)

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateRunning
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.active[wt.SelfLink().String()] = wt
		s.want.jobState.pod.list[wt.SelfLink().String()] = pod
		delete(s.want.jobState.task.queue, task.SelfLink().String())

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle have not task in queue"}

		job := getJobAsset(types.StateRunning, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateQueued

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateRunning
		s.want.jobState.task.list[wt.SelfLink().String()] = wt

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle provision task with pod"}

		job := getJobAsset(types.StateRunning, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)
		pod, _ := podCreate(task)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.queue[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Error = true
		wt.Status.Message = errors.New(fmt.Sprintf("pod %s can not be manage: node not attached", pod.Meta.SelfLink.String())).Error()

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateRunning
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.finished = append(s.want.jobState.task.finished, wt)
		s.want.jobState.pod.list = make(map[string]*types.Pod, 0)
		delete(s.want.jobState.task.queue, task.SelfLink().String())

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle provision task not queue"}

		job := getJobAsset(types.StateRunning, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateQueued

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateRunning
		s.want.jobState.task.list[wt.SelfLink().String()] = wt

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle provision task and pod in status destroy"}

		job := getJobAsset(types.StateRunning, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)
		pod, _ := podCreate(task)
		pod.Status.State = types.StateDestroy

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.queue[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateProvision

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateRunning
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.queue[wt.SelfLink().String()] = wt

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle provision task with pod and finished more 5"}

		job := getJobAsset(types.StateRunning, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)
		task1 := getTaskAsset(job, types.StateExited, types.EmptyString)
		task2 := getTaskAsset(job, types.StateExited, types.EmptyString)
		task3 := getTaskAsset(job, types.StateExited, types.EmptyString)
		task4 := getTaskAsset(job, types.StateExited, types.EmptyString)
		task5 := getTaskAsset(job, types.StateExited, types.EmptyString)
		pod, _ := podCreate(task)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.queue[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod
		s.args.jobState.task.finished = append(s.args.jobState.task.finished, task1)
		s.args.jobState.task.finished = append(s.args.jobState.task.finished, task2)
		s.args.jobState.task.finished = append(s.args.jobState.task.finished, task3)
		s.args.jobState.task.finished = append(s.args.jobState.task.finished, task4)
		s.args.jobState.task.finished = append(s.args.jobState.task.finished, task5)

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Error = true
		wt.Status.Message = errors.New(fmt.Sprintf("pod %s can not be manage: node not attached", pod.Meta.SelfLink.String())).Error()

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateRunning
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.finished = append(s.want.jobState.task.finished, wt)
		s.want.jobState.pod.list = make(map[string]*types.Pod, 0)
		s.want.jobState.task.finished = s.want.jobState.task.finished[1:]
		delete(s.want.jobState.task.queue, task.SelfLink().String())

		return s
	}())

	stg := envs.Get().GetStorage()

	for _, tt := range tests {
		testTaskObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.task)

		if tt.args.task.Status.State == types.StateProvision && len(tt.args.jobState.task.active) > 0 {

			list := types.NewPodList()

			filter := stg.Filter().Pod().ByTask(tt.args.task.Meta.Namespace, tt.args.job.Meta.Name, tt.args.task.Meta.Name)

			err := stg.List(context.Background(), stg.Collection().Pod(), filter, list, nil)
			if err != nil {
				t.Error(err)
				return
			}

			exists := false

			for _, p := range list.Items {
				_, sl := p.Meta.SelfLink.Parent()
				if sl.String() == tt.args.task.SelfLink().String() {
					exists = true
					break
				}
			}

			if !exists {
				t.Error("pod not created")
			}

		}

	}

}

func TestHandleTaskStateRunning(t *testing.T) {
	type suit struct {
		name string
		args struct {
			jobState *JobState
			task     *types.Task
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle running task"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateRunning, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.active[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateRunning

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.active[wt.SelfLink().String()] = wt

		return s
	}())

	for _, tt := range tests {
		testTaskObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.task)
	}
}

func TestHandleTaskStateError(t *testing.T) {
	type suit struct {
		name string
		args struct {
			jobState *JobState
			task     *types.Task
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle error task"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateError, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Error = true

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.finished = append(s.want.jobState.task.finished, wt)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle error task and set in queued"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateError, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.queue[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Error = true

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.finished = append(s.want.jobState.task.finished, wt)
		delete(s.want.jobState.task.queue, task.SelfLink().String())

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle error task and set in active"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateError, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.active[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Error = true

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.finished = append(s.want.jobState.task.finished, wt)
		delete(s.want.jobState.task.active, task.SelfLink().String())

		return s
	}())

	for _, tt := range tests {
		testTaskObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.task)
	}
}

func TestHandleTaskStateCanceled(t *testing.T) {
	type suit struct {
		name string
		args struct {
			jobState *JobState
			task     *types.Task
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle canceled task"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateCanceled, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Canceled = true

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.finished = append(s.want.jobState.task.finished, wt)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle canceled task in queued"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateCanceled, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.queue[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Canceled = true

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.finished = append(s.want.jobState.task.finished, wt)
		delete(s.want.jobState.task.queue, task.SelfLink().String())

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle canceled task in active"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateCanceled, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.active[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Canceled = true

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.finished = append(s.want.jobState.task.finished, wt)
		delete(s.want.jobState.task.active, task.SelfLink().String())

		return s
	}())

	for _, tt := range tests {
		testTaskObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.task)
	}
}

func TestHandleTaskStateExited(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
			task     *types.Task
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle exited task"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateExited, types.EmptyString)
		task.Status.Done = true

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Done = true

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.finished = append(s.want.jobState.task.finished, wt)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle canceled task in queued"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateExited, types.EmptyString)
		task.Status.Done = true

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.queue[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Done = true

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.finished = append(s.want.jobState.task.finished, wt)
		delete(s.want.jobState.task.queue, task.SelfLink().String())

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle canceled task in active"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateExited, types.EmptyString)
		task.Status.Done = true

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.task.active[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Done = true

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.task.finished = append(s.want.jobState.task.finished, wt)
		delete(s.want.jobState.task.active, task.SelfLink().String())

		return s
	}())

	for _, tt := range tests {
		testTaskObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.task)
	}
}

func TestHandleTaskStateDestroy(t *testing.T) {
	type suit struct {
		name string
		args struct {
			jobState *JobState
			task     *types.Task
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle destroy task with destroy state and without pod"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateDestroy, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateDestroyed

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle destroy task with destroy state and with pod with not destroy state"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateDestroy, types.EmptyString)
		pod, _ := podCreate(task)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateDestroy

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		delete(s.want.jobState.pod.list, wt.SelfLink().String())

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle destroy task with destroy state and with pod with destroy state and without node"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateDestroy, types.EmptyString)
		pod, _ := podCreate(task)
		pod.Status.State = types.StateProvision
		pod.Spec.State.Destroy = true

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateDestroy

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		delete(s.want.jobState.pod.list, wt.SelfLink().String())

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle destroy task with destroy state and with pod with destroy state and with node"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateDestroy, types.EmptyString)
		pod, _ := podCreate(task)
		pod.Status.State = types.StateProvision
		pod.Spec.State.Destroy = true
		pod.Meta.Node = "local"

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateDestroy
		wp := getPodCopy(pod)
		wp.Status.State = types.StateDestroy
		wp.Spec.State.Destroy = true

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.pod.list[wt.SelfLink().String()] = wp

		return s
	}())

	for _, tt := range tests {
		testTaskObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.task)
	}
}

func TestHandleTaskStateDestroyed(t *testing.T) {
	type suit struct {
		name string
		args struct {
			jobState *JobState
			task     *types.Task
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle destroyed task without pod"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateDestroyed, types.EmptyString)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wt.Status.State = types.StateDestroyed

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list = make(map[string]*types.Task, 0)
		s.want.jobState.task.queue = make(map[string]*types.Task, 0)
		s.want.jobState.task.active = make(map[string]*types.Task, 0)
		s.want.jobState.task.finished = make([]*types.Task, 0)
		s.want.jobState.pod.list = make(map[string]*types.Pod, 0)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle destroyed task without pod"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateDestroyed, types.EmptyString)
		pod, _ := podCreate(task)

		s.args.task = task
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateDestroy
		wp := getPodCopy(pod)

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.task.list[wt.SelfLink().String()] = wt
		s.want.jobState.pod.list[wt.SelfLink().String()] = wp

		return s
	}())

	for _, tt := range tests {
		testTaskObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.task)
	}
}

func TestTaskStatusState(t *testing.T) {
	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			jobState *JobState
			task     *types.Task
			pod      *types.Pod
		}
		want struct {
			t *types.Task
		}
		wantErr  bool
		preHook  func() error
		postHook func() error
	}

	var tests []suit

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset(job, types.StateCreated, types.EmptyString)

		pod, _ := podCreate(task)
		assert.NotNil(t, pod)
		pod.Status.State = types.StateProvision

		s := suit{name: "set task created -> provision state (pod: provision)"}
		s.ctx = ctx

		s.args.jobState = js
		s.args.task = task
		s.args.pod = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateProvision

		s.want.t = wt

		s.preHook = func() error {
			return nil
		}
		s.postHook = func() error {

			if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
				return err
			}

			err := stg.Del(ctx, stg.Collection().Pod(), "")
			if !assert.NoError(t, err) {
				return err
			}

			return nil
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)

		pod, _ := podCreate(task)
		assert.NotNil(t, pod)
		pod.Status.State = types.StateProvision

		s := suit{name: "set task provision -> provision state (pod: provision)"}
		s.ctx = ctx

		s.args.jobState = js
		s.args.task = task
		s.args.pod = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateProvision

		s.want.t = wt

		s.preHook = func() error {
			return nil
		}
		s.postHook = func() error {

			if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
				return err
			}

			err := stg.Del(ctx, stg.Collection().Pod(), "")
			if !assert.NoError(t, err) {
				return err
			}

			return nil
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)

		pod, _ := podCreate(task)
		pod.Status.State = types.StateError

		s := suit{name: "set task provision -> exited state (pod: error)"}
		s.ctx = ctx
		s.args.jobState = js
		s.args.task = task
		s.args.pod = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateError
		wt.Status.Error = true

		s.want.t = wt

		s.preHook = func() error {
			return nil
		}
		s.postHook = func() error {

			if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
				return err
			}

			err := stg.Del(ctx, stg.Collection().Pod(), "")
			if !assert.NoError(t, err) {
				return err
			}

			return nil
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)

		pod, _ := podCreate(task)
		pod.Status.State = types.StateDestroy

		s := suit{name: "set task provision -> exited state (pod: destroy)"}
		s.ctx = ctx
		s.args.jobState = js
		s.args.task = task
		s.args.pod = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Done = true
		s.want.t = wt

		s.preHook = func() error {
			return nil
		}
		s.postHook = func() error {

			if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
				return err
			}

			err := stg.Del(ctx, stg.Collection().Pod(), "")
			if !assert.NoError(t, err) {
				return err
			}

			return nil
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)

		pod, _ := podCreate(task)
		pod.Status.State = types.StateDestroyed

		s := suit{name: "set task provision -> exited state (pod: destroyed)"}
		s.ctx = ctx
		s.args.jobState = js
		s.args.task = task
		s.args.pod = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Done = true
		s.want.t = wt

		s.preHook = func() error {
			return nil
		}
		s.postHook = func() error {

			if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
				return err
			}

			err := stg.Del(ctx, stg.Collection().Pod(), "")
			if !assert.NoError(t, err) {
				return err
			}

			return nil
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)

		pod, _ := podCreate(task)
		pod.Status.State = types.StateExited
		pod.Status.Status = types.StateError

		s := suit{name: "set task provision -> exited state (pod: exited and status error)"}
		s.ctx = ctx
		s.args.jobState = js
		s.args.task = task
		s.args.pod = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Done = true
		s.want.t = wt

		s.preHook = func() error {
			return nil
		}
		s.postHook = func() error {

			if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
				return err
			}

			err := stg.Del(ctx, stg.Collection().Pod(), "")
			if !assert.NoError(t, err) {
				return err
			}

			return nil
		}

		return s
	}())

	tests = append(tests, func() suit {

		ctx := context.Background()

		job := getJobAsset("build", "system")
		job.Status.State = types.StateRunning
		js := NewJobState(cluster.NewClusterState(), job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)
		task.Status.Error = true

		pod, _ := podCreate(task)
		pod.Status.State = types.StateExited

		s := suit{name: "set task provision -> exited state (pod: exited)"}
		s.ctx = ctx
		s.args.jobState = js
		s.args.task = task
		s.args.pod = pod

		wt := getTaskCopy(task)
		wt.Status.State = types.StateExited
		wt.Status.Done = true
		s.want.t = wt

		s.preHook = func() error {
			return nil
		}
		s.postHook = func() error {

			if err := stg.Del(ctx, stg.Collection().Job(), ""); err != nil {
				return err
			}

			err := stg.Del(ctx, stg.Collection().Pod(), "")
			if !assert.NoError(t, err) {
				return err
			}

			return nil
		}

		return s
	}())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			if err := tc.preHook(); err != nil {
				assert.NoError(t, err, "pre hook execute error")
				return
			}
			defer func() {
				if err := tc.postHook(); err != nil {
					assert.NoError(t, err, "post hook execute error")
					return
				}
			}()

			err := taskStatusState(tc.args.jobState, tc.args.task, tc.args.pod)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err, "task finish error")
					return
				}
				return
			}

			assert.Equal(t, tc.want.t.Status.State, tc.args.task.Status.State, "task state mismatch")
			assert.Equal(t, tc.want.t.Status.Canceled, tc.args.task.Status.Canceled, "task canceled mismatch")
			assert.Equal(t, tc.want.t.Status.Error, tc.args.task.Status.Error, "task error mismatch")

		})
	}
}

func TestTaskCreate(t *testing.T) {

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	type suit struct {
		ctx  context.Context
		name string
		args struct {
			j  *types.Job
			mf *types.TaskManifest
		}
		want     *types.Task
		wantErr  bool
		err      string
		preHook  func() error
		postHook func() error
	}

	var tests []suit

	tests = append(tests, func() suit {

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		manifest := getTaskManifestAsset("demo")

		s := suit{name: "successful task creation"}
		s.ctx = context.Background()
		s.args.j = job
		s.args.mf = manifest

		s.want = getTaskAssetWithName(job, types.StateCreated, types.EmptyString, "demo")

		return s
	}())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			err := stg.Del(tc.ctx, stg.Collection().Task(), "")
			if !assert.NoError(t, err) {
				return
			}

			task, err := taskCreate(tc.args.j, tc.args.mf)
			if err != nil {
				if !tc.wantErr {
					assert.NoError(t, err, "task create error")
					return
				}
				return
			}

			if tc.want != nil && task == nil {
				assert.NotNil(t, task, "task can not be nil")
				return
			}

			assert.Equal(t, tc.want.Meta.Namespace, task.Meta.Namespace, "task namespace mismatch")
			assert.Equal(t, tc.want.Meta.Job, task.Meta.Job, "task job mismatch")
			assert.Equal(t, tc.want.Meta.SelfLink.String(), task.Meta.SelfLink.String(), "task self_link mismatch")
			assert.Equal(t, tc.want.Status.State, task.Status.State, "task state mismatch")

			var stgTask = new(types.Task)

			if err := stg.Get(tc.ctx, stg.Collection().Task(), task.SelfLink().String(), stgTask, nil); err != nil {
				t.Log(err)
				return
			}

			assert.Equal(t, tc.want.Meta.Namespace, stgTask.Meta.Namespace, "task namespace mismatch")
			assert.Equal(t, tc.want.Meta.Job, stgTask.Meta.Job, "task job mismatch")
			assert.Equal(t, tc.want.Meta.SelfLink.String(), stgTask.Meta.SelfLink.String(), "task self_link mismatch")
			assert.Equal(t, tc.want.Status.State, stgTask.Status.State, "task state mismatch")
		})
	}

}

func getTaskAsset(job *types.Job, state, message string) *types.Task {
	return getTaskAssetWithName(job, state, message, generator.GetUUIDV4())
}

func getTaskAssetWithName(job *types.Job, state, message, name string) *types.Task {

	t := new(types.Task)

	t.Meta.SetDefault()
	t.Meta.Namespace = job.Meta.Namespace
	t.Meta.Job = job.SelfLink().String()
	t.Meta.Name = name
	t.Meta.SelfLink = *types.NewTaskSelfLink(job.Meta.Namespace, job.Meta.Name, t.Meta.Name)

	t.Status.State = state
	t.Status.Message = message
	t.Status.Dependencies.Volumes = make(map[string]types.StatusDependency, 0)
	t.Status.Dependencies.Secrets = make(map[string]types.StatusDependency, 0)
	t.Status.Dependencies.Configs = make(map[string]types.StatusDependency, 0)

	t.Spec.State = job.Spec.State

	t.Spec.Template.Containers = make(types.SpecTemplateContainers, 0)
	t.Spec.Template.Containers = append(t.Spec.Template.Containers, &types.SpecTemplateContainer{
		Name: "demo",
	})

	return t
}

func getTaskCopy(task *types.Task) *types.Task {
	t := *task
	return &t
}

func compareTaskProperties(old *types.Task, new *types.Task) error {

	compareStringSlice := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}

	// =====================================
	// Check task meta
	// =====================================
	if old.Meta.Namespace != new.Meta.Namespace {
		return errors.New("task meta namespace is different")
	}
	if old.Meta.Job != new.Meta.Job {
		return errors.New("task meta job is different")
	}
	if old.Meta.SelfLink.String() != new.Meta.SelfLink.String() {
		return errors.New("task meta self link is different")
	}

	// =====================================
	// Check task status
	// =====================================
	if old.Status.State != new.Status.State {
		return errors.New("task status state is different")
	}
	if old.Status.Error != new.Status.Error {
		return errors.New("task status error is different")
	}
	if old.Status.Canceled != new.Status.Canceled {
		return errors.New("task status canceled is different")
	}
	if old.Status.Done != new.Status.Done {
		return errors.New("task status done is different")
	}
	if old.Status.Message != new.Status.Message {
		return errors.New("task status message is different")
	}
	if len(old.Status.Dependencies.Configs) != len(new.Status.Dependencies.Configs) {
		return errors.New("task status configs dependency is different")
	}
	for k, v := range old.Status.Dependencies.Configs {
		item, ok := new.Status.Dependencies.Configs[k]
		if !ok {
			return errors.New("task status configs dependency is different")
		}
		if v.Type != item.Type {
			return errors.New("task status configs dependency is different")
		}
		if v.Name != item.Name {
			return errors.New("task status configs dependency is different")
		}
		if v.Status != item.Status {
			return errors.New("task status configs dependency is different")
		}
	}
	if len(old.Status.Dependencies.Volumes) != len(new.Status.Dependencies.Volumes) {
		return errors.New("task status volumes dependency is different")
	}
	for k, v := range old.Status.Dependencies.Volumes {
		item, ok := new.Status.Dependencies.Volumes[k]
		if !ok {
			return errors.New("task status volumes dependency is different")
		}
		if v.Type != item.Type {
			return errors.New("task status volumes dependency is different")
		}
		if v.Name != item.Name {
			return errors.New("task status volumes dependency is different")
		}
		if v.Status != item.Status {
			return errors.New("task status volumes dependency is different")
		}
	}
	if len(old.Status.Dependencies.Secrets) != len(new.Status.Dependencies.Secrets) {
		return errors.New("task status secrets dependency is different")
	}
	for k, v := range old.Status.Dependencies.Secrets {
		item, ok := new.Status.Dependencies.Secrets[k]
		if !ok {
			return errors.New("task status secrets dependency is different")
		}
		if v.Type != item.Type {
			return errors.New("task status secrets dependency is different")
		}
		if v.Name != item.Name {
			return errors.New("task status secrets dependency is different")
		}
		if v.Status != item.Status {
			return errors.New("task status secrets dependency is different")
		}
	}

	if old.Status.Pod.SelfLink != new.Status.Pod.SelfLink {
		return errors.New("task status pod self link is different")
	}
	if old.Status.Pod.State != new.Status.Pod.State {
		return errors.New("task status pod state is different")
	}
	if old.Status.Pod.Status != new.Status.Pod.Status {
		return errors.New("task status pod status is different")
	}
	if len(old.Status.Pod.Runtime.Services) != len(new.Status.Pod.Runtime.Services) {
		return errors.New("task status pod runtime services is different")
	}
	for k := range old.Status.Pod.Runtime.Services {
		_, ok := new.Status.Pod.Runtime.Services[k]
		if !ok {
			return errors.New("task status pod runtime services is different")
		}
		//TODO: check service properties
	}
	if len(old.Status.Pod.Runtime.Pipeline) != len(new.Status.Pod.Runtime.Pipeline) {
		return errors.New("task status pod runtime pipeline is different")
	}
	for k, v := range old.Status.Pod.Runtime.Pipeline {
		item, ok := new.Status.Pod.Runtime.Pipeline[k]
		if !ok {
			return errors.New("task status pod runtime pipeline is different")
		}
		if v.Status != item.Status {
			return errors.New("task status pod runtime pipeline is different")
		}
		if v.Error != item.Error {
			return errors.New("task status pod runtime pipeline is different")
		}
		if v.Message != item.Message {
			return errors.New("task status pod runtime pipeline is different")
		}
		if len(v.Commands) != len(item.Commands) {
			return errors.New("task status pod runtime pipeline is different")
		}
		//TODO: check commands
	}

	// =====================================
	// Check task spec
	// =====================================
	if old.Spec.State.Destroy != new.Spec.State.Destroy {
		return errors.New("task spec state destroy is different")
	}
	if old.Spec.State.Cancel != new.Spec.State.Cancel {
		return errors.New("task spec state cancel is different")
	}
	if old.Spec.State.Maintenance != new.Spec.State.Maintenance {
		return errors.New("task spec state maintenance is different")
	}
	if old.Spec.Runtime.Updated != new.Spec.Runtime.Updated {
		return errors.New("task spec runtime updated is different")
	}
	if len(old.Spec.Runtime.Services) != len(new.Spec.Runtime.Services) {
		return errors.New("task spec runtime services is different")
	}
	if !compareStringSlice(old.Spec.Runtime.Services, new.Spec.Runtime.Services) {
		return errors.New("task spec runtime services is different")
	}
	if len(old.Spec.Runtime.Tasks) != len(new.Spec.Runtime.Tasks) {
		return errors.New("task spec runtime tasks is different")
	}
	for i := range old.Spec.Runtime.Tasks {
		if old.Spec.Runtime.Tasks[i].Name != new.Spec.Runtime.Tasks[i].Name {
			return errors.New("task spec runtime tasks is different")
		}
		if old.Spec.Runtime.Tasks[i].Container != new.Spec.Runtime.Tasks[i].Container {
			return errors.New("task spec runtime tasks is different")
		}
		if !compareStringSlice(old.Spec.Runtime.Tasks[i].Commands, new.Spec.Runtime.Tasks[i].Commands) {
			return errors.New("task spec runtime services is different")
		}
		// TODO: check old.Spec.Container.Tasks[i].EnvVars and new.Spec.Container.Tasks[i].EnvVars
	}
	if old.Spec.Selector.Updated != new.Spec.Selector.Updated {
		return errors.New("task spec selector updated is different")
	}
	if old.Spec.Selector.Node != new.Spec.Selector.Node {
		return errors.New("task spec selector node is different")
	}
	if len(old.Spec.Selector.Labels) != len(new.Spec.Selector.Labels) {
		return errors.New("task spec selector labels is different")
	}
	for k, v := range old.Spec.Selector.Labels {
		item, ok := new.Spec.Selector.Labels[k]
		if !ok {
			return errors.New("task status pod runtime labels is different")
		}
		if v != item {
			return errors.New("task status pod runtime labels is different")
		}
	}
	if old.Spec.Template.Updated != new.Spec.Template.Updated {
		return errors.New("task spec template updated is different")
	}
	if old.Spec.Template.Termination != new.Spec.Template.Termination {
		return errors.New("task spec template termination is different")
	}
	if len(old.Spec.Template.Containers) != len(new.Spec.Template.Containers) {
		return errors.New("task spec template containers is different")
	}
	// TODO compare old.Spec.Template.Containers and new.Spec.Template.Containers properties
	if len(old.Spec.Template.Volumes) != len(new.Spec.Template.Volumes) {
		return errors.New("task spec template volumes is different")
	}
	// TODO compare old.Spec.Template.Volumes and new.Spec.Template.Volumes properties
	return nil
}

func getTaskManifestAsset(name string) *types.TaskManifest {
	t := new(types.TaskManifest)
	_ = json.Unmarshal([]byte(taskManifest), t)
	t.Meta.Name = &name

	return t
}

const taskManifest = `
{
  "Meta": {
    "Name": "demo"
  },
  "Spec": {
    "Container": {
      "services": [
        "s_etcd",
        "s_dind"
      ],
      "tasks": [
        {
          "name": "clone:github.com/lastbackend/lastbackend",
          "container": "r_builder",
          "commands": [
            "lb clone -v github -o lastbackend -n lastbackend -b master /data/"
          ]
        },
        {
          "name": "step:test",
          "container": "p_test",
          "commands": [
            "apt-get -y install openssl",
            "mkdir -pod ${GOPATH}/src/github.com/lastbackend/lastbackend",
            "cp -r /data/. ${GOPATH}/src/github.com/lastbackend/lastbackend",
            "cd ${GOPATH}/src/github.com/lastbackend/lastbackend",
            "make deps",
            "make test"
          ]
        },
        {
          "name": "build:hub.lstbknd.net/lastbackend/lastbackend:master",
          "container": "r_builder",
          "commands": [
            "lb build -i hub.lstbknd.net/lastbackend/lastbackend:master -f ./images/lastbackend/Dockerfile .",
            "lb push hub.lstbknd.net/lastbackend/lastbackend:master"
          ]
        }
      ]
    },
    "Template": {
      "containers": [
        {
          "name": "s_etcd",
          "command": "/usr/local/bin/etcd --data-dir=/etcd-data --name node --initial-advertise-peer-urls http://127.0.0.1:2380 --listen-peer-urls http://127.0.0.1:2380 --advertise-client-urls http://127.0.0.1:2379 --listen-client-urls http://127.0.0.1:2379 --initial-cluster node=http://127.0.0.1:2380",
          "image": {
            "name": "quay.io/coreos/etcd:latest"
          }
        },
        {
          "name": "p_test",
          "workdir": "/data/",
          "env": [
            {
              "name": "ENV_GIT_TOKEN",
              "secret": {
                "name": "vault:lastbackend:token",
                "key": "github"
              }
            },
            {
              "name": "ENV_DOCKER_TOKEN",
              "secret": {
                "name": "vault:lastbackend:token",
                "key": "docker"
              }
            },
            {
              "name": "DOCKER_HOST",
              "value": "tcp://127.0.0.1:2375"
            }
          ],
          "volumes": [
            {
              "name": "data",
              "path": "/data/"
            }
          ],
          "image": {
            "name": "golang:stretch"
          },
          "security": {
            "privileged": false
          }
        },
        {
          "name": "s_dind",
          "image": {
            "name": "docker:dind"
          },
          "security": {
            "privileged": true
          }
        },
        {
          "name": "r_builder",
          "workdir": "/data/",
          "env": [
            {
              "name": "ENV_GIT_TOKEN",
              "secret": {
                "name": "vault:lastbackend:token",
                "key": "github"
              }
            },
            {
              "name": "ENV_DOCKER_TOKEN",
              "secret": {
                "name": "vault:lastbackend:token",
                "key": "docker"
              }
            },
            {
              "name": "DOCKER_HOST",
              "value": "tcp://127.0.0.1:2375"
            }
          ],
          "volumes": [
            {
              "name": "data",
              "path": "/data/"
            }
          ],
          "image": {
            "name": "index.lstbknd.net/lastbackend/builder"
          }
        }
      ],
      "volumes": [
        {
          "name": "data"
        }
      ]
    }
  }
}`
