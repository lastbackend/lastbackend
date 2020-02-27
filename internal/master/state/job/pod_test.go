//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/internal/master/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/generator"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func testPodObserver(t *testing.T, name, werr string, wjs *JobState, js *JobState, pod *types.Pod) {
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
		err := PodObserve(js, pod)
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

func TestHandlePodStateCreated(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
			pod      *types.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod created"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateCreated, types.EmptyString)
		pod := getPodAsset(task, types.StateCreated, types.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = types.StateProvision

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.pod.list[wt.SelfLink().String()] = wp

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.pod)
	}

}

func TestHandlePodStateProvision(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
			pod      *types.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod provision"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)
		pod := getPodAsset(task, types.StateProvision, types.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = types.StateProvision

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.pod.list[wt.SelfLink().String()] = wp

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.pod)
	}

}

func TestHandlePodStateReady(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
			pod      *types.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod ready"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)
		pod := getPodAsset(task, types.StateReady, types.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = types.StateReady

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.pod.list[wt.SelfLink().String()] = wp

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.pod)
	}

}

func TestHandlePodStateError(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
			pod      *types.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod error"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)
		pod := getPodAsset(task, types.StateError, types.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = types.StateError

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.pod.list[wt.SelfLink().String()] = wp

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.pod)
	}

}

func TestHandlePodStateDegradation(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
			pod      *types.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod degradation"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)
		pod := getPodAsset(task, types.StateDegradation, types.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = types.StateDegradation

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		s.want.jobState.pod.list[wt.SelfLink().String()] = wp

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.pod)
	}

}

func TestHandlePodStateDestroy(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
			pod      *types.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod destroy"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)
		pod := getPodAsset(task, types.StateDestroy, types.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = types.StateDestroy

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		delete(s.want.jobState.pod.list, wt.SelfLink().String())

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod destroy and spec destroy true and node empty"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)
		pod := getPodAsset(task, types.StateDestroy, types.EmptyString)
		pod.Spec.State.Destroy = true

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = types.StateDestroyed

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		delete(s.want.jobState.pod.list, wt.SelfLink().String())

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.pod)
	}

}

func TestHandlePodStateDestroyed(t *testing.T) {

	type suit struct {
		name string
		args struct {
			jobState *JobState
			pod      *types.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod destroyed"}

		job := getJobAsset(types.StateWaiting, types.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, types.StateProvision, types.EmptyString)
		pod := getPodAsset(task, types.StateDestroyed, types.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = types.StateDestroyed

		s.want.err = types.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = types.StateWaiting
		delete(s.want.jobState.pod.list, wt.SelfLink().String())

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.pod)
	}

}

func getPodAsset(t *types.Task, state, message string) *types.Pod {

	p := new(types.Pod)

	p.Meta.SetDefault()
	p.Meta.Namespace = t.Meta.Namespace
	p.Meta.Name = strings.Split(generator.GetUUIDV4(), "-")[4][5:]
	p.Meta.Namespace = t.Meta.Namespace

	sl, _ := types.NewPodSelfLink(types.KindTask, t.SelfLink().String(), p.Meta.Name)
	p.Meta.SelfLink = *sl

	p.Status.State = state
	p.Status.Message = message

	if state == types.StateReady {
		p.Status.Running = true
	}

	p.Spec.State = t.Spec.State
	p.Spec.Template = t.Spec.Template

	return p
}

func getPodCopy(pod *types.Pod) *types.Pod {
	p := *pod
	return &p
}
