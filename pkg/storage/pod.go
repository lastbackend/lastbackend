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
	"regexp"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const podStorage = "pods"

// Pod Service type for interface in interfaces folder
type PodStorage struct {
	IPod
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *PodStorage) GetByName(ctx context.Context, app, name string) (*types.Pod, error) {

	log.V(logLevel).Debugf("Storage: Pod: get by name: %s in app: %s", name, app)

	if len(app) == 0 {
		err := errors.New("app can not be empty")
		log.V(logLevel).Errorf("Storage: Pod: get pod err: %s", err.Error())
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("Storage: Pod: get pod err: %s", err.Error())
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
		log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyMeta := keyCreate(podStorage, app, name)
	if err := client.Get(ctx, keyMeta, pod); err != nil {
		log.V(logLevel).Errorf("Storage: Pod: get pod `%s` err: %s", name, err.Error())
		return nil, err
	}

	if pod.Meta.Name == "" {
		return nil, errors.New(store.ErrKeyNotFound)
	}

	keyEndpoints := keyCreate(endpointStorage)
	if err := client.Map(ctx, keyEndpoints, filterEndpoint, endpoints); err != nil && err.Error() != store.ErrKeyNotFound {
		log.V(logLevel).Errorf("Storage: Pod: map endpoints err: %s", err.Error())
		return nil, err
	}

	for pod.Meta.Endpoint = range endpoints {
		break
	}

	return pod, nil
}

func (s *PodStorage) ListByApp(ctx context.Context, app string) (map[string]*types.Pod, error) {

	log.V(logLevel).Debugf("Storage: Pod: get pods list in app: %s", app)

	if len(app) == 0 {
		err := errors.New("app can not be empty")
		log.V(logLevel).Errorf("Storage: Pod: get pod err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	pods := make(map[string]*types.Pod)
	keyList := keyCreate(podStorage, app)
	if err := client.Map(ctx, keyList, "", &pods); err != nil {
		log.V(logLevel).Errorf("Storage: Pod: map pods in app `%s` err: %s", app, err.Error())
		return nil, err
	}

	for _, pod := range pods {
		name := strings.Replace(pod.Meta.Name, ":", "-", -1)
		filterEndpoint := `\b.+` + endpointStorage + `\/` + name + `\..+\b`
		endpoints := make(map[string][]string)
		keyEndpoints := keyCreate(endpointStorage)
		if err := client.Map(ctx, keyEndpoints, filterEndpoint, endpoints); err != nil && err.Error() != store.ErrKeyNotFound {
			log.V(logLevel).Errorf("Storage: Pod: map endpoints err: %s", err.Error())
			return nil, err
		}

		for pod.Meta.Endpoint = range endpoints {
			break
		}
	}

	return pods, nil
}

func (s *PodStorage) ListByService(ctx context.Context, app, service string) ([]*types.Pod, error) {

	log.V(logLevel).Debugf("Storage: Pod: get pods list by service: %s in app: %s", service, app)

	if len(app) == 0 {
		err := errors.New("app can not be empty")
		log.V(logLevel).Errorf("Storage: Pod: get pods list  err: %s", err.Error())
		return nil, err
	}

	if len(service) == 0 {
		err := errors.New("service can not be empty")
		log.V(logLevel).Errorf("Storage: Pod: get pods list err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	pods := make([]*types.Pod, 0)
	keyServiceList := keyCreate(podStorage, app, service)
	if err := client.List(ctx, keyServiceList, "", &pods); err != nil {
		log.V(logLevel).Errorf("Storage: Pod: pods list err: %s", err.Error())
		return nil, err
	}

	for _, pod := range pods {
		filterEndpoint := `\b.+` + endpointStorage + `\/` + pod.Meta.Name + `-` + app + `\..+\b`
		endpoints := make(map[string][]string)
		keyEndpoints := keyCreate(endpointStorage)
		if err := client.Map(ctx, keyEndpoints, filterEndpoint, endpoints); err != nil && err.Error() != store.ErrKeyNotFound {
			log.V(logLevel).Errorf("Storage: Pod: map endpoints err: %s", err.Error())
			return nil, err
		}

		for pod.Meta.Endpoint = range endpoints {
			break
		}
	}

	return pods, nil
}

func (s *PodStorage) Upsert(ctx context.Context, app string, pod *types.Pod) error {

	log.V(logLevel).Debugf("Storage: Pod: upsert pod: %#v in app: %s", pod, app)

	if len(app) == 0 {
		err := errors.New("app can not be empty")
		log.V(logLevel).Errorf("Storage: Pod: upsert pod list  err: %s", err.Error())
		return err
	}

	if pod == nil {
		err := errors.New("pod can not be nil")
		log.V(logLevel).Errorf("Storage: Pod: upsert pod list err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(podStorage, app, pod.Meta.Name)
	if err := client.Upsert(ctx, keyMeta, pod, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Pod: upsert pod err: %s", err.Error())
		return err
	}

	return nil
}

func (s *PodStorage) Update(ctx context.Context, app string, pod *types.Pod) error {

	log.V(logLevel).Debugf("Storage: Pod: update pod: %#v in app: %s", pod, app)

	if len(app) == 0 {
		err := errors.New("app can not be empty")
		log.V(logLevel).Errorf("Storage: Pod: update pod list  err: %s", err.Error())
		return err
	}

	if pod == nil {
		err := errors.New("pod can not be nil")
		log.V(logLevel).Errorf("Storage: Pod: update pod list err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(podStorage, app, pod.Meta.Name)

	if err := client.Update(ctx, keyMeta, pod, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Pod: update pod err: %s", err.Error())
		return err
	}

	return nil
}

func (s *PodStorage) Remove(ctx context.Context, app string, pod *types.Pod) error {

	log.V(logLevel).Debugf("Storage: Pod: remove pod: %#v in app: %s", pod, app)

	if len(app) == 0 {
		err := errors.New("app can not be empty")
		log.V(logLevel).Errorf("Storage: Pod: remove pod list  err: %s", err.Error())
		return err
	}

	if pod == nil {
		err := errors.New("pod can not be nil")
		log.V(logLevel).Errorf("Storage: Pod: remove pod list err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(podStorage, app, pod.Meta.Name)
	tx.Delete(keyMeta)

	KeyNodePod := keyCreate(nodeStorage, pod.Node.ID, "spec", "pods", pod.Meta.Name)
	tx.Delete(KeyNodePod)

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("Storage: Pod: commit transaction err: %s", err.Error())
		return err
	}

	return nil
}

func (s *PodStorage) Watch(ctx context.Context, pod chan *types.Pod) error {

	log.V(logLevel).Debug("Storage: Pod: watch pod")

	const filter = `\b\/` + podStorage + `\/(.+)/(.+)\b`
	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Pod: create client err: %s", err.Error())
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
		log.V(logLevel).Errorf("Storage: Pod: watch pod err: %s", err.Error())
		return err
	}

	return nil
}

func newPodStorage(config store.Config) *PodStorage {
	s := new(PodStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
