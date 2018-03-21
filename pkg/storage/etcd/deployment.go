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

package etcd

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
	"time"
)

const deploymentStorage = "deployments"

type DeploymentStorage struct {
	storage.Deployment
}

// Get deployment by name
func (s *DeploymentStorage) Get(ctx context.Context, namespace, service, name string) (*types.Deployment, error) {

	log.V(logLevel).Debugf("storage:etcd:deployment:> get by name: %s", name)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:deployment:> get by name err: %s", err.Error())
		return nil, err
	}

	if len(service) == 0 {
		err := errors.New("service can not be empty")
		log.V(logLevel).Errorf("storage:etcd:deployment:> get by name err: %s", err.Error())
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("storage:etcd:deployment:> get by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + deploymentStorage + `\/.+\/(?:meta|status|spec)\b`

	var (
		deployment = new(types.Deployment)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyDeployment := keyCreate(deploymentStorage, s.keyCreate(namespace, service, name))
	if err := client.Map(ctx, keyDeployment, filter, deployment); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> err: %s", name, err.Error())
		return nil, err
	}

	if deployment.Meta.Name == "" {
		return nil, errors.New(store.ErrEntityNotFound)
	}

	return deployment, nil
}

// Get deployments by namespace name
func (s *DeploymentStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Deployment, error) {

	log.V(logLevel).Debugf("storage:etcd:deployment:> get list by namespace: %s", namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:deployment:> get list by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + deploymentStorage + `\/.+\/(?:meta|status|spec)\b`

	var (
		deployments = make(map[string]*types.Deployment)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyDeployment := keyCreate(deploymentStorage, fmt.Sprintf("%s:", namespace))
	if err := client.MapList(ctx, keyDeployment, filter, deployments); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> err: %s", namespace, err.Error())
		return nil, err
	}

	return deployments, nil
}

// Get deployments by service name
func (s *DeploymentStorage) ListByService(ctx context.Context, namespace, service string) (map[string]*types.Deployment, error) {

	log.V(logLevel).Debugf("storage:etcd:deployment:> get list by namespace and service: %s:%s", namespace, service)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:deployment:> get list by namespace and service err: %s", err.Error())
		return nil, err
	}

	if len(service) == 0 {
		err := errors.New("service can not be empty")
		log.V(logLevel).Errorf("storage:etcd:deployment:> get list by namespace and service err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + deploymentStorage + `\/.+\/(?:meta|status|spec)\b`

	var (
		deployments = make(map[string]*types.Deployment)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:>  get list by namespace and service err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyDeployment := keyCreate(deploymentStorage, fmt.Sprintf("%s:%s:", namespace, service))
	if err := client.MapList(ctx, keyDeployment, filter, deployments); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> err: %s", namespace, err.Error())
		return nil, err
	}

	return deployments, nil
}

// Update deployment status
func (s *DeploymentStorage) SetStatus(ctx context.Context, deployment *types.Deployment) error {

	log.V(logLevel).Debugf("storage:etcd:deployment:> update deployment status: %#v", deployment)

	if err := s.checkDeploymentExists(ctx, deployment); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:>: update deployment err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(deploymentStorage, s.keyGet(deployment), "status")
	if err := client.Upsert(ctx, key, deployment.Status, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:>: update deployment err: %s", err.Error())
		return err
	}

	return nil
}

// Update deployment spec
func (s *DeploymentStorage) SetSpec(ctx context.Context, deployment *types.Deployment) error {

	log.V(logLevel).Debugf("storage:etcd:deployment:> update deployment spec: %#v", deployment)

	if err := s.checkDeploymentExists(ctx, deployment); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:>: update deployment err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(deploymentStorage, s.keyGet(deployment), "spec")
	if err := client.Upsert(ctx, key, deployment.Spec, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:>: update deployment err: %s", err.Error())
		return err
	}

	return nil
}

// Insert new deployment
func (s *DeploymentStorage) Insert(ctx context.Context, deployment *types.Deployment) error {

	log.V(logLevel).Debugf("storage:etcd:deployment:> insert deployment: %#v", deployment)

	if err := s.checkDeploymentArgument(deployment); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> insert deployment err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(deploymentStorage, s.keyGet(deployment), "meta")
	if err := tx.Create(keyMeta, deployment.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> insert deployment err: %s", err.Error())
		return err
	}

	keyStatus := keyCreate(deploymentStorage, s.keyGet(deployment), "status")
	if err := tx.Create(keyStatus, deployment.Spec, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> insert deployment err: %s", err.Error())
		return err
	}

	keySpec := keyCreate(deploymentStorage, s.keyGet(deployment), "spec")
	if err := tx.Create(keySpec, deployment.Spec, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> insert deployment err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> insert deployment err: %s", err.Error())
		return err
	}

	return nil
}

// Update deployment info
func (s *DeploymentStorage) Update(ctx context.Context, deployment *types.Deployment) error {

	if err := s.checkDeploymentExists(ctx, deployment); err != nil {
		return err
	}

	deployment.Meta.Updated = time.Now()
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> update deployment err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(deploymentStorage, s.keyGet(deployment), "meta")
	if err := client.Upsert(ctx, keyMeta, deployment.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> update deployment err: %s", err.Error())
		return err
	}

	return nil
}

// Remove deployment from storage
func (s *DeploymentStorage) Remove(ctx context.Context, deployment *types.Deployment) error {

	if err := s.checkDeploymentExists(ctx, deployment); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> remove err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(deploymentStorage, s.keyGet(deployment))
	if err := client.DeleteDir(ctx, keyMeta); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> remove deployment err: %s", err.Error())
		return err
	}

	return nil
}

// Watch deployment spec changes
func (s *DeploymentStorage) Watch(ctx context.Context, deployment chan *types.Deployment) error {

	log.V(logLevel).Debug("storage:etcd:deployment:> watch deployment")

	const filter = `\b\/` + deploymentStorage + `\/(.+):(.+):(.+)/.+\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> watch deployment err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(deploymentStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2], keys[3]); err == nil {
			deployment <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> watch deployment err: %s", err.Error())
		return err
	}

	return nil
}

// Watch deployment spec changes
func (s *DeploymentStorage) WatchSpec(ctx context.Context, deployment chan *types.Deployment) error {

	log.V(logLevel).Debug("storage:etcd:deployment:> watch deployment by spec")

	const filter = `\b\/` + deploymentStorage + `\/(.+):(.+):(.+)/spec\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> watch deployment by spec err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(deploymentStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2], keys[3]); err == nil {
			deployment <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> watch deployment by spec err: %s", err.Error())
		return err
	}

	return nil
}

// Watch deployment status changes
func (s *DeploymentStorage) WatchStatus(ctx context.Context, deployment chan *types.Deployment) error {

	log.V(logLevel).Debug("storage:etcd:deployment:> watch deployment by spec")

	const filter = `\b\/` + deploymentStorage + `\/(.+):(.+):(.+)/status\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> watch deployment by spec err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(deploymentStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2], keys[3]); err == nil {
			deployment <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> watch deployment by spec err: %s", err.Error())
		return err
	}

	return nil
}

// Clear deployment database
func (s *DeploymentStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:deployment:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, deploymentStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> clear err: %s", err.Error())
		return err
	}

	return nil
}

// keyCreate util function
func (s *DeploymentStorage) keyCreate(namespace, service, name string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, service, name)
}

// keyGet util function
func (s *DeploymentStorage) keyGet(d *types.Deployment) string {
	return s.keyCreate(d.Meta.Namespace, d.Meta.Service, d.Meta.Name)
}

func newDeploymentStorage() *DeploymentStorage {
	s := new(DeploymentStorage)
	return s
}

// checkDeploymentArgument - check if argument is valid for manipulations
func (s *DeploymentStorage) checkDeploymentArgument(deployment *types.Deployment) error {

	if deployment == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if deployment.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

// checkDeploymentArgument - check if deployment exists in store
func (s *DeploymentStorage) checkDeploymentExists(ctx context.Context, deployment *types.Deployment) error {

	if err := s.checkDeploymentArgument(deployment); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:deployment:> check deployment exists")

	if _, err := s.Get(ctx, deployment.Meta.Namespace, deployment.Meta.Service, deployment.Meta.Name); err != nil {
		log.V(logLevel).Debugf("storage:etcd:deployment:> check deployment exists err: %s", err.Error())
		return err
	}

	return nil
}
