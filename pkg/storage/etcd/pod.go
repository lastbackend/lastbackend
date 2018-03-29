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

package etcd

import (
	"context"
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
	"time"
)

const podStorage = "pods"

// Pod Service type for interface in interfaces folder
type PodStorage struct {
	storage.Pod
}

// Get pod from storage
func (s *PodStorage) Get(ctx context.Context, namespace, service, deployment, name string) (*types.Pod, error) {

	log.V(logLevel).Debugf("storage:etcd:pod:> get by name: %s", name)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:pod:> get by name err: %s", err.Error())
		return nil, err
	}

	if len(service) == 0 {
		err := errors.New("service can not be empty")
		log.V(logLevel).Errorf("storage:etcd:pod:> get by name err: %s", err.Error())
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("storage:etcd:pod:> get by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + podStorage + `\/.+\/(?:meta|status|spec|status)\b`
	var (
		pod = new(types.Pod)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> get by name err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyMeta := keyCreate(podStorage, s.keyCreate(namespace, service, deployment, name))
	if err := client.Map(ctx, keyMeta, filter, pod); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> get by name err: %s", name, err.Error())
		return nil, err
	}

	if pod.Meta.Name == "" {
		return nil, errors.New(store.ErrEntityNotFound)
	}

	return pod, nil
}

// ListByNamespace returns pod list from storage by namespace
func (s *PodStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Pod, error) {

	log.V(logLevel).Debugf("storage:etcd:pod:> get list by namespace: %s", namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:pod:> get list by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + podStorage + `\/.+\/(?:meta|status|spec|status)\b`

	var (
		pods = make(map[string]*types.Pod)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> get list by namespace err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	key := keyCreate(podStorage, fmt.Sprintf("%s:", namespace))
	if err := client.MapList(ctx, key, filter, pods); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> err: %s", namespace, err.Error())
		return nil, err
	}

	return pods, nil
}

// ListByService returns pod list from storage by namespace and service names
func (s *PodStorage) ListByService(ctx context.Context, namespace, service string) (map[string]*types.Pod, error) {

	log.V(logLevel).Debugf("storage:etcd:pod:> get list by namespace and service: %s:%s", namespace, service)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:pod:> get list by namespace and service err: %s", err.Error())
		return nil, err
	}

	if len(service) == 0 {
		err := errors.New("service can not be empty")
		log.V(logLevel).Errorf("storage:etcd:pod:> get list by namespace and service err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + podStorage + `\/.+\/(?:meta|status|spec|status)\b`

	var (
		pods = make(map[string]*types.Pod)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:>  get list by namespace and service err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	key := keyCreate(podStorage, fmt.Sprintf("%s:%s:", namespace, service))
	if err := client.MapList(ctx, key, filter, pods); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> err: %s", namespace, err.Error())
		return nil, err
	}

	return pods, nil
}

// ListByDeployment returns pod list from storage by namespace, service and deployment names
func (s *PodStorage) ListByDeployment(ctx context.Context, namespace, service, deployment string) (map[string]*types.Pod, error) {

	log.V(logLevel).Debugf("storage:etcd:pod:> get list by namespace, service and deployment: %s:%s:%s", namespace, service, deployment)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:pod:> get list by namespace, service and deployment err: %s", err.Error())
		return nil, err
	}

	if len(service) == 0 {
		err := errors.New("service can not be empty")
		log.V(logLevel).Errorf("storage:etcd:pod:> get list by namespace, service and deployment err: %s", err.Error())
		return nil, err
	}

	if len(deployment) == 0 {
		err := errors.New("service can not be empty")
		log.V(logLevel).Errorf("storage:etcd:pod:> get list by namespace, service and deployment err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + podStorage + `\/.+\/(?:meta|status|spec|status)\b`

	var (
		pods = make(map[string]*types.Pod)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:>  get list by namespace, service and deployment err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	key := keyCreate(podStorage, fmt.Sprintf("%s:%s:%s:", namespace, service, deployment))
	if err := client.MapList(ctx, key, filter, pods); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> err: %s", err.Error())
		return nil, err
	}

	return pods, nil
}

// Update pod meta
func (s *PodStorage) SetMeta(ctx context.Context, pod *types.Pod) error {

	log.V(logLevel).Debugf("storage:etcd:pod:> update pod meta: %#v", pod)

	if err := s.checkPodExists(ctx, pod); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:>: update pod err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(podStorage, s.keyGet(pod), "meta")
	if err := client.Upsert(ctx, key, pod.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:>: update pod err: %s", err.Error())
		return err
	}

	return nil
}

// Update pod spec
func (s *PodStorage) SetSpec(ctx context.Context, pod *types.Pod) error {

	log.V(logLevel).Debugf("storage:etcd:pod:> update pod spec: %#v", pod)

	if err := s.checkPodExists(ctx, pod); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:>: update pod err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(podStorage, s.keyGet(pod), "spec")
	if err := client.Upsert(ctx, key, pod.Spec, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:>: update pod err: %s", err.Error())
		return err
	}

	return nil
}

// Update pod status
func (s *PodStorage) SetStatus(ctx context.Context, pod *types.Pod) error {
	log.V(logLevel).Debugf("storage:etcd:pod:> update pod status: %#v", pod)

	if err := s.checkPodExists(ctx, pod); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:>: update pod err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(podStorage, s.keyGet(pod), "status")
	if err := client.Upsert(ctx, key, pod.Status, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:>: update pod err: %s", err.Error())
		return err
	}

	return nil
}

// Insert new pod into storage
func (s *PodStorage) Insert(ctx context.Context, pod *types.Pod) error {

	log.V(logLevel).Debugf("storage:etcd:pod:> insert pod: %#v", pod)

	if err := s.checkPodArgument(pod); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> insert pod err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(podStorage, s.keyGet(pod), "meta")
	if err := tx.Create(keyMeta, pod.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> insert pod err: %s", err.Error())
		return err
	}

	keySpec := keyCreate(podStorage, s.keyGet(pod), "spec")
	if err := tx.Create(keySpec, pod.Spec, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> insert pod err: %s", err.Error())
		return err
	}

	keyStatus := keyCreate(podStorage, s.keyGet(pod), "status")
	if err := tx.Create(keyStatus, pod.Status, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> insert pod err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> insert pod err: %s", err.Error())
		return err
	}

	return nil
}

// Update pod in storage
func (s *PodStorage) Update(ctx context.Context, pod *types.Pod) error {

	log.V(logLevel).Debugf("storage:etcd:pod:> update pod: %#v", pod)

	if err := s.checkPodExists(ctx, pod); err != nil {
		return err
	}

	pod.Meta.Updated = time.Now()
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> update pod err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(podStorage, s.keyGet(pod), "meta")
	if err := client.Upsert(ctx, keyMeta, pod.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> update pod err: %s", err.Error())
		return err
	}

	return nil
}

// Remove pod from storage
func (s *PodStorage) Destroy(ctx context.Context, pod *types.Pod) error {

	log.V(logLevel).Debugf("storage:etcd:pod:> update pod: %#v", pod)

	if err := s.checkPodExists(ctx, pod); err != nil {
		return err
	}

	pod.Meta.Updated = time.Now()
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> update pod err: %s", err.Error())
		return err
	}
	defer destroy()

	keySpec := keyCreate(podStorage, s.keyGet(pod), "spec")
	if err := client.Upsert(ctx, keySpec, pod.Spec, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> update pod err: %s", err.Error())
		return err
	}

	return nil
}

// Remove pod from storage
func (s *PodStorage) Remove(ctx context.Context, pod *types.Pod) error {

	if err := s.checkPodExists(ctx, pod); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> remove err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(podStorage, s.keyGet(pod))
	if err := client.DeleteDir(ctx, key); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> remove pod err: %s", err.Error())
		return err
	}

	return nil
}

// Watch pod changes
func (s *PodStorage) Watch(ctx context.Context, pod chan *types.Pod) error {

	log.V(logLevel).Debug("storage:etcd:pod:> watch pod")

	const filter = `\b\/` + podStorage + `\/(.+):(.+):(.+):(.+)/.+\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> watch pod err: %s", err.Error())
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

		if action == types.STORAGEDELEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2], keys[3], keys[4]); err == nil {
			pod <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> watch pod err: %s", err.Error())
		return err
	}

	return nil
}

// Watch pod spec changes
func (s *PodStorage) WatchSpec(ctx context.Context, pod chan *types.Pod) error {

	log.V(logLevel).Debug("storage:etcd:pod:> watch pod")

	const filter = `\b\/` + podStorage + `\/(.+):(.+):(.+):(.+)/spec\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> watch pod err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(podStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 5 {
			return
		}

		if action == types.STORAGEDELEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2], keys[3], keys[4]); err == nil {
			pod <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> watch pod err: %s", err.Error())
		return err
	}

	return nil
}

// Watch pod status changes
func (s *PodStorage) WatchStatus(ctx context.Context, pod chan *types.Pod) error {
	log.V(logLevel).Debug("storage:etcd:pod:> watch pod")

	const filter = `\b\/` + podStorage + `\/(.+):(.+):(.+):(.+)/status\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> watch pod err: %s", err.Error())
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

		if action == types.STORAGEDELEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2], keys[3], keys[4]); err == nil {
			pod <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> watch pod err: %s", err.Error())
		return err
	}

	return nil
}

// Clear pod storage
func (s *PodStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:pod:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, podStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:pod:> clear err: %s", err.Error())
		return err
	}

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
func (s *PodStorage) checkPodExists(ctx context.Context, pod *types.Pod) error {

	if err := s.checkPodArgument(pod); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:pod:> check pod exists")

	if _, err := s.Get(ctx, pod.Meta.Namespace, pod.Meta.Service, pod.Meta.Deployment, pod.Meta.Name); err != nil {
		log.V(logLevel).Debugf("storage:etcd:pod:> check pod exists err: %s", err.Error())
		return err
	}

	return nil
}
