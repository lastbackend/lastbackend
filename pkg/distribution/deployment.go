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
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"strings"
)

const (
	logDeploymentPrefix = "distribution:deployment"
)

type IDeployment interface {
	Create(service *types.Service) (*types.Deployment, error)
	Get(namespace, service, name string) (*types.Deployment, error)
	ListByNamespace(namespace string) (map[string]*types.Deployment, error)
	ListByService(namespace, service string) (map[string]*types.Deployment, error)
	SetSpec(dt *types.Deployment, opts *request.DeploymentUpdateOptions) error
	SetStatus(dt *types.Deployment) error
	Cancel(dt *types.Deployment) error
	Destroy(dt *types.Deployment) error
	Remove(dt *types.Deployment) error
	Watch(dt chan *types.Deployment) error
	WatchSpec(dt chan *types.Deployment) error
}

// Deployment - distribution model
type Deployment struct {
	context context.Context
	storage storage.Storage
}

// Get deployment info by namespace service and deployment name
func (d *Deployment) Get(namespace, service, name string) (*types.Deployment, error) {

	log.Debugf("%s:get:> namespace %s and service %s by name %s", logDeploymentPrefix, namespace, service, name)

	dt, err := d.storage.Deployment().Get(d.context, namespace, service, name)
	if err != nil {

		if err.Error() == store.ErrEntityNotFound {
			log.V(logLevel).Warnf("%s:get:> in namespace %s by name %s not found", logDeploymentPrefix, namespace, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> in namespace %s by name %s error: %s", logDeploymentPrefix, namespace, name, err)
		return nil, err
	}

	return dt, nil
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

	if err := d.storage.Deployment().Insert(d.context, deployment); err != nil {
		log.Errorf("%s:create:> distribution create in service: %s err: %s", logDeploymentPrefix, service.Meta.Name, err.Error())
		return nil, err
	}

	return deployment, nil
}

// ListByService - list of deployments by service
func (d *Deployment) ListByNamespace(namespace string) (map[string]*types.Deployment, error) {

	log.Debugf("%s:listbynamespace:> in namespace: %s", namespace)

	dl, err := d.storage.Deployment().ListByNamespace(d.context, namespace)
	if err != nil {
		log.Errorf("%s:listbynamespace:> in namespace: %s err: %s", logDeploymentPrefix, namespace, err.Error())
		return nil, err
	}

	return dl, nil
}

// ListByService - list of deployments by service
func (d *Deployment) ListByService(namespace, service string) (map[string]*types.Deployment, error) {

	log.Debugf("%s:listbyservice:> in namespace: %s and service %s", logDeploymentPrefix, namespace, service)

	dl, err := d.storage.Deployment().ListByService(d.context, namespace, service)
	if err != nil {
		log.Errorf("%s:listbyservice:> in namespace: %s and service %s err: %s", logDeploymentPrefix, namespace, service, err.Error())
		return nil, err
	}

	return dl, nil
}

// Scale deployment
func (d *Deployment) SetSpec(dt *types.Deployment, opts *request.DeploymentUpdateOptions) error {

	log.Debugf("%s:setspec:> set spec for deployment %s", logDeploymentPrefix, dt.Meta.Name)

	if dt.Spec.Replicas != *opts.Replicas {
		dt.Spec.Replicas = *opts.Replicas
		if err := d.storage.Deployment().SetSpec(d.context, dt); err != nil {
			log.Errorf("%s:setspec:> set spec for deployment %s err: %s", logDeploymentPrefix, dt.Meta.Name, err.Error())
			return err
		}
	}

	return nil
}

// Set state for deployment
func (d *Deployment) SetStatus(dt *types.Deployment) error {

	log.Debugf("%s:setstatus:> set state for deployment %s", logDeploymentPrefix, dt.Meta.Name)

	if err := d.storage.Deployment().SetStatus(d.context, dt); err != nil {
		log.Errorf("%s:setstatus:> set state for deployment %s err: %s", logDeploymentPrefix, dt.Meta.Name, err.Error())
		return err
	}

	return nil
}

// Cancel deployment
func (d *Deployment) Cancel(dt *types.Deployment) error {

	log.Debugf("%s:cancel:> cancel deployment %s", logDeploymentPrefix, dt.Meta.Name)

	// mark deployment for destroy
	dt.Spec.State.Destroy = true
	if err := d.storage.Deployment().SetSpec(d.context, dt); err != nil {
		log.Debugf("%s:destroy: destroy deployment %s err: %s", logDeploymentPrefix, dt.Meta.Name, err.Error())
		return err
	}

	// mark deployment for cancel
	dt.Status.SetCancel()

	if err := d.storage.Deployment().SetStatus(d.context, dt); err != nil {
		log.Debugf("%s:cancel:> cancel deployment %s err: %s", logDeploymentPrefix, dt.Meta.Name, err.Error())
		return err
	}

	return nil
}

// Destroy deployment
func (d *Deployment) Destroy(dt *types.Deployment) error {

	log.Debugf("%s:destroy:> destroy deployment %s", logDeploymentPrefix, dt.Meta.Name)

	// mark deployment for destroy
	dt.Spec.State.Destroy = true
	if err := d.storage.Deployment().SetSpec(d.context, dt); err != nil {
		log.Debugf("%s:destroy:> destroy deployment %s err: %s", logDeploymentPrefix, dt.Meta.Name, err.Error())
		return err
	}

	dt.Status.SetDestroy()

	if err := d.storage.Deployment().SetStatus(d.context, dt); err != nil {
		log.Debugf("%s:destroy:> destroy deployment %s err: %s", logDeploymentPrefix, dt.Meta.Name, err.Error())
		return err
	}

	return nil
}

// Destroy deployment
func (d *Deployment) Remove(dt *types.Deployment) error {

	log.Debugf("%s:remove:> remove deployment %s", logDeploymentPrefix, dt.Meta.Name)
	if err := d.storage.Deployment().Remove(d.context, dt); err != nil {
		log.Debugf("%s:remove:> remove deployment %s err: %s", logDeploymentPrefix, dt.Meta.Name, err.Error())
		return err
	}

	return nil
}

// Watch deployment changes
func (d *Deployment) Watch(dt chan *types.Deployment) error {

	log.Debugf("%s:watch:> watch deployments", logDeploymentPrefix)
	if err := d.storage.Deployment().Watch(d.context, dt); err != nil {
		log.Debugf("%s:watch:> watch deployment err: %s", logDeploymentPrefix, err.Error())
		return err
	}

	return nil
}

// Watch deployment by spec changing
func (d *Deployment) WatchSpec(dt chan *types.Deployment) error {

	log.Debugf("%s:watchspec:> watch deployments by spec changes", logDeploymentPrefix)
	if err := d.storage.Deployment().WatchSpec(d.context, dt); err != nil {
		log.Debugf("%s:watchspec:> watch deployment by spec changes err: %s", logDeploymentPrefix, err.Error())
		return err
	}

	return nil
}

func NewDeploymentModel(ctx context.Context, stg storage.Storage) IDeployment {
	return &Deployment{ctx, stg}
}
