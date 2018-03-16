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
)

const deploymentStorage string = "deployments"

type DeploymentStorage struct {
	storage.Deployment
}

// Get deployment by name
func (s *DeploymentStorage) Get(ctx context.Context, namespace, name string) (*types.Deployment, error) {

	log.V(logLevel).Debugf("Storage: Deployment: get by name: %s in namespace: %s", name, namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("Storage: Deployment: get deployment err: %s", err.Error())
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("Storage: Deployment: get deployment err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + deploymentStorage + `\/.+\/(?:meta|state)\b`

	var (
		filterDeploymentEndpoint = `\b.+` + endpointStorage + `\/` + name + `-` + namespace + `\..+\b`
		endpoints                = make(map[string][]string)
		deployment               = new(types.Deployment)
	)

	deployment.Spec = types.DeploymentSpec{}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Deployment: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyDeployment := keyCreate(deploymentStorage, namespace, name)
	if err := client.Map(ctx, keyDeployment, filter, deployment); err != nil {
		log.V(logLevel).Errorf("Storage: Deployment: map deployment `%s` err: %s", name, err.Error())
		return nil, err
	}

	if deployment.Meta.Name == "" {
		return nil, errors.New(store.ErrEntityNotFound)
	}

	keySpec := keyCreate(deploymentStorage, namespace, name, "spec")
	if err := client.Map(ctx, keySpec, "", &deployment.Spec); err != nil && err.Error() != store.ErrEntityNotFound {
		log.V(logLevel).Errorf("Storage: Deployment: Map deployment `%s` spec err: %s", name, err.Error())
		return nil, err
	}

	keyEndpoints := keyCreate(endpointStorage)
	if err := client.Map(ctx, keyEndpoints, filterDeploymentEndpoint, endpoints); err != nil && err.Error() != store.ErrEntityNotFound {
		log.V(logLevel).Errorf("Storage: Deployment: map deployment endpoint `%s` meta err: %s", name, err.Error())
		return nil, err
	}

	return deployment, nil
}

// Update deployment state
func (s *DeploymentStorage) updateState(ctx context.Context, deployment *types.Deployment) error {

	log.V(logLevel).Debugf("Storage: Deployment: update deployment state: %#v", deployment)

	if deployment == nil {
		err := errors.New("deployment can not be nil")
		log.V(logLevel).Errorf("Storage: Deployment: update deployment state err: %s", err.Error())
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Deployment: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyState := keyCreate(deploymentStorage, deployment.Meta.Namespace, deployment.Meta.Name, "state")
	if err := client.Upsert(ctx, keyState, deployment.State, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Deployment: upsert state err: %s", err.Error())
		return err
	}

	keyDeploymentController := keyCreate(systemStorage, types.KindController, deploymentStorage, fmt.Sprintf("%s:%s", deployment.Meta.Namespace, deployment.Meta.Name))
	if err := client.Upsert(ctx, keyDeploymentController, &deployment.State, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Deployment: upsert deployments controller err: %s", err.Error())
		return err
	}

	return nil
}

func (s *DeploymentStorage) SpecWatch(ctx context.Context, deployment chan *types.Deployment) error {

	log.V(logLevel).Debug("Storage: Deployment: watch deployment by spec")

	const filter = `\b\/` + deploymentStorage + `\/(.+)\/(.+)\/spec/.+\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Deployment: create client err: %s", err.Error())
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

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			s.updateState(ctx, d)
			deployment <- d
		}

	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("Storage: Deployment: watch deployment spec err: %s", err.Error())
		return err
	}

	return nil
}

func newDeploymentStorage() *DeploymentStorage {
	s := new(DeploymentStorage)
	return s
}
