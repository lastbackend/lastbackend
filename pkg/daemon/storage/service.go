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
	"time"
)

const serviceStorage string = "services"

// Service Service type for interface in interfaces folder
type ServiceStorage struct {
	IService
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get service by id
func (s *ServiceStorage) GetByID(ctx context.Context, namespaceID, serviceID string) (*types.Service, error) {
	var service = new(types.Service)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, namespaceStorage, namespaceID, serviceStorage, serviceID, "meta")
	if err := client.Get(ctx, keyMeta, &service.Meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	keyConfig := s.util.Key(ctx, namespaceStorage, namespaceID, serviceStorage, serviceID, "config")
	if err := client.Get(ctx, keyConfig, &service.Config); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	keyPods := s.util.Key(ctx, namespaceStorage, namespaceID, serviceStorage, serviceID, "pods")
	if err := client.List(ctx, keyPods, "", &service.Pods); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return service, nil
}

func (s *ServiceStorage) GetByPodID(ctx context.Context, uuid string) (*types.Service, error) {
	var key string
	var service = new(types.Service)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}

	defer destroy()

	keyHelper := s.util.Key(ctx, "helper", serviceStorage, "pods", uuid)
	if err := client.Get(ctx, keyHelper, &key); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	keyMeta := s.util.Key(ctx, key, "meta")
	if err := client.Get(ctx, keyMeta, &service.Meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	keyConfig := s.util.Key(ctx, key, "config")
	if err := client.Get(ctx, keyConfig, &service.Config); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	keyPods := s.util.Key(ctx, key, "pods")
	if err := client.List(ctx, keyPods, "", &service.Pods); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}


	return service, nil
}

// Get service by name
func (s *ServiceStorage) GetByName(ctx context.Context, namespaceID string, name string) (*types.Service, error) {
	var id string

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, "helper", namespaceStorage, namespaceID, serviceStorage, name)
	if err := client.Get(ctx, key, &id); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	return s.GetByID(ctx, namespaceID, id)
}

// List services
func (s *ServiceStorage) ListByNamespace(ctx context.Context, namespaceID string) (*types.ServiceList, error) {

	const filter = `\b(.+)` + serviceStorage + `\/[a-z0-9-]{36}\/meta\b`

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyServiceList := s.util.Key(ctx, namespaceStorage, namespaceID, serviceStorage)
	metaList := []types.Meta{}
	if err := client.List(ctx, keyServiceList, filter, &metaList); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	if metaList == nil {
		return nil, nil
	}

	serviceList := new(types.ServiceList)
	for _, meta := range metaList {
		service := types.Service{}
		service.Meta.ID = meta.ID
		service.Meta.Name = meta.Name
		service.Meta.Description = meta.Description
		service.Meta.Labels = meta.Labels
		service.Meta.Created = meta.Created
		service.Meta.Updated = meta.Updated

		*serviceList = append(*serviceList, service)
	}

	return serviceList, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(ctx context.Context, service *types.Service) (*types.Service, error) {

	client, destroy, err := s.Client()
	if err != nil {
		return service, err
	}
	defer destroy()

	var namespace string
	keyNamespace := s.util.Key(ctx, "helper", namespaceStorage, service.Meta.Namespace)
	if err := client.Get(ctx, keyNamespace, &namespace); err != nil {
		return nil, err
	}

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", namespaceStorage, namespace, serviceStorage, service.Meta.Name)
	if err := tx.Create(keyHelper, &service.Meta.ID, 0); err != nil {
		return nil, err
	}

	keyMeta := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "meta")
	if err := tx.Create(keyMeta, &service.Meta, 0); err != nil {
		return nil, err
	}

	keyConfig := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "config")
	if err := tx.Create(keyConfig, &service.Config, 0); err != nil {
		return nil, err
	}

	for _, pod := range service.Pods {
		KeyPod := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "pods", pod.Meta.ID)
		if err := tx.Create(KeyPod, &pod, 0); err != nil {
			return nil, err
		}

		keyHelper := s.util.Key(ctx, "helper", serviceStorage, "pods", pod.Meta.ID)
		if err := tx.Create(keyHelper, s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID), 0); err != nil {
			return nil, err
		}

		KeyNodePod := s.util.Key(ctx, nodeStorage, pod.Meta.Hostname, "spec", "pods", pod.Meta.ID)
		if err := tx.Create(KeyNodePod, &types.PodNodeSpec{
			Meta:  pod.Meta,
			State: pod.State,
			Spec:  pod.Spec,
		}, 0); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return service, nil
}

