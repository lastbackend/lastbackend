//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logRoutePrefix = "distribution:route"
)

type Route struct {
	context context.Context
	storage storage.IStorage
}

func (r *Route) Runtime() (*models.System, error) {

	log.Debugf("%s:get:> get route runtime info", logRoutePrefix)
	runtime, err := r.storage.Info(r.context, r.storage.Collection().Pod(), "")
	if err != nil {
		log.Errorf("%s:get:> get runtime info error: %s", logRoutePrefix, err)
		return &runtime.System, err
	}
	return &runtime.System, nil
}

func (r *Route) List() (*models.RouteList, error) {

	log.Debugf("%s:listspec:> list specs", logRoutePrefix)

	list := models.NewRouteList()

	//TODO: change map to list
	err := r.storage.List(r.context, r.storage.Collection().Route(), models.EmptyString, list, nil)
	if err != nil {
		log.Error("%s:listbynamespace:> list route err: %v", logRoutePrefix, err)
		return list, err
	}

	return list, nil
}

func (r *Route) ListByNamespace(namespace string) (*models.RouteList, error) {

	log.Debug("%s:listbynamespace:> list route", logRoutePrefix)

	list := models.NewRouteList()

	err := r.storage.List(r.context, r.storage.Collection().Route(), r.storage.Filter().Route().ByNamespace(namespace), list, nil)
	if err != nil {
		log.Error("%s:listbynamespace:> list route err: %v", logRoutePrefix, err)
		return list, err
	}

	log.Debugf("%s:listbynamespace:> list route result: %d", logRoutePrefix, len(list.Items))

	return list, nil
}

func (r *Route) Get(namespace, name string) (*models.Route, error) {

	log.Debug("%s:get:> get route by id %s/%s", logRoutePrefix, namespace, name)

	route := new(models.Route)
	key := models.NewRouteSelfLink(namespace, name).String()

	err := r.storage.Get(r.context, r.storage.Collection().Route(), key, &route, nil)
	if err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			log.Warnf("%s:get:> in namespace %s by name %s not found", logRoutePrefix, namespace, name)
			return nil, nil
		}

		log.Errorf("%s:get:> in namespace %s by name %s error: %v", logRoutePrefix, namespace, name, err)
		return nil, err
	}

	return route, nil
}

func (r *Route) Add(namespace *models.Namespace, route *models.Route) (*models.Route, error) {

	log.Debugf("%s:create:> create route %#v", logRoutePrefix, route.Meta.Name)

	route.Status.State = models.StateCreated

	if err := r.storage.Put(r.context, r.storage.Collection().Route(),
		route.SelfLink().String(), route, nil); err != nil {
		log.Errorf("%s:create:> insert route err: %v", logRoutePrefix, err)
		return nil, err
	}

	return route, nil
}

func (r *Route) Set(route *models.Route) (*models.Route, error) {

	log.Debugf("%s:update:> update route %s", logRoutePrefix, route.Meta.Name)

	if err := r.storage.Set(r.context, r.storage.Collection().Route(),
		route.SelfLink().String(), route, nil); err != nil {
		log.Errorf("%s:update:> update route err: %v", logRoutePrefix, err)
		return nil, err
	}

	return route, nil
}

func (r *Route) Del(route *models.Route) error {

	log.Debugf("%s:remove:> remove route %#v", logRoutePrefix, route)

	if err := r.storage.Del(r.context, r.storage.Collection().Route(),
		route.SelfLink().String()); err != nil {
		log.Errorf("%s:remove:> remove route  err: %v", logRoutePrefix, err)
		return err
	}

	return nil
}

