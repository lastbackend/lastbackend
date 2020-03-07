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

package model

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"time"

	"encoding/json"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logDeploymentPrefix = "distribution:deployment"
)

// Deployment - distribution model
type Deployment struct {
	context context.Context
	storage storage.Storage
}

func (d *Deployment) Runtime() (*types.System, error) {

	log.V(logLevel).Debugf("%s:get:> get deployment runtime info", logDeploymentPrefix)
	runtime, err := d.storage.Info(d.context, d.storage.Collection().Deployment(), "")
	if err != nil {
		log.V(logLevel).Errorf("%s:get:> get runtime info error: %s", logDeploymentPrefix, err)
		return &runtime.System, err
	}
	return &runtime.System, nil
}

// Get deployment info by namespace service and deployment name
func (d *Deployment) Get(namespace, service, name string) (*types.Deployment, error) {

	log.V(logLevel).Debugf("%s:get:> namespace %s and service %s by name %s", logDeploymentPrefix, namespace, service, name)

	dp := new(types.Deployment)

	err := d.storage.Get(d.context, d.storage.Collection().Deployment(), types.NewDeploymentSelfLink(namespace, service, name).String(), &dp, nil)
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
func (d *Deployment) Create(service *types.Service, name string) (*types.Deployment, error) {

	log.V(logLevel).Debugf("%s:create:> distribution create in service: %s", logDeploymentPrefix, service.Meta.Name)

	deployment := new(types.Deployment)

	deployment.Meta.Namespace = service.Meta.Namespace
	deployment.Meta.Service = service.Meta.Name
	deployment.Meta.Name = name
	deployment.Meta.Created = time.Now()
	deployment.Meta.Updated = time.Now()

	deployment.Meta.Labels = make(map[string]string, 0)
	for k, v := range service.Meta.Labels {
		deployment.Meta.Labels[k] = v
	}

	deployment.Meta.SelfLink = *types.NewDeploymentSelfLink(deployment.Meta.Namespace, deployment.Meta.Service, deployment.Meta.Name)

	deployment.Spec = types.DeploymentSpec{
		Replicas: service.Spec.Replicas,
		Template: service.Spec.Template,
		Selector: service.Spec.Selector,
	}

	deployment.Status.SetCreated()

	if err := d.storage.Put(d.context, d.storage.Collection().Deployment(),
		deployment.SelfLink().String(), deployment, nil); err != nil {
		log.Errorf("%s:create:> distribution create in service: %s err: %v", logDeploymentPrefix, service.Meta.Name, err)
		return nil, err
	}

	return deployment, nil
}

func (d *Deployment) Insert(dep *types.Deployment) error {

	log.V(logLevel).Debugf("%s:create:> distribution insert new deployment in service: %s", logDeploymentPrefix, dep.Meta.Service)

	if err := d.storage.Put(d.context, d.storage.Collection().Deployment(),
		dep.SelfLink().String(), dep, nil); err != nil {
		log.Errorf("%s:create:> distribution create deployment in service: %s err: %v", logDeploymentPrefix, dep.Meta.Service, err)
		return err
	}

	return nil
}

// ListByService - list of deployments by service
func (d *Deployment) ListByNamespace(namespace string) (*types.DeploymentList, error) {

	log.V(logLevel).Debugf("%s:listbynamespace:> in namespace: %s", namespace)

	q := d.storage.Filter().Deployment().ByNamespace(namespace)
	dl := types.NewDeploymentList()

	err := d.storage.List(d.context, d.storage.Collection().Deployment(), q, dl, nil)
	if err != nil {
		log.Errorf("%s:listbynamespace:> in namespace: %s err: %v", logDeploymentPrefix, namespace, err)
		return nil, err
	}

	return dl, nil
}

// ListByService - list of deployments by service
func (d *Deployment) ListByService(namespace, service string) (*types.DeploymentList, error) {

	log.V(logLevel).Debugf("%s:listbyservice:> in namespace: %s and service %s", logDeploymentPrefix, namespace, service)

	q := d.storage.Filter().Deployment().ByService(namespace, service)
	dl := types.NewDeploymentList()

	err := d.storage.List(d.context, d.storage.Collection().Deployment(), q, dl, nil)
	if err != nil {
		log.Errorf("%s:listbyservice:> in namespace: %s and service %s err: %v", logDeploymentPrefix, namespace, service, err)
		return nil, err
	}

	return dl, nil
}

// Update deployment
func (d *Deployment) Update(dt *types.Deployment) error {

	log.V(logLevel).Debugf("%s:update:> update deployment %s", logDeploymentPrefix, dt.Meta.Name)

	if err := d.storage.Set(d.context, d.storage.Collection().Deployment(),
		dt.SelfLink().String(), dt, nil); err != nil {
		log.Errorf("%s:update:> update for deployment %s err: %v", logDeploymentPrefix, dt.Meta.Name, err)
		return err
	}

	return nil
}

// Cancel deployment
func (d *Deployment) Cancel(dt *types.Deployment) error {

	log.V(logLevel).Debugf("%s:cancel:> cancel deployment %s", logDeploymentPrefix, dt.Meta.Name)

	// mark deployment for destroy
	dt.Spec.State.Destroy = true
	// mark deployment for cancel
	dt.Status.SetCancel()

	if err := d.storage.Set(d.context, d.storage.Collection().Deployment(),
		dt.SelfLink().String(), dt, nil); err != nil {
		log.V(logLevel).Debugf("%s:destroy: destroy deployment %s err: %v", logDeploymentPrefix, dt.Meta.Name, err)
		return err
	}

	return nil
}

// Destroy deployment
func (d *Deployment) Destroy(dt *types.Deployment) error {

	log.V(logLevel).Debugf("%s:destroy:> destroy deployment %s", logDeploymentPrefix, dt.Meta.Name)

	// mark deployment for destroy
	dt.Spec.State.Destroy = true
	// mark deployment for destroy
	dt.Status.SetDestroy()

	if err := d.storage.Set(d.context, d.storage.Collection().Deployment(),
		dt.SelfLink().String(), dt, nil); err != nil {
		log.V(logLevel).Debugf("%s:destroy:> destroy deployment %s err: %v", logDeploymentPrefix, dt.Meta.Name, err)
		return err
	}

	return nil
}

// Destroy deployment
func (d *Deployment) Remove(dt *types.Deployment) error {

	log.V(logLevel).Debugf("%s:remove:> remove deployment %s", logDeploymentPrefix, dt.Meta.Name)
	if err := d.storage.Del(d.context, d.storage.Collection().Deployment(),
		dt.SelfLink().String()); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove deployment %s err: %v", logDeploymentPrefix, dt.Meta.Name, err)
		return err
	}

	return nil
}

// Watch deployment changes
func (d *Deployment) Watch(dt chan types.DeploymentEvent, rev *int64) error {

	done := make(chan bool)
	watcher := storage.NewWatcher()

	log.V(logLevel).Debugf("%s:watch:> watch deployments", logDeploymentPrefix)

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

				res := types.DeploymentEvent{}
				res.Action = e.Action
				res.Name = e.Name

				deployment := new(types.Deployment)

				if err := json.Unmarshal(e.Data.([]byte), deployment); err != nil {
					log.Errorf("%s:> parse data err: %v", logDeploymentPrefix, err)
					continue
				}

				res.Data = deployment

				dt <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := d.storage.Watch(d.context, d.storage.Collection().Deployment(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewDeploymentModel(ctx context.Context, stg storage.Storage) *Deployment {
	return &Deployment{ctx, stg}
}
