//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package v1

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type PodViewHelper struct{}

func (pv *PodViewHelper) New(pod *types.Pod) PodView {
	p := PodView{}
	p.ID = pod.Meta.Name
	p.Meta = p.toMeta(pod.Meta)
	p.State = p.toState(pod.State)
	p.Spec = p.toSpec(pod.Spec)
	p.Status = p.toStatus(pod.Status)
	return p
}

func (pv *PodView) toMeta(pod types.PodMeta) PodMeta {
	meta := PodMeta{}
	meta.Name = pod.Name
	meta.Description = pod.Description
	meta.SelfLink = pod.SelfLink
	meta.Namespace = pod.Namespace
	meta.Deployment = pod.Deployment
	meta.SelfLink = pod.SelfLink
	meta.Node = pod.Node
	meta.Status = pod.Status
	meta.Updated = pod.Updated
	meta.Created = pod.Created

	return meta
}

func (pv *PodView) toState(pod types.PodState) PodState {
	return PodState{
		Scheduled: pod.Scheduled,
		Provision: pod.Provision,
		Error:     pod.Error,
		Created:   pod.Created,
		Pulling:   pod.Pulling,
		Running:   pod.Running,
		Stopped:   pod.Stopped,
		Destroy:   pod.Destroy,
	}
}

func (pv *PodView) toSpec(pod types.PodSpec) PodSpec {
	return PodSpec{
		Volumes:     pod.Volumes,
		Containers:  pod.Containers,
		Termination: pod.Termination,
	}
}

func (pv *PodView) toStatus(pod types.PodStatus) PodStatus {
	var status = PodStatus{
		Stage:   pod.Stage,
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
	return status
}