func (r *Route) Watch(ch chan models.RouteEvent, rev *int64) error {

	log.Debugf("%s:watch:> watch routes", logRoutePrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-r.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := models.RouteEvent{}
				res.Action = e.Action
				res.Name = e.Name

				route := new(models.Route)

				if err := json.Unmarshal(e.Data.([]byte), route); err != nil {
					log.Errorf("%s:> parse data err: %v", logRoutePrefix, err)
					continue
				}

				res.Data = route

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := r.storage.Watch(r.context, r.storage.Collection().Route(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func (r *Route) ManifestMap(ingress string) (*models.RouteManifestMap, error) {
	log.Debug("%s:ManifestMap:> get route manifest map by ingress %s", logRoutePrefix, ingress)

	var (
		mf = models.NewRouteManifestMap()
	)

	if err := r.storage.Map(r.context, r.storage.Collection().Manifest().Route(ingress), models.EmptyString, mf, nil); err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s:ManifestMap:> err: %s", logRoutePrefix, err.Error())
			return nil, err
		}

		return nil, nil
	}

	return mf, nil
}

func (r *Route) ManifestGet(ingress, route string) (*models.RouteManifest, error) {
	log.Debug("%s:ManifestGet:> get route manifest by name %s", logRoutePrefix, route)

	var (
		mf = new(models.RouteManifest)
	)

	if err := r.storage.Get(r.context, r.storage.Collection().Manifest().Route(ingress), route, &mf, nil); err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return mf, nil
}

func (r *Route) ManifestAdd(ingress, route string, manifest *models.RouteManifest) error {

	log.Debugf("%s:ManifestAdd:> ", logRoutePrefix)

	if err := r.storage.Put(r.context, r.storage.Collection().Manifest().Route(ingress), route, manifest, nil); err != nil {
		log.Errorf("%s:ManifestAdd:> err :%s", logRoutePrefix, err.Error())
		return err
	}

	return nil
}

func (r *Route) ManifestSet(ingress, route string, manifest *models.RouteManifest) error {
	log.Debugf("%s:ManifestSet:> ", logRoutePrefix)

	if err := r.storage.Set(r.context, r.storage.Collection().Manifest().Route(ingress), route, manifest, nil); err != nil {
		log.Errorf("%s:ManifestSet:> err :%s", logRoutePrefix, err.Error())
		return err
	}

	return nil
}

func (r *Route) ManifestDel(ingress, route string) error {
	log.Debugf("%s:ManifestDel:> %s on ingress %s", logRoutePrefix, route, ingress)

	if err := r.storage.Del(r.context, r.storage.Collection().Manifest().Route(ingress), route); err != nil {
		log.Errorf("%s:ManifestDel:> err :%s", logRoutePrefix, err.Error())
		return err
	}

	return nil
}

func (r *Route) ManifestWatch(ingress string, ch chan models.RouteManifestEvent, rev *int64) error {
	log.Debugf("%s:watch:> watch routes manifest", logRoutePrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	var f, c string

	if ingress != models.EmptyString {
		f = fmt.Sprintf(`\b.+\/%s\/%s\/(.+)\b`, ingress, storage.RouteKind)
		c = r.storage.Collection().Manifest().Route(ingress)
	} else {
		f = fmt.Sprintf(`\b.+\/(.+)\/%s\/(.+)\b`, storage.RouteKind)
		c = r.storage.Collection().Manifest().Ingress()
	}

	rg, err := regexp.Compile(f)
	if err != nil {
		log.Errorf("%s:> filter compile err: %v", logRoutePrefix, err.Error())
		return err
	}

	go func() {
		for {
			select {
			case <-r.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				keys := rg.FindStringSubmatch(e.Storage.Key)
				if len(keys) == 0 {
					continue
				}

				res := models.RouteManifestEvent{}
				res.Action = e.Action
				res.Name = e.Name
				res.SelfLink = e.SelfLink

				if ingress != models.EmptyString {
					res.Ingress = ingress
				} else {
					res.Ingress = keys[1]
				}

				manifest := new(models.RouteManifest)

				if err := json.Unmarshal(e.Data.([]byte), manifest); err != nil {
					log.Errorf("%s:> parse data err: %v", logRoutePrefix, err)
					continue
				}

				res.Data = manifest

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := r.storage.Watch(r.context, c, watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewRouteModel(ctx context.Context, stg storage.IStorage) *Route {
	return &Route{ctx, stg}
}
