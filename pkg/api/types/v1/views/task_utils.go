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

package views

import "github.com/lastbackend/lastbackend/pkg/distribution/types"

type TaskView struct{}

func (tw *TaskView) New(obj *types.Task) *Task {
	t := new(Task)

	t.ToMeta(obj.Meta)
	t.ToStatus(obj.Status)
	t.ToSpec(obj.Spec)

	return t
}

func (t *Task) ToMeta(obj types.TaskMeta) {
	tm := TaskMeta{}

	tm.Namespace = obj.Namespace
	tm.Job = obj.Job
	tm.Name = obj.Name

	tm.SelfLink = obj.SelfLink.String()
	tm.Description = obj.Description

	tm.Labels = obj.Labels
	tm.Created = obj.Created
	tm.Updated = obj.Updated

	t.Meta = tm
}

func (t *Task) ToStatus(obj types.TaskStatus) {
	ts := TaskStatus{
		State:    obj.State,
		Error:    obj.Error,
		Canceled: obj.Canceled,
		Done:     obj.Done,
		Message:  obj.Message,
		Pod: TaskStatusPod{
			SelfLink: obj.Pod.SelfLink,
			Status:   obj.Pod.Status,
			State:    obj.Pod.State,
		},
	}

	ts.Pod.Runtime = PodStatusRuntime{
		Services: make(PodContainers, 0),
		Pipeline: make([]PodStatusPipelineStep, 0),
	}

	for _, container := range obj.Pod.Runtime.Services {
		cv := new(ContainerView)
		ts.Pod.Runtime.Services = append(ts.Pod.Runtime.Services, cv.NewPodContainer(container))
	}

	for name, step := range obj.Pod.Runtime.Pipeline {

		s := PodStatusPipelineStep{
			Name:    name,
			Status:  step.Status,
			Error:   step.Error,
			Message: step.Message,
		}

		for _, container := range step.Commands {
			cv := new(ContainerView)
			s.Commands = append(s.Commands, cv.NewPodContainer(container))
		}

		ts.Pod.Runtime.Pipeline = append(ts.Pod.Runtime.Pipeline, s)
	}

	t.Status = ts
}

func (t *Task) ToSpec(obj types.TaskSpec) {
	mv := new(ManifestView)
	ts := TaskSpec{
		Template: mv.NewManifestSpecTemplate(obj.Template),
		Selector: mv.NewManifestSpecSelector(obj.Selector),
		Runtime:  mv.NewManifestSpecRuntime(obj.Runtime),
	}
	t.Spec = ts
}

func (tw *TaskView) NewList(obj *types.TaskList) *TaskList {

	if obj == nil {
		return nil
	}

	tl := make(TaskList, 0)
	for _, v := range obj.Items {
		tl = append(tl, tw.New(v))
	}

	return &tl
}
