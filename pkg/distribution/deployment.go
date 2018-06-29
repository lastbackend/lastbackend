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
	"strings"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/generator"

	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

const (
	logDeploymentPrefix = "distribution:deployment"
)

type IDeployment interface {
	Create(service *types.Service) (*types.Deployment, error)
	Get(namespace, service, name string) (*types.Deployment, error)
	ListByNamespace(namespace string) ([]*types.Deployment, error)
	ListByService(namespace, service string) ([]*types.Deployment, error)
	Update(dt *types.Deployment, opts *types.DeploymentUpdateOptions) error
	Cancel(dt *types.Deployment) error
	Destroy(dt *types.Deployment) error
	Remove(dt *types.Deployment) error
	Watch(dt chan *types.Deployment)
}

// Deployment - distribution model
type Deployment struct {
	context context.Context
	storage storage.Storage
}

// Get deployment info by namespace service and deployment name
func (d *Deployment) Get(namespace, service, name string) (*types.Deployment, error) {

	log.Debugf("%s:get:> namespace %s and service %s by name %s", logDeploymentPrefix, namespace, service, name)

	dp := new(types.Deployment)

	err := d.storage.Get(d.context, storage.DeploymentKind, d.storage.Key().Deployment(namespace, service, name), &dp)
	if err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> in namespace %s by name %s not found", logDeploymentPrefix, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> in namespace %s by name %s error: %s", logDeploymentPrefix, namespace, name, err)
		return nil, err
	}

	return dp, nil
}

// Create new deployment
func (d *Deployment) Create(service *types.Service) (*types.Deployment, error) {

	log.Debugf("%s:create:> distribution create in service: %s", logDeploymentPrefix, service.Meta.Name)

	deployment := new(types.Deployment)

	deployment.Meta.Namespace = service.Meta.Namespace
	deployment.Meta.Service = service.Meta.Name
	deployment.Meta.Status = types.StateCreated
	deployment.Meta.Name = strings.Split(generator.GetUUIDV4(), "-")[4][5:]

	deployment.SelfLink()

	deployment.Spec = types.DeploymentSpec{
		Replicas: service.Spec.Replicas,
		Strategy: service.Spec.Strategy,
		Template: service.Spec.Template,
		Triggers: service.Spec.Triggers,
		Selector: service.Spec.Selector,
	}

	deployment.Spec.Meta.SetDefault()
	deployment.Spec.Meta.Name = service.Spec.Meta.Name

	deployment.Status.SetProvision()

	if err := d.storage.Put(d.context, storage.DeploymentKind,
		d.storage.Key().Deployment(deployment.Meta.Namespace, deployment.Meta.Service, deployment.Meta.Name), deployment, nil); err != nil {
		log.Errorf("%s:create:> distribution create in service: %s err: %v", logDeploymentPrefix, service.Meta.Name, err)
		return nil, err
	}

	return deployment, nil
}

// ListByService - list of deployments by service
func (d *Deployment) ListByNamespace(namespace string) ([]*types.Deployment, error) {

	log.Debugf("%s:listbynamespace:> in namespace: %s", namespace)

	q := d.storage.Filter().Deployment().ByNamespace(namespace)
	dl := make([]*types.Deployment, 0)

	err := d.storage.List(d.context, storage.DeploymentKind, q, &dl)
	if err != nil {
		log.Errorf("%s:listbynamespace:> in namespace: %s err: %v", logDeploymentPrefix, namespace, err)
		return nil, err
	}

	return dl, nil
}

// ListByService - list of deployments by service
func (d *Deployment) ListByService(namespace, service string) ([]*types.Deployment, error) {

	log.Debugf("%s:listbyservice:> in namespace: %s and service %s", logDeploymentPrefix, namespace, service)

	q := d.storage.Filter().Deployment().ByService(namespace, service)
	dl := make([]*types.Deployment, 0)

	err := d.storage.List(d.context, storage.DeploymentKind, q, &dl)
	if err != nil {
		log.Errorf("%s:listbyservice:> in namespace: %s and service %s err: %v", logDeploymentPrefix, namespace, service, err)
		return nil, err
	}

	return dl, nil
}

