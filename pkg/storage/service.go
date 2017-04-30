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
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
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

	const filter = `\b(.+)` + serviceStorage + `\/[a-z0-9-]{36}\/(meta|state)\b`
	var service = new(types.Service)
	service.Spec = make(map[string]*types.ServiceSpec)
	service.Pods = make(map[string]*types.Pod)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, namespaceStorage, namespaceID, serviceStorage, serviceID)
	if err := client.Map(ctx, key, filter, service); err != nil {
		return nil, err
	}

	keySpec := s.util.Key(ctx, namespaceStorage, namespaceID, serviceStorage, serviceID, "specs")
	if err := client.Map(ctx, keySpec, "", &service.Spec); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return service, nil
		}
		return nil, err
	}

	keyPods := s.util.Key(ctx, namespaceStorage, namespaceID, serviceStorage, serviceID, "pods")
	if err := client.Map(ctx, keyPods, "", &service.Pods); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return service, nil
		}
		return nil, err
	}

	return service, nil
}

// Get service by Pod ID
func (s *ServiceStorage) GetByPodID(ctx context.Context, uuid string) (*types.Service, error) {
	const filter = `\b(.+)` + serviceStorage + `\/[a-z0-9-]{36}\/(meta)\b`

	var key string
	var service = new(types.Service)
	service.Pods = make(map[string]*types.Pod)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}

	defer destroy()

	keyHelper := s.util.Key(ctx, "helper", serviceStorage, "pods", uuid)
	if err := client.Get(ctx, keyHelper, &key); err != nil {
		return nil, err
	}

	if err := client.Map(ctx, key, filter, service); err != nil {
		return nil, err
	}

	keyPods := s.util.Key(ctx, key, "pods")
	if err := client.Map(ctx, keyPods, "", &service.Pods); err != nil {
		return nil, err
	}

	return service, nil
}

// Get service by name
func (s *ServiceStorage) GetByName(ctx context.Context, namespaceID string, name string) (*types.Service, error) {

	var (
		id string
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, "helper", namespaceStorage, namespaceID, serviceStorage, name)
	if err := client.Get(ctx, key, &id); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, namespaceID, id)
}

// Get service by name
func (s *ServiceStorage) UpdateState(ctx context.Context, service *types.Service) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	service.State.Resources = types.ServiceResourcesState{}

	for _,s := range service.Spec {
		service.State.Resources.Memory+= int(s.Memory)*service.Meta.Replicas
	}

	service.State.Replicas = types.ServiceReplicasState{}

	for _, p := range service.Pods {
		service.State.Replicas.Total++
		switch p.State.State {
		case types.StateCreated: service.State.Replicas.Created++
		case types.StateStarted: service.State.Replicas.Running++
		case types.StateStopped: service.State.Replicas.Stopped++
		case types.StateError: service.State.Replicas.Errored++
		}

		if p.State.Provision {
			service.State.Replicas.Provision++
		}

		if p.State.Ready {
			service.State.Replicas.Ready++
		}
	}

	var namespace string
	keyNamespace := s.util.Key(ctx, "helper", namespaceStorage, service.Meta.Namespace)
	if err := client.Get(ctx, keyNamespace, &namespace); err != nil {
		return err
	}

	tx := client.Begin(ctx)

	keyState := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "state")
	if err := tx.Upsert(keyState, service.State, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil


}

// List services
func (s *ServiceStorage) ListByNamespace(ctx context.Context, namespaceID string) ([]*types.Service, error) {

	const filter = `\b(.+)` + serviceStorage + `\/[a-z0-9-]{36}\/(meta|state)\b`

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	var services = []*types.Service{}

	keyServices := s.util.Key(ctx, namespaceStorage, namespaceID, serviceStorage)
	if err := client.List(ctx, keyServices, filter, &services); err != nil {
		return nil, err
	}

	return services, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(ctx context.Context, service *types.Service) error {

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

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", namespaceStorage, namespace, serviceStorage, service.Meta.Name)
	if err := tx.Create(keyHelper, &service.Meta.ID, 0); err != nil {
		return err
	}

	keyMeta := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "meta")
	if err := tx.Create(keyMeta, &service.Meta, 0); err != nil {
		return err
	}

	for _, spec := range service.Spec {
		keyConfig := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "specs", spec.Meta.ID)
		if err := tx.Create(keyConfig, &spec, 0); err != nil {
			return err
		}
	}

	for _, pod := range service.Pods {
		KeyPod := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "pods", pod.Meta.ID)
		if err := tx.Create(KeyPod, &pod, 0); err != nil {
			return err
		}

		keyHelper := s.util.Key(ctx, "helper", serviceStorage, "pods", pod.Meta.ID)
		if err := tx.Create(keyHelper, s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID), 0); err != nil {
			return err
		}

		KeyNodePod := s.util.Key(ctx, nodeStorage, pod.Meta.Hostname, "spec", "pods", pod.Meta.ID)
		if err := tx.Create(KeyNodePod, &types.PodNodeSpec{
			Meta:  pod.Meta,
			Spec:  pod.Spec,
			State: pod.State,
		}, 0); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Update service in storage
