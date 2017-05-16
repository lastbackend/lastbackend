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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/controller/context"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

func PodCreate(svc *types.Service) *types.Pod {
	var (
		log = context.Get().GetLogger()
	)

	log.Debug("Create new pod state on service")

	pod := new(types.Pod)
	pod.Meta.SetDefault()
	pod.Meta.Name = podNameGenerate(svc)
	pod.State.Provision = true

	return pod
}

func PodSetSpec(p *types.Pod, spec map[string]*types.ServiceSpec) {

	var (
		log = context.Get().GetLogger()
		ids = make(map[string]struct{})
	)

	for id := range spec {
		ids[id] = struct{}{}
	}

	if p.Spec.State == types.StateDestroy {
		return
	}

	if len(p.Spec.Containers) != len(spec) {
		p.Spec = podSpecGenerate(spec)
		p.Meta.Updated = time.Now()
		return
	}

	for id := range p.Spec.Containers {
		if _, ok := ids[id]; !ok {
			log.Debug("Pod update")
			p.Spec = podSpecGenerate(spec)
			p.Meta.Updated = time.Now()
			return
		}
	}
}

func PodRemove(p *types.Pod) {

	var (
		log = context.Get().GetLogger()
	)

	log.Debugf("Mark pod for deletion: %s", p.Meta.Name)
	p.State.Provision = true
	p.State.Ready = false
	p.State.State = types.StateDestroy
	p.Spec.State = types.StateDestroy

	for _, c := range p.Containers {
		c.State = types.StateDestroy
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
		cs.Meta.ID = spc.Meta.ID
		cs.Meta.Labels = spc.Meta.Labels

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

func podNameGenerate(svc *types.Service) string {
	var name, hash string
	for {

		hash = strings.Split(uuid.NewV4().String(), "-")[4]
		name = fmt.Sprintf("%s:%s:%s", svc.Meta.Namespace, svc.Meta.Name, hash[5:])
		if _, ok := svc.Pods[name]; !ok {
			break
		}
	}

	return name
}
