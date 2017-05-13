//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package pod

import (
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/controller/context"
	"time"
	"github.com/satori/go.uuid"
)

func PodClone (p *types.Pod) *types.Pod {

	var (
		log = context.Get().GetLogger()
	)

	log.Debug("Create new pod state on service")

	pod := new(types.Pod)
	pod.Meta.SetDefault()
	pod.State.Provision = true
	pod.Spec = p.Spec

	return pod
}

func PodCreate(spec map[string]*types.ServiceSpec) *types.Pod {
	var (
		log = context.Get().GetLogger()
	)

	log.Debug("Create new pod state on service")

	pod := new(types.Pod)
	pod.Meta.SetDefault()
	pod.State.Provision = true
	pod.Spec = podSpecGenerate(spec)

	return pod
}

func PodUpdate(p *types.Pod, spec map[string]*types.ServiceSpec) {

	var (
		log = context.Get().GetLogger()
	)

	log.Debug("Pod update")

	p.Meta.Updated = time.Now()
	p.Spec = podSpecGenerate(spec)

}

func PodRemove(p *types.Pod) {

	var (
		log = context.Get().GetLogger()
	)

	log.Debugf("Mark pod for deletion: %s", p.Meta.Name)
	p.State.Provision = true
	p.State.Ready = false
	p.Spec.State = types.StateDestroy

	for _, c := range p.Containers {
		c.State = types.StateProvision
	}
}

func podSpecGenerate(spec map[string]*types.ServiceSpec) types.PodSpec {

	var (
		log = context.Get().GetLogger()
	)

	log.Debug("Generate new node pod spec")

	var s = types.PodSpec{}
	s.ID = uuid.NewV4().String()
	s.Created = time.Now()
	s.Updated = time.Now()
	s.Containers = make(map[string]*types.ContainerSpec)

	for _, spc := range spec {

		cs := new(types.ContainerSpec)
		cs.Meta.SetDefault()

		cs.Image = types.ImageSpec{
			Name: spc.Image,
			Pull: true,
		}

		for _, port := range spc.Ports {
			cs.Ports = append(cs.Ports, types.ContainerPortSpec{
				ContainerPort: port.Container,
				Protocol:      port.Protocol,
			})
		}

		cs.Command = spc.Command
		cs.Entrypoint = spc.Entrypoint
		cs.Envs = spc.EnvVars
		cs.Quota = types.ContainerQuotaSpec{
			Memory: spc.Memory,
		}

		cs.RestartPolicy = types.ContainerRestartPolicySpec{
			Name:    "always",
			Attempt: 0,
		}

		s.Containers[cs.Meta.ID] = cs
	}

	s.State = types.StateStarted

	return s
}