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

package views

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
)

type PodView struct{}

func (pv *PodView) New(pod *types.Pod) *Pod {
	p := new(Pod)
	p.SetMeta(pod.Meta)
	p.SetStatus(pod.Status)
	p.SetSpec(pod.Spec)

	return p
}

func (p *Pod) SetMeta(pod types.PodMeta) {
	meta := PodMeta{}
	meta.Namespace = pod.Namespace
	meta.Name = pod.Name

	meta.Description = pod.Description
	meta.SelfLink = pod.SelfLink.String()

	meta.Node = pod.Node
	meta.Status = pod.Status
	meta.Updated = pod.Updated
	meta.Created = pod.Created

	p.Meta = meta
}

func (p *Pod) SetSpec(pod types.PodSpec) {
	mv := new(ManifestView)
	p.Spec = PodSpec{
		State: PodSpecState{
			Destroy:     pod.State.Destroy,
			Maintenance: pod.State.Maintenance,
		},
		Template: mv.NewManifestSpecTemplate(pod.Template),
	}
}

func (p *Pod) SetStatus(pod types.PodStatus) {
	var status = PodStatus{
		State:   pod.State,
		Message: pod.Message,
	}

	status.Network.HostIP = pod.Network.HostIP
	status.Network.PodIP = pod.Network.PodIP

	status.Steps = make(PodSteps, 0)
	for key, step := range pod.Steps {
		status.Steps[key] = PodStep{
			Ready:     step.Ready,
			Timestamp: step.Timestamp,
		}
	}
	status.Runtime = PodStatusRuntime{
		Services: make(PodContainers, 0),
		Pipeline: make([]PodStatusPipelineStep, 0),
	}

	for _, container := range pod.Runtime.Services {
		cv := new(ContainerView)
		status.Runtime.Services = append(status.Runtime.Services, cv.NewPodContainer(container))
	}

	for name, step := range pod.Runtime.Pipeline {

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

	}

	p.Status = status
}

func (pv *PodView) NewList(obj *types.PodList) *PodList {
	pl := make(PodList, 0)
	for _, d := range obj.Items {
		v := new(PodView)
		p := v.New(d)
		pl = append(pl, p)
	}
	return &pl
}

func (pl *PodList) ToJson() ([]byte, error) {
	return json.Marshal(pl)
}
