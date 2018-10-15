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
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

const (
	logRoutePrefix = "distribution:route"
)

type Route struct {
	context context.Context
	storage storage.Storage
}

func (n *Route) Get(namespace, name string) (*types.Route, error) {

	log.V(logLevel).Debug("%s:get:> get route by id %s/%s", logRoutePrefix, namespace, name)

	route := new(types.Route)

	err := n.storage.Get(n.context, n.storage.Collection().Route(), n.storage.Key().Route(namespace, name), &route, nil)
	if err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> in namespace %s by name %s not found", logRoutePrefix, namespace, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> in namespace %s by name %s error: %v", logRoutePrefix, namespace, name, err)
		return nil, err
	}

	return route, nil
}

func (n *Route) List() (*types.RouteList, error) {

	log.V(logLevel).Debugf("%s:listspec:> list specs", logRoutePrefix)

	list := types.NewRouteList()

	//TODO: change map to list
	err := n.storage.List(n.context, n.storage.Collection().Route(), types.EmptyString, list, nil)
	if err != nil {
		log.V(logLevel).Error("%s:listbynamespace:> list route err: %v", logRoutePrefix, err)
		return list, err
	}

	return list, nil
}

func (n *Route) ListByNamespace(namespace string) (*types.RouteList, error) {

	log.V(logLevel).Debug("%s:listbynamespace:> list route", logRoutePrefix)

	list := types.NewRouteList()

	err := n.storage.List(n.context, n.storage.Collection().Route(), n.storage.Filter().Route().ByNamespace(namespace), list, nil)
	if err != nil {
		log.V(logLevel).Error("%s:listbynamespace:> list route err: %v", logRoutePrefix, err)
		return list, err
	}

	log.V(logLevel).Debugf("%s:listbynamespace:> list route result: %d", logRoutePrefix, len(list.Items))

	return list, nil
}

func (n *Route) Create(namespace *types.Namespace, route *types.Route) (*types.Route, error) {

	log.V(logLevel).Debugf("%s:create:> create route %#v", logRoutePrefix, route.Meta.Name)

	route.Meta.SetDefault()
	route.Status.State = types.StatusInitialized
	route.Spec.Domain = fmt.Sprintf("%s.%s.%s", strings.ToLower(route.Meta.Name), strings.ToLower(namespace.Meta.Name),  viper.GetString("domain.external"))
	route.SelfLink()

	if err := n.storage.Put(n.context, n.storage.Collection().Route(),
		n.storage.Key().Route(route.Meta.Namespace, route.Meta.Name), route, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert route err: %v", logRoutePrefix, err)
		return nil, err
	}

	return route, nil
}

func (n *Route) Update(route *types.Route) (*types.Route, error) {

	log.V(logLevel).Debugf("%s:update:> update route %s", logRoutePrefix, route.Meta.Name)
	route.Status.State = types.StateProvision

	if err := n.storage.Set(n.context, n.storage.Collection().Route(),
		n.storage.Key().Route(route.Meta.Namespace, route.Meta.Name), route, nil); err != nil {
		log.V(logLevel).Errorf("%s:update:> update route err: %v", logRoutePrefix, err)
		return nil, err
	}

	return route, nil
}

func (n *Route) SetStatus(route *types.Route, status *types.RouteStatus) error {

	if route == nil {
		log.V(logLevel).Warnf("%s:setstatus:> invalid argument %v", logRoutePrefix, route)
		return nil
	}

	log.V(logLevel).Debugf("%s:setstate:> set state route %s -> %#v", logRoutePrefix, route.Meta.Name, status)

	route.Status = *status
	if err := n.storage.Set(n.context, n.storage.Collection().Route(),
		n.storage.Key().Route(route.Meta.Namespace, route.Meta.Name), route, nil); err != nil {
		log.Errorf("%s:setstatus:> pod set status err: %v", logRoutePrefix, err)
		return err
	}

	return nil
}

func (n *Route) Remove(route *types.Route) error {

	log.V(logLevel).Debugf("%s:remove:> remove route %#v", logRoutePrefix, route)

	if err := n.storage.Del(n.context, n.storage.Collection().Route(),
		n.storage.Key().Route(route.Meta.Namespace, route.Meta.Name)); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove route  err: %v", logRoutePrefix, err)
		return err
	}

	return nil
}

func (n *Route) Watch(ch chan types.RouteEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch routes", logRoutePrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-n.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.RouteEvent{}
				res.Action = e.Action
				res.Name = e.Name

				route := new(types.Route)

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
	if err := n.storage.Watch(n.context, n.storage.Collection().Route(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewRouteModel(ctx context.Context, stg storage.Storage) *Route {
	return &Route{ctx, stg}
}
