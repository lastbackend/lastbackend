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
	"github.com/lastbackend/lastbackend/internal/master/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/generator"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func testPodObserver(t *testing.T, name, werr string, wjs *JobState, js *JobState, pod *models.Pod) {
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
		if werr != models.EmptyString {

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
			pod      *models.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod created"}

		job := getJobAsset(models.StateWaiting, models.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, models.StateCreated, models.EmptyString)
		pod := getPodAsset(task, models.StateCreated, models.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = models.StateProvision

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateWaiting
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
			pod      *models.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod provision"}

		job := getJobAsset(models.StateWaiting, models.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, models.StateProvision, models.EmptyString)
		pod := getPodAsset(task, models.StateProvision, models.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = models.StateProvision

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateWaiting
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
			pod      *models.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod ready"}

		job := getJobAsset(models.StateWaiting, models.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, models.StateProvision, models.EmptyString)
		pod := getPodAsset(task, models.StateReady, models.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = models.StateReady

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateWaiting
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
			pod      *models.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod error"}

		job := getJobAsset(models.StateWaiting, models.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, models.StateProvision, models.EmptyString)
		pod := getPodAsset(task, models.StateError, models.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = models.StateError

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateWaiting
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
			pod      *models.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod degradation"}

		job := getJobAsset(models.StateWaiting, models.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, models.StateProvision, models.EmptyString)
		pod := getPodAsset(task, models.StateDegradation, models.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = models.StateDegradation

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateWaiting
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
			pod      *models.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod destroy"}

		job := getJobAsset(models.StateWaiting, models.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, models.StateProvision, models.EmptyString)
		pod := getPodAsset(task, models.StateDestroy, models.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = models.StateDestroy

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateWaiting
		delete(s.want.jobState.pod.list, wt.SelfLink().String())

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod destroy and spec destroy true and node empty"}

		job := getJobAsset(models.StateWaiting, models.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, models.StateProvision, models.EmptyString)
		pod := getPodAsset(task, models.StateDestroy, models.EmptyString)
		pod.Spec.State.Destroy = true

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = models.StateDestroyed

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateWaiting
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
			pod      *models.Pod
		}
		want struct {
			err      string
			jobState *JobState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle pod destroyed"}

		job := getJobAsset(models.StateWaiting, models.EmptyString)
		js := getJobStateAsset(job)
		task := getTaskAsset(job, models.StateProvision, models.EmptyString)
		pod := getPodAsset(task, models.StateDestroyed, models.EmptyString)

		s.args.pod = pod
		s.args.jobState = js
		s.args.jobState.task.list[task.SelfLink().String()] = task
		s.args.jobState.pod.list[task.SelfLink().String()] = pod

		wt := getTaskCopy(task)
		wp := getPodCopy(pod)
		wp.Status.State = models.StateDestroyed

		s.want.err = models.EmptyString
		s.want.jobState = getJobStateCopy(s.args.jobState)
		s.want.jobState.job.Status.State = models.StateWaiting
		delete(s.want.jobState.pod.list, wt.SelfLink().String())

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.jobState, tt.args.jobState, tt.args.pod)
	}

}

func getPodAsset(t *models.Task, state, message string) *models.Pod {

	p := new(models.Pod)

	p.Meta.SetDefault()
	p.Meta.Namespace = t.Meta.Namespace
	p.Meta.Name = strings.Split(generator.GetUUIDV4(), "-")[4][5:]
	p.Meta.Namespace = t.Meta.Namespace

	sl, _ := models.NewPodSelfLink(models.KindTask, t.SelfLink().String(), p.Meta.Name)
	p.Meta.SelfLink = *sl

	p.Status.State = state
	p.Status.Message = message

	if state == models.StateReady {
		p.Status.Running = true
	}

	p.Spec.State = t.Spec.State
	p.Spec.Template = t.Spec.Template

	return p
}

func getPodCopy(pod *models.Pod) *models.Pod {
	p := *pod
	return &p
}
