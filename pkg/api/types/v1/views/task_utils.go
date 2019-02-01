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

func (tw *TaskView) New(obj *types.Task, pods *types.PodList) *Task {
	t := new(Task)

	t.ToMeta(obj.Meta)
	t.ToStatus(obj.Status)
	t.ToSpec(obj.Spec)

	if pods != nil {
		t.Pods = make(PodList, 0)
		t.JoinPods(pods)
	}

	return t
}

func (t *Task) ToMeta(obj types.TaskMeta) {
	tm := TaskMeta{}

	tm.Namespace = obj.Namespace
	tm.Job = obj.Job
	tm.Name = obj.Name

	tm.SelfLink = obj.SelfLink
	tm.Description = obj.Description

	tm.Labels = obj.Labels
	tm.Created = obj.Created
	tm.Updated = obj.Updated

	t.Meta = tm
}

func (t *Task) ToStatus(obj types.TaskStatus) {
	ts := TaskStatus{
		State:   obj.State,
		Message: obj.Message,
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

func (t *Task) JoinPods(pods *types.PodList) {

	for _, p := range pods.Items {

		if p.Meta.Namespace != t.Meta.Namespace {
			continue
		}

		if p.Meta.Parent.Kind != types.KindTask {
			continue
		}

		if p.Meta.Parent.SelfLink != t.Meta.SelfLink {
			continue
		}

		t.Pods[p.Meta.SelfLink] = new(PodView).New(p)
	}
}

func (tw *TaskView) NewList(obj *types.TaskList, pods *types.PodList) *TaskList {

	if obj == nil {
		return nil
	}

	tl := make(TaskList, 0)
	for _, v := range obj.Items {
		tl = append(tl, tw.New(v, pods))
	}

	return &tl
}
