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

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3/store"
	"github.com/spf13/viper"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd"
	stgtypes "github.com/lastbackend/lastbackend/pkg/storage/types"
	"encoding/json"
)

const (
	logServicePrefix = "distribution:service"
)

type IService interface {
	Get(namespace, service string) (*types.Service, error)
	List(namespace string) (map[string]*types.Service, error)
	Create(namespace *types.Namespace, opts *types.ServiceCreateOptions) (*types.Service, error)
	Update(service *types.Service, opts *types.ServiceUpdateOptions) (*types.Service, error)
	Destroy(service *types.Service) (*types.Service, error)
	Remove(service *types.Service) error
	SetStatus(service *types.Service) error
	Watch(ch chan types.ServiceEvent) error
}

type Service struct {
	context context.Context
	storage storage.Storage
}

// Get service by namespace and service name
func (s *Service) Get(namespace, service string) (*types.Service, error) {

	log.V(logLevel).Debugf("%s:get:> get in namespace %s by name %s", logServicePrefix, namespace, service)

	svc := new(types.Service)

	err := s.storage.Get(s.context, storage.ServiceKind, etcd.BuildServiceKey(namespace, service), &svc)
	if err != nil {

		if err.Error() == store.ErrEntityNotFound {
			log.V(logLevel).Warnf("%s:get:> get in namespace %s by name %s not found", logServicePrefix, namespace, service)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> get in namespace %s by name %s error: %v", logServicePrefix, namespace, service, err)
		return nil, err
	}

	return svc, nil
}

// List method return map of services in selected namespace
func (s *Service) List(namespace string) (map[string]*types.Service, error) {

	log.V(logLevel).Debugf("%s:list:> by namespace %s", logServicePrefix, namespace)

	items := make(map[string]*types.Service, 0)
	q := etcd.BuildServiceQuery(namespace)

	err := s.storage.Map(s.context, storage.ServiceKind, q, &items)
	if err != nil {
		log.V(logLevel).Error("%s:list:> by namespace %s err: %v", logServicePrefix, namespace, err)
		return nil, err
	}

	log.V(logLevel).Debugf("%s:list:> by namespace %s result: %d", logServicePrefix, namespace, len(items))

	return items, nil
}

// Create new service model in namespace
func (s *Service) Create(namespace *types.Namespace, opts *types.ServiceCreateOptions) (*types.Service, error) {

	log.V(logLevel).Debugf("%s:create:> create new service %#v", logServicePrefix, opts)

	service := new(types.Service)
	switch true {
	case opts == nil:
		return nil, errors.New("opts can not be nil")
	case opts.Name == nil || *opts.Name == "":
		return nil, errors.New("name is required")
	case opts.Image == nil || *opts.Image == "":
		return nil, errors.New("image is required")
	}

	// prepare meta data for service
	service.Meta.SetDefault()
	service.Meta.Name = *opts.Name
	service.Meta.Namespace = namespace.Meta.Name
	service.Meta.Endpoint = strings.ToLower(fmt.Sprintf("%s-%s.%s", *opts.Name, namespace.Meta.Name, viper.GetString("domain.internal")))

	if opts.Description != nil {
		service.Meta.Description = *opts.Description
	}

	service.SelfLink()

	// prepare default template spec
	c := types.SpecTemplateContainer{}
	c.SetDefault()
	c.Role = types.ContainerRolePrimary
	c.Image.Name = *opts.Image

	// prepare spec data for service
	service.Spec.SetDefault()
	service.Spec.Template.Containers = append(service.Spec.Template.Containers, c)

	if opts.Spec != nil {
		service.Spec.Update(opts.Spec)
	}

	service.SelfLink()

	if err := s.storage.Create(s.context, storage.ServiceKind, service.Meta.SelfLink, service, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert service err: %v", logServicePrefix, err)
		return nil, err
	}

	return service, nil
}

// Update service in namespace
func (s *Service) Update(service *types.Service, opts *types.ServiceUpdateOptions) (*types.Service, error) {

	log.V(logLevel).Debugf("%s:update:> %#v -> %#v", logServicePrefix, service, opts)

	if opts == nil {
		opts = new(types.ServiceUpdateOptions)
	}

	if opts.Description != nil {
		service.Meta.Description = *opts.Description

		if err := s.storage.Update(s.context, storage.ServiceKind, service.Meta.SelfLink, service, nil); err != nil {
			log.V(logLevel).Errorf("%s:update:> update service meta err: %v", logServicePrefix, err)
			return nil, err
		}
	}

	if opts.Spec != nil {

		service.Spec.Update(opts.Spec)

		if err := s.storage.Update(s.context, storage.ServiceKind, service.Meta.SelfLink, service, nil); err != nil {
			log.V(logLevel).Errorf("%s:update:> update service spec err: %v", logServicePrefix, err)
			return nil, err
		}
	}

	return service, nil
}

// Destroy method marks service for destroy
func (s *Service) Destroy(service *types.Service) (*types.Service, error) {

	log.V(logLevel).Debugf("%s:destroy:> destroy service %s", logServicePrefix, service.SelfLink())

	service.Status.State = types.StateDestroy
	service.Spec.State.Destroy = true

	if err := s.storage.Update(s.context, storage.ServiceKind, service.Meta.SelfLink, service, nil); err != nil {
		log.V(logLevel).Errorf("%s:destroy:> destroy service err: %v", logServicePrefix, err)
		return nil, err
	}
	return service, nil
}

// Remove service from storage
func (s *Service) Remove(service *types.Service) error {

	log.V(logLevel).Debugf("%s:remove:> remove service %#v", logServicePrefix, service)

	err := s.storage.Remove(s.context, storage.ServiceKind, service.Meta.SelfLink)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove service err: %v", logServicePrefix, err)
		return err
	}

	return nil
}

// Set state for deployment
func (s *Service) SetStatus(service *types.Service) error {

	log.Debugf("%s:setstatus:> set state for service %s", logServicePrefix, service.Meta.Name)

	if err := s.storage.Update(s.context, storage.ServiceKind, service.Meta.SelfLink, service, nil); err != nil {
		log.Errorf("%s:setstatus:> set state for service %s err: %v", logServicePrefix, service.Meta.Name, err)
		return err
	}

	return nil
}

// Watch service changes
func (s *Service) Watch(ch chan types.ServiceEvent) error {

	log.Debugf("%s:watch:> watch service by spec changes", logServicePrefix)

	done := make(chan bool)
	event := make(chan *stgtypes.WatcherEvent)

	go func() {
		for {
			select {
			case <-s.context.Done():
				done <- true
				return
			case e := <-event:
				if e.Data == nil {
					continue
				}

				res := types.ServiceEvent{}
				res.Action = e.Action
				res.Name = e.Name

				service := new(types.Service)

				if err := json.Unmarshal(e.Data.([]byte), *service); err != nil {
					log.Errorf("%s:> parse data err: %v", logServicePrefix, err)
					continue
				}

				res.Data = service

				ch <- res
			}
		}
	}()

	if err := s.storage.Watch(s.context, storage.ServiceKind, event); err != nil {
		return err
	}

	return nil
}

// NewServiceModel returns new service management model
func NewServiceModel(ctx context.Context, stg storage.Storage) IService {
	return &Service{ctx, stg}
}
