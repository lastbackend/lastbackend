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

package route

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"net/http"
)

const (
	logPrefix = "api:handler:route"
	logLevel  = 3
)

func Fetch(ctx context.Context, namespace, name string) (*types.Route, *errors.Err) {

	vm := distribution.NewRouteModel(ctx, envs.Get().GetStorage())
	vol, err := vm.Get(namespace, name)

	if err != nil {
		log.V(logLevel).Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("route").InternalServerError(err)
	}

	if vol == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("route").NotFound()
	}

	return vol, nil
}

func Apply(ctx context.Context, ns *types.Namespace, mf *request.RouteManifest) (*types.Route, *errors.Err) {

	if mf.Meta.Name == nil {
		return nil, errors.New("route").BadParameter("meta.name")
	}

	vol, err := Fetch(ctx, ns.Meta.Name, *mf.Meta.Name)
	if err != nil {
		if err.Code != http.StatusText(http.StatusNotFound) {
			return nil, errors.New("route").InternalServerError()
		}
	}

	if vol == nil {
		return Create(ctx, ns, mf)
	}

	return Update(ctx, ns, vol, mf)
}

func Create(ctx context.Context, ns *types.Namespace, mf *request.RouteManifest) (*types.Route, *errors.Err) {

	rm := distribution.NewRouteModel(ctx, envs.Get().GetStorage())
	sm := distribution.NewServiceModel(ctx, envs.Get().GetStorage())

	if mf.Meta.Name != nil {

		route, err := rm.Get(ns.Meta.Name, *mf.Meta.Name)
		if err != nil {
			log.V(logLevel).Errorf("%s:create:> get route by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
			return nil, errors.New("route").InternalServerError()
		}

		if route != nil {
			log.V(logLevel).Warnf("%s:create:> route name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
			return nil, errors.New("route").NotUnique("name")
		}
	}

	svc, err := sm.List(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> get services", logPrefix, err.Error())
		return nil, errors.New("route").InternalServerError()
	}

	route := new(types.Route)
	route.Meta.SetDefault()
	route.Meta.Namespace = ns.Meta.Name

	mf.SetRouteMeta(route)
	mf.SetRouteSpec(route, svc)

	if len(route.Spec.Rules) == 0 {
		err := errors.New("route rules are incorrect")
		log.V(logLevel).Errorf("%s:create:> route rules empty", logPrefix, err.Error())
		return nil, errors.New("route").BadParameter("rules", err)

	}

	if _, err := rm.Add(ns, route); err != nil {
		log.V(logLevel).Errorf("%s:create:> create route err: %s", logPrefix, ns.Meta.Name, err.Error())
		return nil, errors.New("route").InternalServerError()
	}

	return route, nil
}

//
func Update(ctx context.Context, ns *types.Namespace, rt *types.Route, mf *request.RouteManifest) (*types.Route, *errors.Err) {

	rm := distribution.NewRouteModel(ctx, envs.Get().GetStorage())
	sm := distribution.NewServiceModel(ctx, envs.Get().GetStorage())

	if mf.Meta.Name != nil {

		route, err := rm.Get(ns.Meta.Name, *mf.Meta.Name)
		if err != nil {
			log.V(logLevel).Errorf("%s:create:> get route by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
			return nil, errors.New("route").InternalServerError()
		}

		if route == nil {
			log.V(logLevel).Warnf("%s:create:> route name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
			return nil, errors.New("route").NotFound()
		}
	}

	svc, err := sm.List(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> get services", logPrefix, err.Error())
		return nil, errors.New("route").InternalServerError()
	}

	mf.SetRouteMeta(rt)
	mf.SetRouteSpec(rt, svc)

	if len(rt.Spec.Rules) == 0 {
		err := errors.New("route rules are incorrect")
		log.V(logLevel).Errorf("%s:create:> route rules empty", logPrefix, err.Error())
		return nil, errors.New("route").BadParameter("rules", err)
	}

	rt, err = rm.Set(rt)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> update route `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
		return nil, errors.New("route").InternalServerError()
	}

	return rt, nil
}
