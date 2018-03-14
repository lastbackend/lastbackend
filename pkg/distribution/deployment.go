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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"strings"
)

type IDeployment interface {
	Create(service *types.Service) (*types.Deployment, error)
	Get(namespace, service, name string) (*types.Deployment, error)
	ListByNamespace(namespace string) ([]*types.Deployment, error)
	ListByService(namespace, service string) ([]*types.Deployment, error)
	Scale(dt *types.Deployment, opts types.DeploymentOptions) error
	SetState(dt *types.Deployment) error
	Cancel(dt *types.Deployment) error
	Destroy(dt *types.Deployment) error
}

// Deployment - distribution model
type Deployment struct {
	context context.Context
	storage storage.Storage
}

// Create new deployment
func (d *Deployment) Create(service *types.Service) (*types.Deployment, error) {
	log.Debug("Deployment: Create: generate deployment for service")

	deployment := new(types.Deployment)
	deployment.State.Provision = true
	deployment.Meta.Namespace = service.Meta.Namespace
	deployment.Meta.Service = service.Meta.Name
	deployment.Meta.Status = types.StateCreated
	deployment.Meta.Name = strings.Split(generator.GetUUIDV4(), "-")[4][5:]
	deployment.Meta.SelfLink = fmt.Sprintf("%s/deployment/%s", service.Meta.SelfLink, deployment.Meta.Name)
	deployment.Meta.Endpoint = strings.ToLower(fmt.Sprintf("%s/deployment/%s", service.Meta.SelfLink, deployment.Meta.Name))

	deployment.Spec = types.DeploymentSpec{
		Replicas: service.Spec.Replicas,
		Strategy: service.Spec.Strategy,
		Template: service.Spec.Template,
		Triggers: service.Spec.Triggers,
		Selector: service.Spec.Selector,
	}

	if err := d.storage.Deployment().Insert(d.context, deployment); err != nil {
		log.Errorf("Deployment: storage insert error: %s", err.Error())
		return nil, err
	}

	return deployment, nil
}

// Check deployment is in ready stage
func (d *Deployment) Get(namespace, service, name string) (*types.Deployment, error) {

	log.Debugf("Deployment: get deployment by id: %s/%s/%s", namespace, service, name)

	dt, err := d.storage.Deployment().Get(d.context, namespace, name)
	if err != nil {
		log.Errorf("Can not get deployment by id: %s", err.Error())
		return nil, err
	}

	return dt, nil
}

// ListByService - list of deployments by service
func (d *Deployment) ListByNamespace(namespace string) ([]*types.Deployment, error) {

	log.Debug("Deployment: List By Service: get deployments list")

	dl, err := d.storage.Deployment().ListByNamespace(d.context, namespace)
	if err != nil {
		log.Errorf("Can not get deployment by id: %s", err.Error())
		return nil, err
	}

	return dl, nil
}

// ListByService - list of deployments by service
func (d *Deployment) ListByService(namespace, service string) ([]*types.Deployment, error) {

	log.Debug("Deployment: List By Service: get deployments list")

	dl, err := d.storage.Deployment().ListByService(d.context, namespace, service)
	if err != nil {
		log.Errorf("Can not get deployment by id: %s", err.Error())
		return nil, err
	}

	return dl, nil
}

// Scale deployment
func (d *Deployment) Scale(dt *types.Deployment, opts types.DeploymentOptions) error {

	log.Debugf("Deployment: SetState for deployment %s", dt.Meta.Name)

	if dt.Spec.Replicas != opts.Replicas {
		dt.Spec.Replicas = opts.Replicas
		if err := d.storage.Deployment().Update(d.context, dt); err != nil {
			log.Errorf("Can not get deployment by id: %s", err.Error())
			return err
		}
	}

	return nil
}

// Set state for deployment
func (d *Deployment) SetState(dt *types.Deployment) error {

	log.Debugf("Deployment: SetState for deployment %s", dt.Meta.Name)

	if err := d.storage.Deployment().SetState(d.context, dt); err != nil {
		log.Errorf("Can not get deployment by id: %s", err.Error())
		return err
	}

	return nil
}

// Cancel deployment
func (d *Deployment) Cancel(dt *types.Deployment) error {

	log.Debugf("Deployment: Cancel deployment %s", dt.Meta.Name)

	dt.State.Active = false
	dt.State.Cancel = true
	dt.State.Provision = true
	dt.Spec.Replicas = 0

	if err := d.storage.Deployment().Update(d.context, dt); err != nil {
		log.Errorf("Can not set deployment state err: %s", err.Error())
		return err
	}

	if err := d.storage.Deployment().SetState(d.context, dt); err != nil {
		log.Errorf("Can not set deployment state state err: %s", err.Error())
		return err
	}

	return nil
}

// Destroy deployment
func (d *Deployment) Destroy(dt *types.Deployment) error {

	log.Debugf("Deployment: Destroy deployment %s", dt.Meta.Name)

	dt.State.Active = false
	dt.State.Destroy = true
	dt.State.Provision = true
	dt.Spec.Replicas = 0

	if err := d.storage.Deployment().Update(d.context, dt); err != nil {
		log.Errorf("Can not set deployment state err: %s", err.Error())
		return err
	}

	if err := d.storage.Deployment().SetState(d.context, dt); err != nil {
		log.Errorf("Can not set deployment state state err: %s", err.Error())
		return err
	}

	return nil
}

func NewDeploymentModel(ctx context.Context, stg storage.Storage) IDeployment {
	return &Deployment{ctx, stg}
}
