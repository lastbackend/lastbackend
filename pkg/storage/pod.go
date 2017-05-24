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
	"errors"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
	"strings"
)

const podStorage = "pods"

// Pod Service type for interface in interfaces folder
type PodStorage struct {
	IPod
	log    logger.ILogger
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *PodStorage) GetByName(ctx context.Context, namespace, name string) (*types.Pod, error) {

	s.log.V(logLevel).Debugf("Storage: Pod: get by name: %s in namespace: %s", name, namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		s.log.V(logLevel).Errorf("Storage: Pod: get pod err: %s", err.Error())
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		s.log.V(logLevel).Errorf("Storage: Pod: get pod err: %s", err.Error())
		return nil, err
	}

	var (
		pod            = new(types.Pod)
		podName        = strings.Replace(pod.Meta.Name, ":", "-", -1)
		filterEndpoint = `\b.+` + endpointStorage + `\/` + podName + `\..+\b`
		endpoints      = make(map[string][]string)
	)

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyMeta := keyCreate(podStorage, namespace, name)
	if err := client.Get(ctx, keyMeta, pod); err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: get pod `%s` err: %s", name, err.Error())
		return nil, err
	}

	if pod.Meta.Name == "" {
		return nil, errors.New(store.ErrKeyNotFound)
	}

	keyEndpoints := keyCreate(endpointStorage)
	if err := client.Map(ctx, keyEndpoints, filterEndpoint, endpoints); err != nil && err.Error() != store.ErrKeyNotFound {
		s.log.V(logLevel).Errorf("Storage: Pod: map endpoints err: %s", err.Error())
		return nil, err
	}

	for pod.Meta.Endpoint = range endpoints {
		break
	}

	return pod, nil
}

func (s *PodStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Pod, error) {

	s.log.V(logLevel).Debugf("Storage: Pod: get pods list in namespace: %s", namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		s.log.V(logLevel).Errorf("Storage: Pod: get pod err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	pods := make(map[string]*types.Pod)
	keyList := keyCreate(podStorage, namespace)
	if err := client.Map(ctx, keyList, "", &pods); err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: map pods in namespace `%s` err: %s", namespace, err.Error())
		return pods, err
	}

	for _, pod := range pods {
		name := strings.Replace(pod.Meta.Name, ":", "-", -1)
		filterEndpoint := `\b.+` + endpointStorage + `\/` + name + `\..+\b`
		endpoints := make(map[string][]string)
		keyEndpoints := keyCreate(endpointStorage)
		if err := client.Map(ctx, keyEndpoints, filterEndpoint, endpoints); err != nil && err.Error() != store.ErrKeyNotFound {
			s.log.V(logLevel).Errorf("Storage: Pod: map endpoints err: %s", err.Error())
			return pods, err
		}

		for pod.Meta.Endpoint = range endpoints {
			break
		}
	}

	return pods, nil
}

func (s *PodStorage) ListByService(ctx context.Context, namespace, service string) ([]*types.Pod, error) {

	s.log.V(logLevel).Debugf("Storage: Pod: get pods list by service: %s in namespace: %s", service, namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		s.log.V(logLevel).Errorf("Storage: Pod: get pods list  err: %s", err.Error())
		return nil, err
	}

	if len(service) == 0 {
		err := errors.New("service can not be empty")
		s.log.V(logLevel).Errorf("Storage: Pod: get pods list err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	pods := make([]*types.Pod, 0)
	keyServiceList := keyCreate(podStorage, namespace, service)
	if err := client.List(ctx, keyServiceList, "", &pods); err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: pods list err: %s", err.Error())
		return nil, err
	}

	for _, pod := range pods {
		filterEndpoint := `\b.+` + endpointStorage + `\/` + pod.Meta.Name + `-` + namespace + `\..+\b`
		endpoints := make(map[string][]string)
		keyEndpoints := keyCreate(endpointStorage)
		if err := client.Map(ctx, keyEndpoints, filterEndpoint, endpoints); err != nil && err.Error() != store.ErrKeyNotFound {
			s.log.V(logLevel).Errorf("Storage: Pod: map endpoints err: %s", err.Error())
			return nil, err
		}

		for pod.Meta.Endpoint = range endpoints {
			break
		}
	}

	return pods, nil
}

func (s *PodStorage) Upsert(ctx context.Context, namespace string, pod *types.Pod) error {

	s.log.V(logLevel).Debugf("Storage: Pod: upsert pod: %#v in namespace: %s", pod, namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		s.log.V(logLevel).Errorf("Storage: Pod: upsert pod list  err: %s", err.Error())
		return err
	}

	if pod == nil {
		err := errors.New("pod can not be nil")
		s.log.V(logLevel).Errorf("Storage: Pod: upsert pod list err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(podStorage, namespace, pod.Meta.Name)
	if err := client.Upsert(ctx, keyMeta, pod, nil, 0); err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: upsert pod err: %s", err.Error())
		return err
	}

	return nil
}

func (s *PodStorage) Update(ctx context.Context, namespace string, pod *types.Pod) error {

	s.log.V(logLevel).Debugf("Storage: Pod: update pod: %#v in namespace: %s", pod, namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		s.log.V(logLevel).Errorf("Storage: Pod: update pod list  err: %s", err.Error())
		return err
	}

	if pod == nil {
		err := errors.New("pod can not be nil")
		s.log.V(logLevel).Errorf("Storage: Pod: update pod list err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(podStorage, namespace, pod.Meta.Name)

	if err := client.Update(ctx, keyMeta, pod, nil, 0); err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: update pod err: %s", err.Error())
		return err
	}

	return nil
}

func (s *PodStorage) Remove(ctx context.Context, namespace string, pod *types.Pod) error {

	s.log.V(logLevel).Debugf("Storage: Pod: remove pod: %#v in namespace: %s", pod, namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		s.log.V(logLevel).Errorf("Storage: Pod: remove pod list  err: %s", err.Error())
		return err
	}

	if pod == nil {
		err := errors.New("pod can not be nil")
		s.log.V(logLevel).Errorf("Storage: Pod: remove pod list err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(podStorage, namespace, pod.Meta.Name)
	tx.Delete(keyMeta)

	KeyNodePod := keyCreate(nodeStorage, pod.Node.ID, "spec", "pods", pod.Meta.Name)
	tx.Delete(KeyNodePod)

	if err := tx.Commit(); err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: commit transaction err: %s", err.Error())
		return err
	}

	return nil
}

func (s *PodStorage) Watch(ctx context.Context, pod chan *types.Pod) error {

	s.log.V(logLevel).Debug("Storage: Pod: watch pod")

	const filter = `\b\/` + podStorage + `\/(.+)/(.+)\b`
	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
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

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		s.log.V(logLevel).Errorf("Storage: Pod: watch pod err: %s", err.Error())
		return err
	}

	return nil
}

func newPodStorage(config store.Config, log logger.ILogger, util IUtil) *PodStorage {
	s := new(PodStorage)
	s.log = log
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config, log)
	}
	return s
}
