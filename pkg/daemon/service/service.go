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
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	ctx "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/node"
	"github.com/lastbackend/lastbackend/pkg/daemon/service/routes/request"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"reflect"
	"strconv"
	"time"
)

type service struct {
	Context   context.Context
	Namespace types.Meta
}

func New(ctx context.Context, namespace types.Meta) *service {
	return &service{
		Context:   ctx,
		Namespace: namespace,
	}
}

func (s *service) List() (*types.ServiceList, error) {
	var storage = ctx.Get().GetStorage()
	return storage.Service().ListByNamespace(s.Context, s.Namespace.ID)
}

func (s *service) Get(service string) (*types.Service, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     *types.Service
	)

	if validator.IsUUID(service) {
		svc, err = storage.Service().GetByID(s.Context, s.Namespace.ID, service)
	} else {
		svc, err = storage.Service().GetByName(s.Context, s.Namespace.ID, service)
	}

	if err != nil {
		log.Error("Error: find service by name", err.Error())
		return svc, err
	}

	return svc, nil
}

func (s *service) Create(rq *request.RequestServiceCreateS) (*types.Service, error) {

	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     = new(types.Service)
	)

	svc.Meta = types.ServiceMeta{}
	svc.Meta.ID = uuid.NewV4().String()
	svc.Meta.Name = rq.Name
	svc.Meta.Region = rq.Region
	svc.Meta.Namespace = s.Namespace.Name
	svc.Meta.Description = rq.Description
	svc.Meta.Updated = time.Now()
	svc.Meta.Created = time.Now()

	config, err := createConfig(rq.Config)
	if err != nil {
		log.Errorf("Error: create config from request opts : %s", err.Error())
		return svc, err
	}

	svc.Config = *config

	log.Debugf("Service: Create: add pods : %d", svc.Config.Replicas)
	for i := 0; i < svc.Config.Replicas; i++ {
		log.Debug("Service: Create: add new pod")
		if err := s.AddPod(svc); err != nil {
			log.Errorf("Service: Create: add new pod error: %s", err.Error())
			return svc, err
		}
	}

	svc, err = storage.Service().Insert(s.Context, svc)
	if err != nil {
		log.Errorf("Error: insert service to db : %s", err.Error())
		return svc, err
	}

	return svc, nil
}

func (s *service) Update(service *types.Service, rq *request.RequestServiceUpdateS) (*types.Service, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     *types.Service
	)

	log.Debug("Service: Update: update start")

	if rq.Name != "" {
		service.Meta.Name = rq.Name
	}

	if rq.Description != "" {
		service.Meta.Description = rq.Description
	}

	if rq.Domains != nil {
		service.Domains = rq.Domains
	}

	if err := updateConfig(rq.Config, &service.Config); err != nil {
		log.Error("Error: update service config from request opts", err)
		return svc, err
	}

	// Update pod spec
	spec := s.GenerateSpec(service)

	log.Debugf("Service: Update: pods count: %d", len(service.Pods))
	for _, pod := range service.Pods {
		log.Debugf("Service: Update: pod %s update", pod.Meta.ID)
		pod.Spec = spec
	}

	svc, err = storage.Service().Update(s.Context, service)
	if err != nil {
		log.Error("Error: insert service to db", err)
		return svc, err
	}

	return svc, nil
}

func (s *service) Remove(service *types.Service) error {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	if err := storage.Service().Remove(s.Context, service); err != nil {
		log.Error("Error: insert service to db", err)
		return err
	}
	return nil
}

func (s *service) AddPod(service *types.Service) error {

	var (
		log = ctx.Get().GetLogger()
	)

	log.Debug("Create new pod state on service")

	pod := new(types.Pod)
	pod.State.State = "running"
	pod.Meta.ID = uuid.NewV4().String()
	pod.Meta.Created = time.Now()
	pod.Meta.Updated = time.Now()
	pod.Spec = s.GenerateSpec(service)

	n, err := node.New().Allocate(s.Context, pod.Spec)
	if err != nil {
		return err
	}

	log.Debugf("Service: Add pod: Node meta: %s", n.Meta)
	pod.Meta.Hostname = n.Meta.Hostname

	service.Pods = append(service.Pods, pod)

	return nil
}

