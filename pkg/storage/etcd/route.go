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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
	"time"
)

const routeStorage = "routes"

type RouteStorage struct {
	storage.Route
}

// Get route by name
func (s *RouteStorage) Get(ctx context.Context, namespace, name string) (*types.Route, error) {

	log.V(logLevel).Debugf("storage:etcd:route:> get by name: %s", name)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:route:> get by name err: %s", err.Error())
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("storage:etcd:route:> get by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + routeStorage + `\/.+\/(?:meta|status|spec)\b`
	var (
		route = new(types.Route)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> get by name err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyMeta := keyCreate(routeStorage, s.keyCreate(namespace, name))
	if err := client.Map(ctx, keyMeta, filter, route); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> get by name err: %s", name, err.Error())
		return nil, err
	}

	if route.Meta.Name == "" {
		return nil, errors.New(store.ErrEntityNotFound)
	}

	return route, nil
}

// Get routes by namespace name
func (s *RouteStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Route, error) {

	log.V(logLevel).Debugf("storage:etcd:route:> get list by namespace: %s", namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:route:> get list by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + routeStorage + `\/.+\/(?:meta|status|spec)\b`

	var (
		routes = make(map[string]*types.Route)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> get list by namespace err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	key := keyCreate(routeStorage, fmt.Sprintf("%s:", namespace))
	if err := client.MapList(ctx, key, filter, routes); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> err: %s", namespace, err.Error())
		return nil, err
	}

	return routes, nil

}

// Update route spec
func (s *RouteStorage) SetSpec(ctx context.Context, route *types.Route) error {
	log.V(logLevel).Debugf("storage:etcd:route:> update route spec: %#v", route)

	if err := s.checkRouteExists(ctx, route); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:>: update route err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(routeStorage, s.keyGet(route), "spec")
	if err := client.Upsert(ctx, key, route.Spec, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:>: update route err: %s", err.Error())
		return err
	}

	return nil
}

// Update route status
func (s *RouteStorage) SetStatus(ctx context.Context, route *types.Route) error {

	log.V(logLevel).Debugf("storage:etcd:route:> update route status: %#v", route)

	if err := s.checkRouteExists(ctx, route); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:>: update route err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(routeStorage, s.keyGet(route), "status")
	if err := client.Upsert(ctx, key, route.Status, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:>: update route err: %s", err.Error())
		return err
	}

	return nil
}

// Insert new route
func (s *RouteStorage) Insert(ctx context.Context, route *types.Route) error {

	log.V(logLevel).Debugf("storage:etcd:route:> insert route: %#v", route)

	if err := s.checkRouteArgument(route); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> insert route err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(routeStorage, s.keyGet(route), "meta")
	if err := tx.Create(keyMeta, route.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> insert route err: %s", err.Error())
		return err
	}

	keyStatus := keyCreate(routeStorage, s.keyGet(route), "status")
	if err := tx.Create(keyStatus, route.Status, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> insert route err: %s", err.Error())
		return err
	}

	keySpec := keyCreate(routeStorage, s.keyGet(route), "spec")
	if err := tx.Create(keySpec, route.Spec, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> insert route err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> insert route err: %s", err.Error())
		return err
	}

	return nil
}

// Update route info
func (s *RouteStorage) Update(ctx context.Context, route *types.Route) error {

	log.V(logLevel).Debugf("storage:etcd:route:> update route: %#v", route)

	if err := s.checkRouteExists(ctx, route); err != nil {
		return err
	}

	route.Meta.Updated = time.Now()
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> update route err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(routeStorage, s.keyGet(route), "meta")
	if err := client.Upsert(ctx, keyMeta, route.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> update route err: %s", err.Error())
		return err
	}

	return nil
}

// Remove route from storage
func (s *RouteStorage) Remove(ctx context.Context, route *types.Route) error {

	log.V(logLevel).Debugf("storage:etcd:route:> remove route: %#v", route)

	if err := s.checkRouteExists(ctx, route); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> remove err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(routeStorage, s.keyGet(route))
	if err := client.DeleteDir(ctx, key); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> remove route err: %s", err.Error())
		return err
	}

	return nil
}

// Watch route changes
func (s *RouteStorage) Watch(ctx context.Context, route chan *types.Route) error {

	log.V(logLevel).Debug("storage:etcd:route:> watch route")

	const filter = `\b\/` + routeStorage + `\/(.+):(.+)/.+\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> watch route err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(routeStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			route <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> watch route err: %s", err.Error())
		return err
	}

	return nil
}

// Watch route spec changes
func (s *RouteStorage) WatchSpec(ctx context.Context, route chan *types.Route) error {

	log.V(logLevel).Debug("storage:etcd:route:> watch route")

	const filter = `\b\/` + routeStorage + `\/(.+):(.+)/spec\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> watch route err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(routeStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			route <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> watch route err: %s", err.Error())
		return err
	}

	return nil
}

// Watch route status changes
func (s *RouteStorage) WatchStatus(ctx context.Context, route chan *types.Route) error {

	log.V(logLevel).Debug("storage:etcd:route:> watch route")

	const filter = `\b\/` + routeStorage + `\/(.+):(.+)/status\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> watch route err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(routeStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			route <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> watch route err: %s", err.Error())
		return err
	}

	return nil
}

// Clear route storage
func (s *RouteStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:route:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, routeStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:route:> clear err: %s", err.Error())
		return err
	}

	return nil
}

// keyCreate util function
func (s *RouteStorage) keyCreate(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

// keyCreate util function
func (s *RouteStorage) keyGet(r *types.Route) string {
	return r.SelfLink()
}

func newRouteStorage() *RouteStorage {
	s := new(RouteStorage)
	return s
}

// checkRouteArgument - check if argument is valid for manipulations
func (s *RouteStorage) checkRouteArgument(route *types.Route) error {

	if route == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if route.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

// checkRouteArgument - check if route exists in store
func (s *RouteStorage) checkRouteExists(ctx context.Context, route *types.Route) error {

	if err := s.checkRouteArgument(route); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:route:> check route exists")

	if _, err := s.Get(ctx, route.Meta.Namespace, route.Meta.Name); err != nil {
		log.V(logLevel).Debugf("storage:etcd:route:> check route exists err: %s", err.Error())
		return err
	}

	return nil
}