// Update deployment
func (d *Deployment) Update(dt *types.Deployment, opts *types.DeploymentUpdateOptions) error {

	log.Debugf("%s:update:> update deployment %s", logDeploymentPrefix, dt.Meta.Name)

	var isChanged = false

	switch true {
	case opts.Replicas != nil && dt.Spec.Replicas != *opts.Replicas:
		dt.Spec.Replicas = *opts.Replicas
		isChanged = true
		break
	case opts.Status != nil:
		dt.Status.State = opts.Status.State
		dt.Status.Message = opts.Status.Message
		isChanged = true
		break
	}

	if isChanged {
		if err := d.storage.Set(d.context, storage.DeploymentKind,
			d.storage.Key().Deployment(dt.Meta.Namespace, dt.Meta.Service, dt.Meta.Name), dt, nil); err != nil {
			log.Errorf("%s:update:> update for deployment %s err: %v", logDeploymentPrefix, dt.Meta.Name, err)
			return err
		}
	}

	return nil
}

// Cancel deployment
func (d *Deployment) Cancel(dt *types.Deployment) error {

	log.Debugf("%s:cancel:> cancel deployment %s", logDeploymentPrefix, dt.Meta.Name)

	// mark deployment for destroy
	dt.Spec.State.Destroy = true
	// mark deployment for cancel
	dt.Status.SetCancel()

	if err := d.storage.Set(d.context, storage.DeploymentKind,
		d.storage.Key().Deployment(dt.Meta.Namespace, dt.Meta.Service, dt.Meta.Name), dt, nil); err != nil {
		log.Debugf("%s:destroy: destroy deployment %s err: %v", logDeploymentPrefix, dt.Meta.Name, err)
		return err
	}

	return nil
}

// Destroy deployment
func (d *Deployment) Destroy(dt *types.Deployment) error {

	log.Debugf("%s:destroy:> destroy deployment %s", logDeploymentPrefix, dt.Meta.Name)

	// mark deployment for destroy
	dt.Spec.State.Destroy = true
	// mark deployment for destroy
	dt.Status.SetDestroy()

	if err := d.storage.Set(d.context, storage.DeploymentKind,
		d.storage.Key().Deployment(dt.Meta.Namespace, dt.Meta.Service, dt.Meta.Name), dt, nil); err != nil {
		log.Debugf("%s:destroy:> destroy deployment %s err: %v", logDeploymentPrefix, dt.Meta.Name, err)
		return err
	}

	return nil
}

// Destroy deployment
func (d *Deployment) Remove(dt *types.Deployment) error {

	log.Debugf("%s:remove:> remove deployment %s", logDeploymentPrefix, dt.Meta.Name)
	if err := d.storage.Del(d.context, storage.DeploymentKind,
		d.storage.Key().Deployment(dt.Meta.Namespace, dt.Meta.Service, dt.Meta.Name)); err != nil {
		log.Debugf("%s:remove:> remove deployment %s err: %v", logDeploymentPrefix, dt.Meta.Name, err)
		return err
	}

	return nil
}

// Watch deployment changes
func (d *Deployment) Watch(dt chan *types.Deployment) {

	done := make(chan bool)
	watcher := storage.NewWatcher()

	log.Debugf("%s:watch:> watch deployments", logDeploymentPrefix)

	go func() {
		for {
			select {
			case <-d.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				deployment := new(types.Deployment)

				if err := json.Unmarshal(e.Data.([]byte), *deployment); err != nil {
					log.Errorf("%s:> parse data err: %v", logDeploymentPrefix, err)
					continue
				}

				dt <- deployment
			}
		}
	}()

	go d.storage.Watch(d.context, storage.DeploymentKind, watcher)

	<-done
}

func NewDeploymentModel(ctx context.Context, stg storage.Storage) IDeployment {
	return &Deployment{ctx, stg}
}