func (s *service) DelPod(service *types.Service) error {

	var (
		log = ctx.Get().GetLogger()
		pod *types.Pod
	)

	log.Debug("Create new pod state on service")

	for i := len(service.Pods); i >= 0; i-- {
		pod = service.Pods[i-1]
		if pod.State.State != "deleting" {
			break
		}
	}

	pod.State.State = "deleting"
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

func (s *service) Scale(c context.Context, service *types.Service) (*types.Service, error) {
	var (
		log      = ctx.Get().GetLogger()
		storage  = ctx.Get().GetStorage()
		pod      *types.Pod
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
		log.Debug("Service: Replicas not needed, replicas equal")
		return service, nil
	}

	if replicas < service.Config.Replicas {
		for i := 0; i < (service.Config.Replicas - replicas); i++ {
			if err := s.AddPod(service); err != nil {
				return service, err
			}
		}
	}

	if replicas > service.Config.Replicas {
		for i := 0; i < (replicas - service.Config.Replicas); i++ {
			if err := s.DelPod(service); err != nil {
				return service, err
			}
		}
	}

	svc, err := storage.Service().Update(s.Context, service)
	if err != nil {
		log.Error("Error: insert service to db", err)
		return svc, err
	}

	return svc, nil
}

func (s *service) GenerateSpec(service *types.Service) types.PodSpec {

	var (
		log = ctx.Get().GetLogger()
	)

	log.Debug("Generate new node pod spec")
	var spec = types.PodSpec{}
	spec.ID = uuid.NewV4().String()
	spec.Created = time.Now()
	spec.Updated = time.Now()

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
	cs.Quota = types.ContainerQuotaSpec{
		Memory: service.Config.Memory,
	}

	cs.RestartPolicy = types.ContainerRestartPolicySpec{
		Name:    "always",
		Attempt: 0,
	}

	spec.Containers = append(spec.Containers, cs)

	var state = new(types.PodState)
	state.State = service.State.State

	return spec
}

func createConfig(opts map[string]interface{}) (*types.ServiceConfig, error) {
	config := new(types.ServiceConfig)
	if err := patchConfig(opts, config); err != nil {
		return nil, err
	}
	return config, nil
}

func updateConfig(opts map[string]interface{}, config *types.ServiceConfig) error {
	if config == nil {
		config = new(types.ServiceConfig)
	}

	tmp := make(map[string]interface{})
	for k, v := range opts {
		tmp[k] = v
	}
	delete(tmp, "image")

	return patchConfig(tmp, config)
}

func patchConfig(opts map[string]interface{}, config *types.ServiceConfig) error {

	config.Replicas = int(1)
	config.Memory   = int64(32)

	if val, ok := opts["replicas"]; ok {
		switch reflect.ValueOf(val).Kind() {
		case reflect.Float64:
			config.Replicas = int(val.(float64))
		case reflect.String:
			i, err := strconv.Atoi(val.(string))
			if err != nil {
				return err
			}
			config.Replicas = i
		default:
			return errors.New("replicas incorrect format")
		}
	}

	if val, ok := opts["memory"]; ok {

		switch reflect.ValueOf(val).Kind() {
		case reflect.Float64:
			config.Memory = int64(val.(float64))
		case reflect.String:
			i, err := strconv.ParseInt(val.(string), 10, 64)
			if err != nil {
				return err
			}
			config.Memory = i
		default:
			return errors.New("memory incorrect format")
		}
	}

	if val, ok := opts["image"]; ok {
		if reflect.ValueOf(val).Kind() != reflect.String {
			return errors.New("image incorrect format")
		}
		config.Image = val.(string)
	}

	if val, ok := opts["entrypoint"]; ok {
		if err := json.Unmarshal([]byte(val.(string)), &config.Entrypoint); err != nil {
			return errors.New("entrypoint incorrect format")
		}
	}

	if val, ok := opts["command"]; ok {
		if err := json.Unmarshal([]byte(val.(string)), &config.Command); err != nil {
			return errors.New("command incorrect format")
		}
	}

	if val, ok := opts["env"]; ok {
		if err := json.Unmarshal([]byte(val.(string)), &config.EnvVars); err != nil {
			return errors.New("env incorrect format")
		}
	}

	if val, ok := opts["ports"]; ok && reflect.ValueOf(val).Kind() == reflect.Struct {
		if err := json.Unmarshal([]byte(val.(string)), &config.Ports); err != nil {
			return errors.New("ports incorrect format")
		}
	}

	return nil
}
