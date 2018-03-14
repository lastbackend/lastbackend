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

package mock

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
)

const deploymentStorage string = "deployments"

type DeploymentStorage struct {
	storage.Deployment
}

// Get deployment by name
func (s *DeploymentStorage) Get(ctx context.Context, namespace, name string) (*types.Deployment, error) {
	return new(types.Deployment), nil
}

// Update deployment state
func (s *DeploymentStorage) updateState(ctx context.Context, deployment *types.Deployment) error {
	return nil
}

func (s *DeploymentStorage) SpecWatch(ctx context.Context, deployment chan *types.Deployment) error {
	return nil
}

func newDeploymentStorage() *DeploymentStorage {
	s := new(DeploymentStorage)
	return s
}
