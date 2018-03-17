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
	"fmt"
	"strings"
)

// Pod Service type for interface in interfaces folder
type PodStorage struct {
	storage.Pod
	data map[string]*types.Pod
}

func (s *PodStorage) Get(ctx context.Context, name string) (*types.Pod, error) {
	if ns, ok := s.data[name]; ok {
		return ns, nil
	}
	return nil, errors.New(store.ErrEntityNotFound)
}

func (s *PodStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Pod, error) {
	list := make(map[string]*types.Pod, 0)

	prefix := fmt.Sprintf("%s:", namespace)
	for _, d := range s.data {

		if strings.HasPrefix(d.Meta.Name, prefix) {
			list[d.Meta.Name] = d
		}
	}

	return list, nil
}

func (s *PodStorage) ListByService(ctx context.Context, namespace, service string) (map[string]*types.Pod, error) {
	list := make(map[string]*types.Pod, 0)

	prefix := fmt.Sprintf("%s:%s:", namespace, service)

	for _, d := range s.data {
		if strings.HasPrefix(d.Meta.Name, prefix) {
			list[d.Meta.Name] = d
		}
	}

	return list, nil
}

func (s *PodStorage) ListByDeployment(ctx context.Context, namespace, service, deployment string) (map[string]*types.Pod, error) {
	list := make(map[string]*types.Pod, 0)

	prefix := fmt.Sprintf("%s:%s:%s:", namespace, service, deployment)

	for _, d := range s.data {
		if strings.HasPrefix(d.Meta.Name, prefix) {
			list[d.Meta.Name] = d
		}
	}

	return list, nil
}

// Update deployment state
func (s *PodStorage) SetState(ctx context.Context, pod *types.Pod) error {
	if err := s.checkPodExists(pod); err != nil {
		return err
	}

	s.data[pod.Meta.Name].State = pod.State
	return nil
}

func (s *PodStorage) Insert(ctx context.Context, pod *types.Pod) error {
	if err := s.checkPodArgument(pod); err != nil {
		return err
	}

	s.data[pod.Meta.Name] = pod
	return nil
}

func (s *PodStorage) Update(ctx context.Context, pod *types.Pod) error {

	if err := s.checkPodExists(pod); err != nil {
		return err
	}

	s.data[pod.Meta.Name] = pod

	return nil
}

func (s *PodStorage) Remove(ctx context.Context, pod *types.Pod) error {

	if err := s.checkPodExists(pod); err != nil {
		return err
	}

	delete(s.data, pod.Meta.Name)

	return nil
}

func (s *PodStorage) Watch(ctx context.Context, pod chan *types.Pod) error {
	return nil
}

// Watch pod spec changes
func (s *PodStorage) WatchSpec(ctx context.Context, pod chan *types.Pod) error {
	return nil
}

func newPodStorage() *PodStorage {
	s := new(PodStorage)
	s.data = make(map[string]*types.Pod)
	return s
}

func (s *PodStorage) checkPodArgument(pod *types.Pod) error {
	if pod == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if pod.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

func (s *PodStorage) checkPodExists(pod *types.Pod) error {

	if err := s.checkPodArgument(pod); err != nil {
		return err
	}

	if _, ok := s.data[pod.Meta.Name]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}