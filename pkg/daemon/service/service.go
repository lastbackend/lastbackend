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

package service

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	ctx "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/node"
	"github.com/lastbackend/lastbackend/pkg/daemon/pod"
	"github.com/lastbackend/lastbackend/pkg/daemon/service/routes/request"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"github.com/satori/go.uuid"
	"time"
)

type service struct {
	Context   context.Context
	Namespace string
}

func New(ctx context.Context, namespace string) *service {
	return &service{
		Context:   ctx,
		Namespace: namespace,
	}
}

func (s *service) List() (*types.ServiceList, error) {
	var storage = ctx.Get().GetStorage()
	return storage.Service().ListByProject(s.Context, s.Namespace)
}

func (s *service) Get(service string) (*types.Service, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     *types.Service
	)

	if validator.IsUUID(service) {
		svc, err = storage.Service().GetByID(s.Context, s.Namespace, service)
	} else {
		svc, err = storage.Service().GetByName(s.Context, s.Namespace, service)
	}

	if err != nil {
		log.Error("Error: find service by name", err.Error())
		return svc, err
	}

	return svc, nil
}

func (s *service) Create(rq *request.RequestServiceCreateS) (*types.Service, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     *types.Service
	)

	svc, err = storage.Service().Insert(s.Context, s.Namespace, rq.Name, rq.Description, rq.Config)
	if err != nil {
		log.Errorf("Error: insert service to db : %s", err.Error())
		return svc, err
	}

	return svc, nil
}

func (s *service) Update(service *types.Service) (*types.Service, error) {
	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     *types.Service
	)

	svc, err = storage.Service().Update(s.Context, s.Namespace, service)
	if err != nil {
		log.Error("Error: insert service to db", err)
		return svc, err
	}

	n := node.New()
	for _, pod := range svc.Pods {
		// Get node hostname and update pod spec
		if err := n.PodSpecUpdate(s.Context, pod.Meta.Hostname, s.GenerateSpec(svc, pod.Meta)); err != nil {
			log.Errorf("Service: service update: Pod spec upds %s", err.Error())
			return svc, err
		}

	}

	return svc, nil
}

func (s *service) Remove(service *types.Service) error {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	if err := storage.Service().Remove(s.Context, s.Namespace, service); err != nil {
		log.Error("Error: insert service to db", err)
		return err
	}
	return nil
}

func (s *service) AddPod(service *types.Service) error {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.Debug("Create new pod state on service")

	pod := new(types.PodNodeState)
	pod.State.State = "creating"
	pod.Meta.ID = uuid.NewV4().String()
	pod.Meta.Created = time.Now()
	pod.Meta.Updated = time.Now()

	if err := storage.Pod().Insert(s.Context, service.Meta.Namespace, service.Meta.ID, pod); err != nil {
		log.Errorf("Service: Add Pod: insert into storage error: %s", err.Error())
		return err
	}

	service.Pods = append(service.Pods, pod)

	n := node.New()
	n.Allocate(s.Context, s.GenerateSpec(service, pod.Meta))

	return nil
}

func (s *service) DelPod(service *types.Service) error {

	var (
		log     = ctx.Get().GetLogger()
		pod     *types.PodNodeState
		storage = ctx.Get().GetStorage()
	)

	log.Debug("Create new pod state on service")

	for i := len(service.Pods); i >= 0; i-- {
		pod = service.Pods[i-1]
		if pod.State.State != "deleting" {
			break
		}
	}

	pod.State.State = "deleting"
	storage.Pod().Update(s.Context, service.Meta.Namespace, service.Meta.ID, pod)

	return nil
}

func (s *service) SetPods(c context.Context, pods []types.PodNodeState) error {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	for _, pod := range pods {
		svc, err := storage.Service().GetByPodID(c, pod.Meta.ID)
		if err != nil {
			log.Errorf("Error: get pod from db: %s", err)
			return err
		}

		if svc == nil {
			continue
		}

		if p, e := storage.Pod().GetByID(c, svc.Meta.Namespace, svc.Meta.ID, pod.Meta.ID); p == nil || e != nil {

			if err != nil {
				log.Errorf("Error: get pod from db: %s", err)
				return err
			}

			if p == nil {
				log.Warnf("Pod not found, skip setting: %s", pod.Meta.ID)
			}

		}

		if err := storage.Pod().Update(c, svc.Meta.Namespace, svc.Meta.ID, &pod); err != nil {
			log.Errorf("Error: set pod to db: %s", err)
			return err
		}
	}

	return nil
}

func (s *service) Scale(c context.Context, service *types.Service) error {
	var (
		log      = ctx.Get().GetLogger()
		pod      *types.PodNodeState
		replicas int
	)

	for i := 0; i < len(service.Pods); i++ {
		pod = service.Pods[i]
		if pod.State.State == "deleting" {
			continue
		}
		replicas++
	}

	if replicas == service.Config.Replicas {
		log.Debug("Service: Scale: Scale not needed, replicas equal")
		return nil
	}

	if replicas < service.Config.Replicas {
		for i := 0; i < (service.Config.Replicas - replicas); i++ {
			if err := s.AddPod(service); err != nil {
				return err
			}
		}
	}

	if replicas > service.Config.Replicas {
		for i := 0; i < (replicas - service.Config.Replicas); i++ {
			if err := s.DelPod(service); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *service) GenerateSpec(service *types.Service, pod types.PodMeta) *types.PodNodeSpec {

	var (
		log = ctx.Get().GetLogger()
	)

	log.Debug("Generate new node pod spec")
	var spec = new(types.PodNodeSpec)
	spec.Meta = pod
	spec.State.State = "provision"

	spec.Spec.ID = uuid.NewV4().String()
	spec.Spec.Created = time.Now()
	spec.Spec.Updated = time.Now()

	cs := new(types.ContainerSpec)
	cs.Image = types.ImageSpec{
		Name: service.Config.Image,
		Pull: true,
	}

	for _, port := range service.Config.Ports {
		cs.Ports = append(cs.Ports, types.ContainerPortSpec{
			ContainerPort: port.Container,
			Protocol:      port.Protocol,
		})
	}

	cs.Command = service.Config.Command
	cs.Entrypoint = service.Config.Entrypoint
	cs.Envs = service.Config.EnvVars
	cs.Args = service.Config.Args
	cs.Quota = types.ContainerQuotaSpec{
		Memory: service.Config.Memory,
	}

	cs.RestartPolicy = types.ContainerRestartPolicySpec{
		Name:    "always",
		Attempt: 0,
	}

	spec.Spec.Containers = append(spec.Spec.Containers, cs)

	var state = new(types.PodState)
	state.State = service.State.State

	return spec
}
