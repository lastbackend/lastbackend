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

package route
//
//import (
//	"context"
//	"fmt"
//	"net/http"
//	"strings"
//
//	"github.com/lastbackend/lastbackend/internal/api/envs"
//	"github.com/lastbackend/lastbackend/internal/pkg/errors"
//	"github.com/lastbackend/lastbackend/internal/pkg/models"
//	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
//	"github.com/lastbackend/lastbackend/tools/log"
//)
//
//const (
//	logPrefix = "api:handler:route"
//	logLevel  = 3
//)
//
//func Fetch(ctx context.Context, namespace, name string) (*models.Route, *errors.Err) {
//
//	vm := service.NewRouteModel(ctx, envs.Get().GetStorage())
//	vol, err := vm.Get(namespace, name)
//
//	if err != nil {
//		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
//		return nil, errors.New("route").InternalServerError(err)
//	}
//
//	if vol == nil {
//		err := errors.New("namespace not found")
//		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
//		return nil, errors.New("route").NotFound()
//	}
//
//	return vol, nil
//}
//
//func Apply(ctx context.Context, ns *models.Namespace, mf *request.RouteManifest) (*models.Route, *errors.Err) {
//
//	if mf.Meta.Name == nil {
//		return nil, errors.New("route").BadParameter("meta.name")
//	}
//
//	vol, err := Fetch(ctx, ns.Meta.Name, *mf.Meta.Name)
//	if err != nil {
//		if err.Code != http.StatusText(http.StatusNotFound) {
//			return nil, errors.New("route").InternalServerError()
//		}
//	}
//
//	if vol == nil {
//		return Create(ctx, ns, mf)
//	}
//
//	return Update(ctx, ns, vol, mf)
//}
//
//func Create(ctx context.Context, ns *models.Namespace, mf *request.RouteManifest) (*models.Route, *errors.Err) {
//
//	rm := service.NewRouteModel(ctx, envs.Get().GetStorage())
//	sm := service.NewServiceModel(ctx, envs.Get().GetStorage())
//
//	if mf.Meta.Name != nil {
//
//		route, err := rm.Get(ns.Meta.Name, *mf.Meta.Name)
//		if err != nil {
//			log.Errorf("%s:create:> get route by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
//			return nil, errors.New("route").InternalServerError()
//		}
//
//		if route != nil {
//			log.Warnf("%s:create:> route name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
//			return nil, errors.New("route").NotUnique("name")
//		}
//	}
//
//	if err := validateManifest(ctx, mf); err != nil {
//		log.Errorf("%s:create:> route manifest validation err", logPrefix, err.Err().Error())
//		return nil, err
//	}
//
//	svc, err := sm.List(ns.Meta.Name)
//	if err != nil {
//		log.Errorf("%s:create:> get services", logPrefix, err.Error())
//		return nil, errors.New("route").InternalServerError()
//	}
//
//	route := new(models.Route)
//	route.Meta.SetDefault()
//	route.Meta.SelfLink = *models.NewRouteSelfLink(ns.Meta.Name, *mf.Meta.Name)
//	route.Meta.Namespace = ns.Meta.Name
//
//	mf.SetRouteMeta(route)
//	mf.SetRouteSpec(route, ns, svc)
//
//	if route.Spec.Endpoint == models.EmptyString {
//		_, external := envs.Get().GetDomain()
//		route.Spec.Endpoint = fmt.Sprintf("%s.%s.%s", strings.ToLower(route.Meta.Name), strings.ToLower(ns.Meta.Name), external)
//	}
//
//	if len(route.Spec.Rules) == 0 {
//		err := errors.New("route rules are incorrect")
//		log.Errorf("%s:create:> route rules empty: %s", logPrefix, err.Error())
//		return nil, errors.New("route").BadParameter("rules", err)
//	}
//
//	if _, err := rm.Add(ns, route); err != nil {
//		log.Errorf("%s:create:> create route err: %s", logPrefix, ns.Meta.Name, err.Error())
//		return nil, errors.New("route").InternalServerError()
//	}
//
//	return route, nil
//}
//
////
//func Update(ctx context.Context, ns *models.Namespace, rt *models.Route, mf *request.RouteManifest) (*models.Route, *errors.Err) {
//
//	rm := service.NewRouteModel(ctx, envs.Get().GetStorage())
//	sm := service.NewServiceModel(ctx, envs.Get().GetStorage())
//
//	if mf.Meta.Name != nil {
//
//		route, err := rm.Get(ns.Meta.Name, *mf.Meta.Name)
//		if err != nil {
//			log.Errorf("%s:create:> get route by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
//			return nil, errors.New("route").InternalServerError()
//		}
//
//		if route == nil {
//			log.Warnf("%s:create:> route name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
//			return nil, errors.New("route").NotFound()
//		}
//	}
//
//	if err := validateManifest(ctx, mf); err != nil {
//		log.Errorf("%s:update:> route manifest validation err: %s", logPrefix, err.Err().Error())
//		return nil, err
//	}
//
//	svc, err := sm.List(ns.Meta.Name)
//	if err != nil {
//		log.Errorf("%s:create:> get services err: %s", logPrefix, err.Error())
//		return nil, errors.New("route").InternalServerError()
//	}
//
//	mf.SetRouteMeta(rt)
//	mf.SetRouteSpec(rt, ns, svc)
//
//	if rt.Spec.Endpoint == models.EmptyString {
//		_, external := envs.Get().GetDomain()
//		rt.Spec.Endpoint = fmt.Sprintf("%s.%s.%s", strings.ToLower(rt.Meta.Name), strings.ToLower(ns.Meta.Name), external)
//	}
//
//	if len(rt.Spec.Rules) == 0 {
//		err := errors.New("route rules are incorrect")
//		log.Errorf("%s:create:> route rules empty: %s", logPrefix, err.Error())
//		return nil, errors.New("route").BadParameter("rules", err)
//	}
//
//	rt.Status.State = models.StateProvision
//	rt, err = rm.Set(rt)
//	if err != nil {
//		log.Errorf("%s:update:> update route `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
//		return nil, errors.New("route").InternalServerError()
//	}
//
//	return rt, nil
//}
//
//func validateManifest(ctx context.Context, mf *request.RouteManifest) *errors.Err {
//
//	rm := service.NewRouteModel(ctx, envs.Get().GetStorage())
//
//	rl, err := rm.List()
//	if err != nil {
//		log.Errorf("%s:validate:> route manifest validation failed: %s ", logPrefix, err.Error())
//		return errors.New("route").InternalServerError()
//	}
//
//	if mf.Spec.Port != 80 && mf.Spec.Port != 443 {
//		for _, r := range rl.Items {
//			if r.Spec.Port == mf.Spec.Port {
//				return errors.New("route").Allocated("port", errors.Route().NewErrPortAllocated())
//			}
//		}
//	}
//
//	// TODO:  check this. If we want to update route, we need to validate it but endpoint is always allocated
//	//if mf.Spec.Endpoint != types.EmptyString {
//	//	for _, r := range rl.Items {
//	//		if r.Spec.Endpoint == mf.Spec.Endpoint {
//	//			return errors.New("route").Allocated("endpoint", errors.Route().NewErrEndpointAllocated())
//	//		}
//	//	}
//	//}
//
//	return nil
//}
