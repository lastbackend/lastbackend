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
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
)

const podStorage = "pods"

// Namespace Service type for interface in interfaces folder
type PodStorage struct {
	IPod
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *PodStorage) GetByName(ctx context.Context, namespace, name string) (*types.Pod, error) {

	var pod = new(types.Pod)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyMeta := keyCreate(podStorage, namespace, name)
	if err := client.Get(ctx, keyMeta, pod); err != nil {
		return pod, err
	}

	return pod, nil
}

func (s *PodStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Pod, error) {
	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyList := keyCreate(podStorage, namespace)
	pods := make(map[string]*types.Pod)
	if err := client.Map(ctx, keyList, "", &pods); err != nil {
		return pods, err
	}

	return pods, nil
}

func (s *PodStorage) ListByService(ctx context.Context, namespace, service string) ([]*types.Pod, error) {
	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyServiceList := keyCreate(podStorage, namespace, service)
	pods := []*types.Pod{}
	if err := client.List(ctx, keyServiceList, "", &pods); err != nil {
		return pods, err
	}

	return pods, nil
}

func (s *PodStorage) Upsert(ctx context.Context, namespace string, pod *types.Pod) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	keyMeta := keyCreate(podStorage, namespace, pod.Meta.Name)
	if err := client.Upsert(ctx, keyMeta, pod, nil, 0); err != nil {
		return err
	}

	return nil

}

func (s *PodStorage) Update(ctx context.Context, namespace string, pod *types.Pod) error {
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	keyMeta := keyCreate(podStorage, namespace, pod.Meta.Name)
	if err := client.Update(ctx, keyMeta, pod, nil, 0); err != nil {
		return err
	}

	return nil
}

func (s *PodStorage) Remove(ctx context.Context, namespace string, pod *types.Pod) error {
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(podStorage, namespace, pod.Meta.Name)
	tx.Delete(keyMeta)

	KeyNodePod := keyCreate(nodeStorage, pod.Meta.Hostname, "spec", "pods", pod.Meta.Name)
	tx.Delete(KeyNodePod)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *PodStorage) Watch(ctx context.Context, pod chan *types.Pod) error {
	const filter = `\b\/` + podStorage + `\/(.+)/(.+)\b`
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(podStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if p, err := s.GetByName(ctx, keys[1], keys[2]); err == nil {
			pod <- p
		}

	}

	client.Watch(ctx, key, filter, cb)
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
