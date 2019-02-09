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

package distribution

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"

	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"regexp"
)

const (
	logEndpointPrefix = "distribution:endpoint"
)

type Endpoint struct {
	context context.Context
	storage storage.Storage
}

func (e *Endpoint) Get(namespace, service string) (*types.Endpoint, error) {

	log.V(logLevel).Debugf("%s:get:> get endpoint by namespace %s and service %s", logEndpointPrefix, namespace, service)

	item := new(types.Endpoint)

	err := e.storage.Get(e.context, e.storage.Collection().Endpoint(), types.NewEndpointSelfLink(namespace, service).String(), &item, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> get endpoint err: %v", logEndpointPrefix, err)
		return nil, err
	}

	return item, nil
}

func (e *Endpoint) ListByNamespace(namespace string) (*types.EndpointList, error) {
	log.V(logLevel).Debugf("%s:listbynamespace:> in namespace: %s", namespace)

	list := types.NewEndpointList()

	err := e.storage.List(e.context, e.storage.Collection().Endpoint(), e.storage.Filter().Endpoint().ByNamespace(namespace), list, nil)
	if err != nil {
		log.Errorf("%s:listbynamespace:> in namespace: %s err: %v", logEndpointPrefix, namespace, err)
		return nil, err
	}

	return list, nil
}

func (e *Endpoint) Create(namespace, service string, opts *types.EndpointCreateOptions) (*types.Endpoint, error) {
	endpoint := new(types.Endpoint)

	endpoint.Meta.Name = service
	endpoint.Meta.Namespace = namespace
	endpoint.Meta.SetDefault()
	endpoint.Meta.SelfLink = *types.NewEndpointSelfLink(namespace, service)

	endpoint.Status.State = types.StateCreated
	endpoint.Spec.PortMap = make(map[uint16]string, 0)

	for k, v := range opts.Ports {
		endpoint.Spec.PortMap[k] = v
	}

	endpoint.Spec.Policy = opts.Policy
	endpoint.Spec.Strategy.Route = opts.RouteStrategy
	endpoint.Spec.Strategy.Bind = opts.BindStrategy

	endpoint.Spec.IP = opts.IP
	endpoint.Spec.Domain = opts.Domain

	if err := e.storage.Put(e.context, e.storage.Collection().Endpoint(), endpoint.SelfLink().String(), endpoint, nil); err != nil {
		log.Errorf("%s:create:> distribution create endpoint: %s err: %v", logEndpointPrefix, endpoint.SelfLink(), err)
		return nil, err
	}

	return endpoint, nil
}

func (e *Endpoint) Update(endpoint *types.Endpoint, opts *types.EndpointUpdateOptions) (*types.Endpoint, error) {
	log.V(logLevel).Debugf("%s:update:> endpoint: %s", logEndpointPrefix, endpoint.SelfLink())

	if len(opts.Ports) != 0 {
		endpoint.Spec.PortMap = make(map[uint16]string, 0)
		for k, v := range opts.Ports {
			endpoint.Spec.PortMap[k] = v
		}
	}

	if opts.IP != nil {
		endpoint.Spec.IP = *opts.IP
	}

	endpoint.Spec.Policy = opts.Policy
	endpoint.Spec.Strategy.Route = opts.RouteStrategy
	endpoint.Spec.Strategy.Bind = opts.BindStrategy

	if err := e.storage.Set(e.context, e.storage.Collection().Endpoint(),
		endpoint.SelfLink().String(), endpoint, nil); err != nil {
		log.Errorf("%s:create:> distribution update endpoint: %s err: %v", logEndpointPrefix, endpoint.SelfLink(), err)
		return nil, err
	}

	return endpoint, nil
}

func (e *Endpoint) SetSpec(endpoint *types.Endpoint, spec *types.EndpointSpec) (*types.Endpoint, error) {
	endpoint.Spec = *spec
	if err := e.storage.Set(e.context, e.storage.Collection().Endpoint(),
		endpoint.SelfLink().String(), endpoint, nil); err != nil {
		log.Errorf("%s:create:> distribution update endpoint spec: %s err: %v", logEndpointPrefix, endpoint.SelfLink(), err)
		return nil, err
	}
	return endpoint, nil
}

func (e *Endpoint) Remove(endpoint *types.Endpoint) error {
	log.V(logLevel).Debugf("%s:remove:> remove endpoint %s", logEndpointPrefix, endpoint.Meta.Name)
	if err := e.storage.Del(e.context, e.storage.Collection().Endpoint(),
		endpoint.SelfLink().String()); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove endpoint %s err: %v", logEndpointPrefix, endpoint.Meta.Name, err)
		return err
	}

	return nil
}