func (s *ServiceStorage) Update(ctx context.Context, service *types.Service) error {

	service.Meta.Updated = time.Now()

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
	smeta := new(types.Meta)
	if err := client.Get(ctx, keyMeta, smeta); err != nil {
		return err
	}

	tx := client.Begin(ctx)

	if smeta.Name != service.Meta.Name {
		keyHelper1 := s.util.Key(ctx, "helper", namespaceStorage, namespace, serviceStorage, smeta.Name)
		tx.Delete(keyHelper1)

		keyHelper2 := s.util.Key(ctx, "helper", namespaceStorage, namespace, serviceStorage, service.Meta.Name)
		if err := tx.Create(keyHelper2, &service.Meta.ID, 0); err != nil {
			return err
		}
	}

	keyMeta = s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "meta")
	if err := tx.Update(keyMeta, service.Meta, 0); err != nil {
		return err
	}

	specs := make(map[string]*types.ServiceSpec)
	keySpecs := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "specs")
	if err := client.Map(ctx, keySpecs, "", specs); err != nil {
		if err.Error() != store.ErrKeyNotFound {
			return err
		}
	}

	for id := range specs {
		if _, ok := service.Spec[id]; ok {
			continue
		}
		keySpec := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "specs", id)
		tx.DeleteDir(keySpec)
	}

	for id, spec := range service.Spec {
		keySpec := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "specs", id)
		if err := tx.Upsert(keySpec, &spec, 0); err != nil {
			return err
		}
	}

	for _, pod := range service.Pods {

		keyHelper := s.util.Key(ctx, "helper", serviceStorage, "pods", pod.Meta.ID)
		if err := tx.Upsert(keyHelper, s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID), 0); err != nil {
			return err
		}

		KeyPod := s.util.Key(ctx, namespaceStorage, namespace, serviceStorage, service.Meta.ID, "pods", pod.Meta.ID)
		if err := tx.Upsert(KeyPod, &pod, 0); err != nil {
			return err
		}

		KeyNodePod := s.util.Key(ctx, nodeStorage, pod.Meta.Hostname, "spec", "pods", pod.Meta.ID)
		if err := tx.Upsert(KeyNodePod, &types.PodNodeSpec{
			Meta:  pod.Meta,
			Spec:  pod.Spec,
			State: pod.State,
		}, 0); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
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
		return err
	}

	tx := client.Begin(ctx)

	keyUUIDHelper := s.util.Key(context.Background(), "helper", serviceStorage, service.Meta.ID)
	tx.Delete(keyUUIDHelper)

	keyHelper := s.util.Key(ctx, "helper", namespaceStorage, namespace, serviceStorage, meta.Name)
	tx.Delete(keyHelper)

	keySpec := s.util.Key(ctx, "helper", namespaceStorage, namespace, serviceStorage, meta.Name, "specs")
	tx.DeleteDir(keySpec)

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

func (s *ServiceStorage) PodsWatch(ctx context.Context, service chan *types.Service) error {
	const filter = `\b.+` + namespaceStorage + `\/([a-z0-9-]{36})\/` + serviceStorage + `\/([a-z0-9-]{36})\/pods/.+\b`
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := s.util.Key(ctx, namespaceStorage)
	cb := func(key string) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if svc, err := s.GetByID(ctx, keys[1], keys[2]); err == nil {
			s.UpdateState(ctx, svc)
			service <- svc
		}

	}

	client.Watch(ctx, key, filter, cb)
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
