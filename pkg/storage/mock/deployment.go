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

package mock

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"strings"
	"fmt"
)

type DeploymentStorage struct {
	storage.Deployment
	data map[string]*types.Deployment
}

// Get deployment by name
func (s *DeploymentStorage) Get(ctx context.Context, namespace, service, name string) (*types.Deployment, error) {
	if ns, ok := s.data[s.keyCreate(namespace, service, name)]; ok {
		return ns, nil
	}
	return nil, errors.New(store.ErrEntityNotFound)
}

// Get deployments by namespace name
func (s *DeploymentStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Deployment, error) {
	list := make(map[string]*types.Deployment, 0)

	prefix := fmt.Sprintf("%s:", namespace)
	for _, d := range s.data {

		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// Get deployments by service name
func (s *DeploymentStorage) ListByService(ctx context.Context, namespace, service string) (map[string]*types.Deployment, error) {
	list := make(map[string]*types.Deployment, 0)

	prefix := fmt.Sprintf("%s:%s:", namespace, service)

	for _, d := range s.data {
		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// Update deployment state
func (s *DeploymentStorage) SetState(ctx context.Context, deployment *types.Deployment) error {
	if err := s.checkDeploymentExists(deployment); err != nil {
		return err
	}

	s.data[s.keyGet(deployment)].State = deployment.State
	return nil
}

// Insert new deployment
func (s *DeploymentStorage) Insert(ctx context.Context, deployment *types.Deployment) error {

	if err := s.checkDeploymentArgument(deployment); err != nil {
		return err
	}

	s.data[s.keyGet(deployment)] = deployment

	return nil
}

// Update deployment info
func (s *DeploymentStorage) Update(ctx context.Context, deployment *types.Deployment) error {

	if err := s.checkDeploymentExists(deployment); err != nil {
		return err
	}

	s.data[s.keyGet(deployment)] = deployment

	return nil
}

// Remove deployment from storage
func (s *DeploymentStorage) Remove(ctx context.Context, deployment *types.Deployment) error {

	if err := s.checkDeploymentExists(deployment); err != nil {
		return err
	}

	delete(s.data, s.keyGet(deployment))

	return nil
}

// Watch deployment changes
func (s *DeploymentStorage) Watch(ctx context.Context, deployment chan *types.Deployment) error {
	return nil
}

// Watch deployment spec changes
func (s *DeploymentStorage) WatchSpec(ctx context.Context, deployment chan *types.Deployment) error {
	return nil
}

func (s *DeploymentStorage) Clear(ctx context.Context) error {
	s.data = make(map[string]*types.Deployment)
	return nil
}

// keyCreate util function
func (s *DeploymentStorage) keyCreate (namespace, service, name string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, service, name)
}

// keyGet util function
func (s *DeploymentStorage) keyGet (d * types.Deployment) string {
	return d.SelfLink()
}

// newDeploymentStorage returns new storage
func newDeploymentStorage() *DeploymentStorage {
	s := new(DeploymentStorage)
	s.data = make(map[string]*types.Deployment)
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
func (s *DeploymentStorage) checkDeploymentExists(deployment *types.Deployment) error {

	if err := s.checkDeploymentArgument(deployment); err != nil {
		return err
	}

	if _, ok := s.data[s.keyGet(deployment)]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}