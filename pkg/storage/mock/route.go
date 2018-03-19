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

package mock

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"strings"
)

type RouteStorage struct {
	storage.Route
	data map[string]*types.Route
}

// Get route by name
func (s *RouteStorage) Get(ctx context.Context, namespace, name string) (*types.Route, error) {
	if ns, ok := s.data[s.keyCreate(namespace, name)]; ok {
		return ns, nil
	}
	return nil, errors.New(store.ErrEntityNotFound)
}

// Get routes by namespace name
func (s *RouteStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Route, error) {
	list := make(map[string]*types.Route, 0)

	prefix := fmt.Sprintf("%s:", namespace)
	for _, d := range s.data {

		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// Update route state
func (s *RouteStorage) SetState(ctx context.Context, route *types.Route) error {
	if err := s.checkRouteExists(route); err != nil {
		return err
	}

	s.data[s.keyGet(route)].State = route.State
	return nil
}

// Insert new route
func (s *RouteStorage) Insert(ctx context.Context, route *types.Route) error {

	if err := s.checkRouteArgument(route); err != nil {
		return err
	}

	s.data[s.keyGet(route)] = route

	return nil
}

// Update route info
func (s *RouteStorage) Update(ctx context.Context, route *types.Route) error {

	if err := s.checkRouteExists(route); err != nil {
		return err
	}

	s.data[s.keyGet(route)] = route

	return nil
}

// Remove route from storage
func (s *RouteStorage) Remove(ctx context.Context, route *types.Route) error {

	if err := s.checkRouteExists(route); err != nil {
		return err
	}

	delete(s.data, s.keyGet(route))

	return nil
}

// Watch route changes
func (s *RouteStorage) Watch(ctx context.Context, route chan *types.Route) error {
	return nil
}

// Watch route spec changes
func (s *RouteStorage) WatchSpec(ctx context.Context, route chan *types.Route) error {
	return nil
}

// Watch route state changes
func (s *RouteStorage) WatchState(ctx context.Context, route chan *types.Route) error {
	return nil
}

// Clear route storage
func (s *RouteStorage) Clear(ctx context.Context) error {
	s.data = make(map[string]*types.Route)
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

// newRouteStorage returns new storage
func newRouteStorage() *RouteStorage {
	s := new(RouteStorage)
	s.data = make(map[string]*types.Route)
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
func (s *RouteStorage) checkRouteExists(route *types.Route) error {

	if err := s.checkRouteArgument(route); err != nil {
		return err
	}

	if _, ok := s.data[s.keyGet(route)]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}