// Watch endpoint changes
func (e *Endpoint) Watch(ch chan types.EndpointEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch endpoint", logEndpointPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-e.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.EndpointEvent{}
				res.Action = e.Action
				res.Name = e.Name

				endpoint := new(types.Endpoint)

				if err := json.Unmarshal(e.Data.([]byte), endpoint); err != nil {
					log.Errorf("%s:> parse data err: %v", logEndpointPrefix, err)
					continue
				}

				res.Data = endpoint

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := e.storage.Watch(e.context, e.storage.Collection().Endpoint(), watcher, opts); err != nil {
		return err
	}

	return nil
}

// Get network subnet manifests map
func (e *Endpoint) ManifestMap() (*types.EndpointManifestMap, error) {
	log.V(logLevel).Debugf("%s:EndpointManifestMap:> ", logEndpointPrefix)

	var (
		mf = types.NewEndpointManifestMap()
	)

	if err := e.storage.Map(e.context, e.storage.Collection().Manifest().Endpoint(), types.EmptyString, mf, nil); err != nil {
		log.Errorf("%s:EndpointManifestMap:> err :%s", logEndpointPrefix, err.Error())
		return nil, err
	}

	return mf, nil
}

// Get particular network manifest
func (e *Endpoint) ManifestGet(selflink string) (*types.EndpointManifest, error) {
	log.V(logLevel).Debugf("%s:EndpointManifestGet:> ", logEndpointPrefix)

	var (
		mf = new(types.EndpointManifest)
	)

	if err := e.storage.Get(e.context, e.storage.Collection().Manifest().Endpoint(), selflink, &mf, nil); err != nil {
		log.Errorf("%s:EndpointManifestGet:> err :%s", logEndpointPrefix, err.Error())

		if errors.Storage().IsErrEntityNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return mf, nil
}

// Add particular network manifest
func (e *Endpoint) ManifestAdd(selflink string, manifest *types.EndpointManifest) error {

	log.V(logLevel).Debugf("%s:EndpointManifestAdd:> ", logEndpointPrefix)

	if err := e.storage.Put(e.context,
		e.storage.Collection().Manifest().Endpoint(),
		selflink,
		&manifest, nil); err != nil {
		log.Errorf("%s:EndpointManifestAdd:> err :%s", logEndpointPrefix, err.Error())
		return err
	}

	return nil
}

// Set particular network manifest
func (e *Endpoint) ManifestSet(selflink string, manifest *types.EndpointManifest) error {
	log.V(logLevel).Debugf("%s:EndpointManifestSet:> ", logEndpointPrefix)

	if err := e.storage.Set(e.context, e.storage.Collection().Manifest().Endpoint(), selflink, manifest, nil); err != nil {
		log.Errorf("%s:EndpointManifestSet:> err :%s", logEndpointPrefix, err.Error())
		return err
	}

	return nil
}

// Del particular network manifest
func (e *Endpoint) ManifestDel(selflink string) error {
	log.V(logLevel).Debugf("%s:EndpointManifestDel:> ", logEndpointPrefix)

	if err := e.storage.Del(e.context, e.storage.Collection().Manifest().Endpoint(), selflink); err != nil {
		log.Errorf("%s:EndpointManifestDel:> err :%s", logEndpointPrefix, err.Error())
		return err
	}

	return nil
}

// watch subnet manifests
func (e *Endpoint) ManifestWatch(ch chan types.EndpointManifestEvent, rev *int64) error {
	log.V(logLevel).Debugf("%s:EndpointManifestWatch:> watch manifest ", logEndpointPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()
	r, _ := regexp.Compile(`\b.+\/(.+)\b`)

	go func() {
		for {
			select {
			case <-e.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				keys := r.FindStringSubmatch(e.Storage.Key)
				if len(keys) == 0 {
					continue
				}

				res := types.EndpointManifestEvent{}
				res.Action = e.Action
				res.Name = e.Name
				res.SelfLink = e.SelfLink

				manifest := new(types.EndpointManifest)

				if err := json.Unmarshal(e.Data.([]byte), manifest); err != nil {
					log.Errorf("%s:> parse data err: %v", logEndpointPrefix, err)
					continue
				}

				res.Data = manifest

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := e.storage.Watch(e.context, e.storage.Collection().Manifest().Endpoint(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func (e *Endpoint) ManifestGetSelfLink(namespace, service string) string {
	return types.NewEndpointSelfLink(namespace, service).String()
}

func NewEndpointModel(ctx context.Context, stg storage.Storage) *Endpoint {
	return &Endpoint{ctx, stg}
}
