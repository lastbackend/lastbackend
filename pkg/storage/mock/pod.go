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

const podStorage = "pods"

// Pod Service type for interface in interfaces folder
type PodStorage struct {
	storage.Pod
}

func (s *PodStorage) GetByName(ctx context.Context, app, name string) (*types.Pod, error) {
	return new(types.Pod), nil
}

func (s *PodStorage) ListByNamespace(ctx context.Context, app string) ([]*types.Pod, error) {
	return make([]*types.Pod, 0), nil
}

func (s *PodStorage) ListByService(ctx context.Context, namespace, service string) ([]*types.Pod, error) {
	return make([]*types.Pod, 0), nil
}

func (s *PodStorage) Upsert(ctx context.Context, pod *types.Pod) error {
	return nil
}

func (s *PodStorage) Update(ctx context.Context, pod *types.Pod) error {
	return nil
}

func (s *PodStorage) Remove(ctx context.Context, pod *types.Pod) error {
	return nil
}

func (s *PodStorage) Watch(ctx context.Context, pod chan *types.Pod) error {
	return nil
}

func newPodStorage() *PodStorage {
	s := new(PodStorage)
	return s
}