// Update service in storage
func (s *ServiceStorage) Update(ctx context.Context, service *types.Service) (*types.Service, error) {

	service.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return service, err
	}
	defer destroy()

	var namespace string
	keyNamespace := s.util.Key(ctx, "helper", namespaceStorage, service.Meta.Namespace)
	if err := client.Get(ctx, keyNamespace, &namespace); err != nil {
		return nil, err
	}

	keyMeta := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "meta")
	smeta := new(types.Meta)
	if err := client.Get(ctx, keyMeta, smeta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	tx := client.Begin(ctx)

	if smeta.Name != service.Meta.Name {
		keyHelper1 := s.util.Key(ctx, "helper", namespaceStorage, namespace, serviceStorage, smeta.Name)
		tx.Delete(keyHelper1)

		keyHelper2 := s.util.Key(ctx, "helper", namespaceStorage, namespace, serviceStorage, service.Meta.Name)
		if err := tx.Create(keyHelper2, &service.Meta.ID, 0); err != nil {
			return nil, err
		}
	}

	keyMeta = s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "meta")
	if err := tx.Update(keyMeta, service.Meta, 0); err != nil {
		return nil, err
	}

	keyMeta = s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "config")
	if err := tx.Update(keyMeta, service.Config, 0); err != nil {
		return nil, err
	}

	for _, pod := range service.Pods {

		KeyPod := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "pods", pod.Meta.ID)
		if err := tx.Update(KeyPod, &pod, 0); err != nil {
			return nil, err
		}

		KeyNodePod := s.util.Key(ctx, nodeStorage, pod.Meta.Hostname, "spec", "pods", pod.Meta.ID)
		if err := tx.Update(KeyNodePod, &types.PodNodeSpec{
			Meta:  pod.Meta,
			State: pod.State,
			Spec:  pod.Spec,
		}, 0); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return service, nil
}

// Remove service model
func (s *ServiceStorage) Remove(ctx context.Context, service *types.Service) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	var namespace string
	keyNamespace := s.util.Key(ctx, "helper", namespaceStorage, service.Meta.Namespace)
	if err := client.Get(ctx, keyNamespace, &namespace); err != nil {
		return err
	}

	keyMeta := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "meta")
	meta := types.Meta{}
	if err := client.Get(ctx, keyMeta, &meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil
		}
		return err
	}

	tx := client.Begin(ctx)

	keyUUIDHelper := s.util.Key(context.Background(), "helper", serviceStorage, service.Meta.ID)
	tx.Delete(keyUUIDHelper)

	keyHelper := s.util.Key(ctx, "helper", namespaceStorage, namespace, serviceStorage, meta.Name)
	tx.Delete(keyHelper)

	for _, pod := range service.Pods {

		keyHelper := s.util.Key(ctx, "helper", serviceStorage, "pods", pod.Meta.ID)
		tx.Delete(keyHelper)

		KeyPod := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "pods", pod.Meta.ID)
		tx.Delete(KeyPod)

		KeyNodePod := s.util.Key(ctx, nodeStorage, pod.Meta.Hostname, "spec", "pods", pod.Meta.ID)
		tx.Delete(KeyNodePod)
	}

	keyService := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID)
	tx.DeleteDir(keyService)

	return tx.Commit()
}

// Remove services from namespace
func (s *ServiceStorage) RemoveByNamespace(ctx context.Context, namespace string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage)
	if err := client.DeleteDir(ctx, key); err != nil {
		return err
	}

	return nil
}

func newServiceStorage(config store.Config, util IUtil) *ServiceStorage {
	s := new(ServiceStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
