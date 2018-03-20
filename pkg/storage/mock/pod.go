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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"strings"
)

// Pod Service type for interface in interfaces folder
type PodStorage struct {
	storage.Pod
	data map[string]*types.Pod
}

// Get pod from storage
func (s *PodStorage) Get(ctx context.Context, namespace, service, deployment, name string) (*types.Pod, error) {
	if ns, ok := s.data[s.keyCreate(namespace, service, deployment, name)]; ok {
		return ns, nil
	}
	return nil, errors.New(store.ErrEntityNotFound)
}

// ListByNamespace returns pod list from storage by namespace
func (s *PodStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Pod, error) {
	list := make(map[string]*types.Pod, 0)

	prefix := fmt.Sprintf("%s:", namespace)
	for _, d := range s.data {

		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// ListByService returns pod list from storage by namespace and service names
func (s *PodStorage) ListByService(ctx context.Context, namespace, service string) (map[string]*types.Pod, error) {
	list := make(map[string]*types.Pod, 0)

	prefix := fmt.Sprintf("%s:%s:", namespace, service)

	for _, d := range s.data {
		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// ListByDeployment returns pod list from storage by namespace, service and deployment names
func (s *PodStorage) ListByDeployment(ctx context.Context, namespace, service, deployment string) (map[string]*types.Pod, error) {
	list := make(map[string]*types.Pod, 0)

	prefix := fmt.Sprintf("%s:%s:%s:", namespace, service, deployment)

	for _, d := range s.data {
		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// Update deployment spec
func (s *PodStorage) SetSpec(ctx context.Context, pod *types.Pod) error {
	if err := s.checkPodExists(pod); err != nil {
		return err
	}

	s.data[s.keyGet(pod)].Spec = pod.Spec
	return nil
}

// Update deployment state
func (s *PodStorage) SetState(ctx context.Context, pod *types.Pod) error {
	if err := s.checkPodExists(pod); err != nil {
		return err
	}

	s.data[s.keyGet(pod)].State = pod.State
	return nil
}

// Insert new pod into storage
func (s *PodStorage) Insert(ctx context.Context, pod *types.Pod) error {
	if err := s.checkPodArgument(pod); err != nil {
		return err
	}

	s.data[s.keyGet(pod)] = pod
	return nil
}

// Update pod in storage
func (s *PodStorage) Update(ctx context.Context, pod *types.Pod) error {

	if err := s.checkPodExists(pod); err != nil {
		return err
	}

	s.data[s.keyGet(pod)] = pod

	return nil
}

// Remove pod from storage
func (s *PodStorage) Destroy(ctx context.Context, pod *types.Pod) error {

	if err := s.checkPodExists(pod); err != nil {
		return err
	}

	delete(s.data, s.keyGet(pod))

	return nil
}

// Remove pod from storage
func (s *PodStorage) Remove(ctx context.Context, pod *types.Pod) error {

	if err := s.checkPodExists(pod); err != nil {
		return err
	}

	delete(s.data, s.keyGet(pod))

	return nil
}

// Watch pod changes
func (s *PodStorage) Watch(ctx context.Context, pod chan *types.Pod) error {
	return nil
}

// Watch pod spec changes
func (s *PodStorage) WatchSpec(ctx context.Context, pod chan *types.Pod) error {
	return nil
}

// Watch pod state changes
func (s *PodStorage) WatchState(ctx context.Context, pod chan *types.Pod) error {
	return nil
}

// Watch pod status changes
func (s *PodStorage) WatchStatus(ctx context.Context, pod chan *types.Pod) error {
	return nil
}

// Clear pod storage
func (s *PodStorage) Clear(ctx context.Context) error {
	s.data = make(map[string]*types.Pod)
	return nil
}

// keyCreate util function
func (s *PodStorage) keyCreate(namespace, service, deployment, name string) string {
	return fmt.Sprintf("%s:%s:%s:%s", namespace, service, deployment, name)
}

// keyGet util function
func (s *PodStorage) keyGet(p *types.Pod) string {
	return p.SelfLink()
}

// newPodStorage returns new podStorage
func newPodStorage() *PodStorage {
	s := new(PodStorage)
	s.data = make(map[string]*types.Pod)
	return s
}

// checkPodArgument method checks pod argument
func (s *PodStorage) checkPodArgument(pod *types.Pod) error {
	if pod == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if pod.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

// checkPodArgument method checks if pod exists in storage
func (s *PodStorage) checkPodExists(pod *types.Pod) error {

	if err := s.checkPodArgument(pod); err != nil {
		return err
	}

	if _, ok := s.data[s.keyGet(pod)]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}
