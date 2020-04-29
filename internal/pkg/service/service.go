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

package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logServicePrefix = "distribution:service"
)

type Service struct {
	context context.Context
	storage storage.IStorage
}

func (s *Service) Runtime() (*models.System, error) {

	log.Debugf("%s:get:> get services runtime info", logServicePrefix)
	runtime, err := s.storage.Info(s.context, s.storage.Collection().Service(), "")
	if err != nil {
		log.Errorf("%s:get:> get runtime info error: %s", logServicePrefix, err)
		return &runtime.System, err
	}
	return &runtime.System, nil

}

// Get service by namespace and service name
func (s *Service) Get(namespace, name string) (*models.Service, error) {

	log.Debugf("%s:get:> get in namespace %s by name %s", logServicePrefix, namespace, name)

	svc := new(models.Service)
	sl := models.NewServiceSelfLink(namespace, name).String()

	err := s.storage.Get(s.context, s.storage.Collection().Service(), sl, svc, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.Warnf("%s:get:> get in namespace %s by name %s not found", logServicePrefix, namespace, name)
			return nil, nil
		}

		log.Errorf("%s:get:> get in namespace %s by name %s error: %v", logServicePrefix, namespace, name, err)
		return nil, err
	}

	return svc, nil
}

// List method return map of services in selected namespace
func (s *Service) List(namespace string) (*models.ServiceList, error) {

	log.Debugf("%s:list:> by namespace %s", logServicePrefix, namespace)

	list := models.NewServiceList()
	q := s.storage.Filter().Service().ByNamespace(namespace)

	err := s.storage.List(s.context, s.storage.Collection().Service(), q, list, nil)
	if err != nil {
		log.Error("%s:list:> by namespace %s err: %v", logServicePrefix, namespace, err)
		return nil, err
	}

	log.Debugf("%s:list:> by namespace %s result: %d", logServicePrefix, namespace, len(list.Items))

	return list, nil
}

// Create new service model in namespace
func (s *Service) Create(namespace *models.Namespace, svc *models.Service) (*models.Service, error) {

	log.Debugf("%s:create:> create new service %#v", logServicePrefix, svc.Meta)

	svc.Meta.Namespace = namespace.Meta.Name
	svc.Meta.SelfLink = *models.NewServiceSelfLink(svc.Meta.Namespace, svc.Meta.Name)

	svc.Meta.Created = time.Now()
	svc.Meta.Updated = time.Now()

	svc.Status.State = models.StateCreated

	svc.Spec.Network.Updated = time.Now()
	svc.Spec.Template.Updated = time.Now()

	if err := s.storage.Put(s.context, s.storage.Collection().Service(),
		svc.SelfLink().String(), svc, nil); err != nil {
		log.Errorf("%s:create:> insert service err: %v", logServicePrefix, err)
		return nil, err
	}

	return svc, nil
}

// Update service in namespace
func (s *Service) Update(service *models.Service) (*models.Service, error) {

	log.Debugf("%s:update:> %#v -> %#v", logServicePrefix, service)

	if err := s.storage.Set(s.context, s.storage.Collection().Service(),
		service.SelfLink().String(), service, nil); err != nil {
		log.Errorf("%s:update:> update service spec err: %v", logServicePrefix, err)
		return nil, err
	}

	return service, nil
}

// Destroy method marks service for destroy
func (s *Service) Destroy(service *models.Service) (*models.Service, error) {

	log.Debugf("%s:destroy:> destroy service %s", logServicePrefix, service.SelfLink())

	service.Status.State = models.StateDestroy
	service.Spec.State.Destroy = true

	if err := s.storage.Set(s.context, s.storage.Collection().Service(),
		service.SelfLink().String(), service, nil); err != nil {
		log.Errorf("%s:destroy:> destroy service err: %v", logServicePrefix, err)
		return nil, err
	}
	return service, nil
}

// Remove service from storage
func (s *Service) Remove(service *models.Service) error {

	log.Debugf("%s:remove:> remove service %#v", logServicePrefix, service)

	err := s.storage.Del(s.context, s.storage.Collection().Service(),
		service.SelfLink().String())
	if err != nil {
		log.Errorf("%s:remove:> remove service err: %v", logServicePrefix, err)
		return err
	}

	return nil
}

// Set state for deployment
func (s *Service) Set(service *models.Service) error {

	if service == nil {
		return errors.New(errors.ErrStructArgIsNil)
	}

	log.Debugf("%s:setstatus:> set state for service %s", logServicePrefix, service.Meta.Name)

	if err := s.storage.Set(s.context, s.storage.Collection().Service(), service.SelfLink().String(), service, nil); err != nil {
		log.Errorf("%s:setstatus:> set state for service %s err: %v", logServicePrefix, service.Meta.Name, err)
		return err
	}

	return nil
}

// Watch service changes
func (s *Service) Watch(ch chan models.ServiceEvent, rev *int64) error {

	log.Debugf("%s:watch:> watch service by spec changes", logServicePrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-s.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := models.ServiceEvent{}
				res.Action = e.Action
				res.Name = e.Name

				service := new(models.Service)

				if err := json.Unmarshal(e.Data.([]byte), service); err != nil {
					log.Errorf("%s:> parse data err: %v", logServicePrefix, err)
					continue
				}

				res.Data = service

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := s.storage.Watch(s.context, s.storage.Collection().Service(), watcher, opts); err != nil {
		return err
	}

	return nil
}

// NewServiceModel returns new service management model
func NewServiceModel(ctx context.Context, stg storage.IStorage) *Service {
	return &Service{ctx, stg}
}
