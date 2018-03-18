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

package distribution

import (
	"context"
	"fmt"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

type IService interface {
	Get(namespace, service string) (*types.Service, error)
	List(namespace string) (map[string]*types.Service, error)
	Create(namespace *types.Namespace, opts *types.ServiceCreateOptions) (*types.Service, error)
	Update(service *types.Service, opts *types.ServiceUpdateOptions) (*types.Service, error)
	Destroy(service *types.Service) (*types.Service, error)
	Remove(service *types.Service) error
}

type Service struct {
	context context.Context
	storage storage.Storage
}

func (s *Service) List(namespace string) (map[string]*types.Service, error) {

	log.V(logLevel).Debug("Service: list service")

	items, err := s.storage.Service().ListByNamespace(s.context, namespace)
	if err != nil {
		log.Debug(1)
		log.V(logLevel).Error("Service: list service namespace %s err: %s", namespace, err)
		return nil, err
	}

	log.V(logLevel).Debugf("Service: list service namespace %s result: %d", namespace, len(items))

	return items, nil
}

func (s *Service) Get(namespace, service string) (*types.Service, error) {

	log.V(logLevel).Debugf("Service: get service %s", service)

	svc, err := s.storage.Service().Get(s.context, namespace, service)
	if err != nil {

		if err.Error() == store.ErrEntityNotFound {
			log.V(logLevel).Warnf("Service: Get: service by name `%s` not found", service)
			return nil, nil
		}

		log.V(logLevel).Errorf("Service: get service by name `%s` err: %s", service, err)
		return nil, err
	}

	return svc, nil
}

func (s *Service) Create(namespace *types.Namespace, opts *types.ServiceCreateOptions) (*types.Service, error) {

	log.V(logLevel).Debugf("Service: create service %#v", opts)

	if opts == nil {
		opts = new(types.ServiceCreateOptions)
	}

	// ------------------------------------------------
	// Create service
	// ------------------------------------------------

	service := new(types.Service)
	service.Meta.SetDefault()
	service.Spec.SetDefault()
	service.Meta.Name = *opts.Name
	service.Meta.SelfLink = fmt.Sprintf("%s/service/%s", namespace.Meta.SelfLink, *opts.Name)
	service.Meta.Endpoint = strings.ToLower(fmt.Sprintf("%s-%s.%s", *opts.Name, namespace.Meta.Endpoint, viper.GetString("domain.internal")))
	service.Meta.Namespace = namespace.Meta.Name

	if opts.Replicas != nil {
		service.Spec.Replicas = types.DEFAULT_SERVICE_REPLICAS
	}

	c := types.SpecTemplateContainer{}
	c.SetDefault()
	c.Role = types.ContainerRolePrimary

	service.Spec.Template.Containers = append(service.Spec.Template.Containers, c)

	if err := s.storage.Service().Insert(s.context, service); err != nil {
		log.V(logLevel).Errorf("Service: insert service err: %s", err)
		return nil, err
	}

	if opts.Spec != nil {
		patchSpec(&service.Spec, *opts.Spec)
	}

	return service, nil
}

func (s *Service) Update(service *types.Service, opts *types.ServiceUpdateOptions) (*types.Service, error) {

	log.V(logLevel).Debugf("Service: update service %#v -> %#v", service, opts)

	if opts == nil {
		opts = new(types.ServiceUpdateOptions)
	}

	if opts.Description != nil {
		service.Meta.Description = *opts.Description
	}

	if opts.Spec != nil {
		patchSpec(&service.Spec, *opts.Spec)
	}

	if opts.Replicas != nil {
		service.Spec.Replicas = *opts.Replicas
	}

	if err := s.storage.Service().Update(s.context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service err: %s", err)
		return nil, err
	}

	return service, nil
}

func (s *Service) Destroy(service *types.Service) (*types.Service, error) {

	log.V(logLevel).Debugf("Service: remove service %#v", service)

	service.State.Destroy = true
	service.Spec.Replicas = int(0) // Delete all pods

	err := s.storage.Service().Update(s.context, service)
	if err != nil {
		log.V(logLevel).Errorf("Service: update service state err: %s", err)
		return nil, err
	}

	return service, nil
}

func (s *Service) Remove(service *types.Service) error {

	log.V(logLevel).Debugf("Service: remove service %#v", service)

	err := s.storage.Service().Remove(s.context, service)
	if err != nil {
		log.V(logLevel).Errorf("Service: update service state err: %s", err)
		return err
	}

	return nil
}

func patchSpec(tpl *types.ServiceSpec, spec types.ServiceOptionsSpec) bool {

	var (
		m bool
		f bool
		i int
	)

	c := types.SpecTemplateContainer{}

	for i, t := range tpl.Template.Containers {
		if t.Role == types.ContainerRolePrimary {
			c = tpl.Template.Containers[i]
			f = true
		}
	}

	if !f {
		c.SetDefault()
		c.Role = types.ContainerRolePrimary
		m = true
	}

	if spec.Command != nil {
		c.Exec.Command = strings.Split(*spec.Command, " ")
		m = true
	}

	if spec.Entrypoint != nil {
		c.Exec.Entrypoint = strings.Split(*spec.Entrypoint, " ")
		m = true
	}

	if spec.Ports != nil {
		c.Ports = types.SpecTemplateContainerPorts{}
		for _, val := range *spec.Ports {
			c.Ports = append(c.Ports, types.SpecTemplateContainerPort{
				Protocol:      val.Protocol,
				ContainerPort: val.Internal,
			})
		}
		m = true
	}

	if spec.EnvVars != nil {
		c.EnvVars = types.SpecTemplateContainerEnvs{}

		for _, e := range *spec.EnvVars {
			match := strings.Split(e, "=")
			env := types.SpecTemplateContainerEnv{Name: match[0]}
			if len(match) == 2 {
				env.Value = match[1]
			}
			c.EnvVars = append(c.EnvVars, env)
		}
		m = true
	}

	if spec.Memory != nil {
		c.Resources.Limits.RAM = *spec.Memory
		m = true
	}

	if !f {
		tpl.Template.Containers = append(tpl.Template.Containers, c)
		m = true
	} else {
		tpl.Template.Containers[i] = c
	}

	return m
}

func NewServiceModel(ctx context.Context, stg storage.Storage) IService {
	return &Service{ctx, stg}
}
