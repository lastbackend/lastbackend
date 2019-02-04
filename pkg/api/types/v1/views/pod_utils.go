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

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type PodView struct{}

func (pv *PodView) New(pod *types.Pod) Pod {
	p := Pod{}
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
	meta.SelfLink = pod.SelfLink

	meta.Parent.Kind = pod.Parent.Kind
	meta.Parent.SelfLink = pod.Parent.SelfLink

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

	status.Containers = make(PodContainers, 0)
	for _, container := range pod.Containers {
		cv := new(ContainerView)
		status.Containers = append(status.Containers, cv.NewPodContainer(container))
	}
	p.Status = status
}
