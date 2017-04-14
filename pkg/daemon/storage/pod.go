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

package storage

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/storage/store"
)

const podStorage = "pods"

// Namespace Service type for interface in interfaces folder
type PodStorage struct {
	IPod
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *PodStorage) GetByID(ctx context.Context, namespace, service, id string) (*types.Pod, error) {

	var pod = new(types.Pod)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	var ns string
	keyNamespace := s.util.Key(ctx, "helper", namespaceStorage, namespace)
	if err := client.Get(ctx, keyNamespace, &ns); err != nil {
		return nil, err
	}

	keyMeta := s.util.Key(ctx, namespaceStorage, ns, serviceStorage, service, podStorage, id)
	if err := client.Get(ctx, keyMeta, &pod); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return pod, nil
}

func (s *PodStorage) ListByService(ctx context.Context, namespace, service string) ([]*types.Pod, error) {
	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyServiceList := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service, podStorage)
	pods := []*types.Pod{}
	if err := client.List(ctx, keyServiceList, "", &pods); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	if pods == nil {
		return nil, nil
	}

	return pods, nil
}

func (s *PodStorage) Insert(ctx context.Context, namespace, service string, pod *types.Pod) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service, podStorage, pod.Meta.ID)
	if err := tx.Create(keyMeta, pod, 0); err != nil {
		return err
	}

	keyHelper := s.util.Key(ctx, "helper", serviceStorage, "pods", pod.Meta.ID)
	if err := tx.Create(keyHelper, &keyMeta, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (s *PodStorage) Update(ctx context.Context, namespace, service string, pod *types.Pod) error {
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	var ns string
	keyNamespace := s.util.Key(ctx, "helper", namespaceStorage, namespace)
	if err := client.Get(ctx, keyNamespace, &ns); err != nil {
		return err
	}

	tx := client.Begin(ctx)

	keyMeta := s.util.Key(ctx, namespaceStorage, ns, serviceStorage, service, podStorage, pod.Meta.ID)
	if err := tx.Update(keyMeta, pod, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *PodStorage) Remove(ctx context.Context, namespace, service string, pod *types.Pod) error {
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	var ns string
	keyNamespace := s.util.Key(ctx, "helper", namespaceStorage, namespace)
	if err := client.Get(ctx, keyNamespace, &ns); err != nil {
		return err
	}

	tx := client.Begin(ctx)

	keyMeta := s.util.Key(ctx, namespaceStorage, ns, serviceStorage, service, podStorage, pod.Meta.ID)
	tx.Delete(keyMeta)

	keyHelper := s.util.Key(ctx, "helper", serviceStorage, "pods", pod.Meta.ID)
	tx.Delete(keyHelper)

	KeyNodePod := s.util.Key(ctx, nodeStorage, pod.Meta.Hostname, "spec", "pods", pod.Meta.ID)
	tx.Delete(KeyNodePod)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func newPodStorage(config store.Config, util IUtil) *PodStorage {
	s := new(PodStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
