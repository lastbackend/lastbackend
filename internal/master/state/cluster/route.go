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

package cluster

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logPrefixRoute = "observer:cluster:route"
)

func routeObserve(ss *ClusterState, d *types.Route) error {

	log.V(logLevel).Debugf("%s:> observe start: %s > %s", logPrefixRoute, d.SelfLink().String(), d.Status.State)

	switch d.Status.State {
	case types.StateCreated:
		if err := handleRouteStateCreated(ss, d); err != nil {
			log.Errorf("%s:> handle route state create err: %s", logPrefixRoute, err.Error())
			return err
		}
		break
	case types.StateProvision:
		if err := handleRouteStateProvision(ss, d); err != nil {
			log.Errorf("%s:> handle route state provision err: %s", logPrefixRoute, err.Error())
			return err
		}
		break
	case types.StateReady:
		if err := handleRouteStateReady(ss, d); err != nil {
			log.Errorf("%s:> handle route state ready err: %s", logPrefixRoute, err.Error())
			return err
		}
		break
	case types.StateError:
		if err := handleRouteStateError(ss, d); err != nil {
			log.Errorf("%s:> handle route state error err: %s", logPrefixRoute, err.Error())
			return err
		}
		break
	case types.StateDestroy:
		if err := handleRouteStateDestroy(ss, d); err != nil {
			log.Errorf("%s:> handle route state destroy err: %s", logPrefixRoute, err.Error())
			return err
		}
		break
	case types.StateDestroyed:
		if err := handleRouteStateDestroyed(ss, d); err != nil {
			log.Errorf("%s:> handle route state destroyed err: %s", logPrefixRoute, err.Error())
			return err
		}
		break
	}

	log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logPrefixRoute, d.SelfLink().String(), d.Status.State)

	return nil
}

func handleRouteStateCreated(cs *ClusterState, v *types.Route) error {
	log.V(logLevel).Debugf("%s:> handleRouteStateCreated: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)

	if err := routeProvision(cs, v); err != nil {
		return err
	}
	return nil
}

func handleRouteStateProvision(cs *ClusterState, v *types.Route) error {
	log.V(logLevel).Debugf("%s:> handleRouteStateProvision: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)

	if err := routeProvision(cs, v); err != nil {
		return err
	}
	return nil
}

func handleRouteStateReady(cs *ClusterState, v *types.Route) error {
	log.V(logLevel).Debugf("%s:> handleRouteStateReady: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)
	return nil
}

func handleRouteStateError(cs *ClusterState, v *types.Route) error {
	log.V(logLevel).Debugf("%s:> handleRouteStateError: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)
	return nil
}

func handleRouteStateDestroy(cs *ClusterState, v *types.Route) error {
	log.V(logLevel).Debugf("%s:> handleRouteStateDestroy: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)

	if err := routeDestroy(cs, v); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleRouteStateDestroyed(cs *ClusterState, v *types.Route) error {
	log.V(logLevel).Debugf("%s:> handleRouteStateDestroyed: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)

	if err := routeRemove(cs, v); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func routeUpdate(stg storage.Storage, v *types.Route, timestamp time.Time) error {

	if timestamp.Before(v.Meta.Updated) {
		vm := model.NewRouteModel(context.Background(), stg)
		if _, err := vm.Set(v); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func routeProvision(cs *ClusterState, route *types.Route) (err error) {

	t := route.Meta.Updated

	defer func() {
		if err == nil {
			err = routeUpdate(cs.storage, route, t)
		}
	}()

	rm := model.NewRouteModel(context.Background(), cs.storage)

	if route.Meta.Ingress != types.EmptyString {
		log.Debugf("%s:> route manifest provision: %s", logPrefixRoute, route.SelfLink().String())

		mf, err := rm.ManifestGet(route.Meta.Ingress, route.SelfLink().String())
		if err != nil {
			log.Errorf("%s:> route manifest create err: %s", logPrefixRoute, err.Error())
			return err
		}

		if mf != nil {
			if !routeManifestCheckEqual(mf, route) {
				if err := routeManifestSet(cs.storage, route); err != nil {
					log.Errorf("%s:> route manifest set err: %s", logPrefixRoute, err.Error())
					return err
				}
			}
		} else {
			if err := routeManifestAdd(cs.storage, route); err != nil {
				log.Errorf("%s:> route manifest add err: %s", logPrefixRoute, err.Error())
				return err
			}
		}

		if route.Status.State != types.StateProvision {
			route.Status.State = types.StateProvision
			route.Meta.Updated = time.Now()
		}

		return nil
	}

	log.Debugf("%s:> route provision > find ingress server for route: %s", logPrefixRoute, route.SelfLink().String())

	if len(cs.ingress.list) == 0 {
		route.Status.State = types.StateError
		route.Status.Message = errors.IngressNotFound
		return nil
	}

	var ing string

	for k, i := range cs.ingress.list {

		if ing == types.EmptyString {
			ing = k
		} else {
			if cs.route.ingress[ing] > cs.route.ingress[k] {
				ing = k
			}
		}

		if !i.Status.Ready {
			continue
		}

		if route.Spec.Selector.Ingress == k {
			route.Meta.Ingress = k
			break
		}
	}

	if route.Meta.Ingress == types.EmptyString && ing != types.EmptyString {
		route.Meta.Ingress = ing
	} else {
		route.Status.State = types.StateError
		route.Status.Message = errors.IngressNotFound
		return nil
	}

	mf, err := rm.ManifestGet(route.Meta.Ingress, route.SelfLink().String())
	if err != nil {
		log.Errorf("%s:> route manifest create err: %s", logPrefixRoute, err.Error())
		return err
	}

	if mf != nil {
		if !routeManifestCheckEqual(mf, route) {
			if err := routeManifestSet(cs.storage, route); err != nil {
				log.Errorf("%s:> route manifest set err: %s", logPrefixRoute, err.Error())
				return err
			}
		}
	} else {
		if err := routeManifestAdd(cs.storage, route); err != nil {
			log.Errorf("%s:> route manifest add err: %s", logPrefixRoute, err.Error())
			return err
		}
	}

	if route.Status.State != types.StateProvision {
		route.Status.State = types.StateProvision
		route.Meta.Updated = time.Now()
	}

	return nil
}

func routeDestroy(cs *ClusterState, route *types.Route) (err error) {

	t := route.Meta.Updated

	defer func() {
		if err == nil {
			err = routeUpdate(cs.storage, route, t)
		}
	}()

	if route.Spec.State == types.StateDestroy {
		if route.Meta.Ingress == types.EmptyString {
			route.Status.State = types.StateDestroyed
			route.Meta.Updated = time.Now()
			return nil
		}
	} else {
		route.Spec.State = types.StateDestroy
		route.Status.State = types.StateDestroy
		route.Meta.Updated = time.Now()
	}

	if route.Status.State != types.StateDestroy {
		route.Status.State = types.StateDestroy
		route.Meta.Updated = time.Now()
	}

	if err = routeManifestSet(cs.storage, route); err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			if route.Meta.Ingress != types.EmptyString {
				return nil
			}

			route.Status.State = types.StateDestroyed
			route.Meta.Updated = time.Now()
			return nil
		}

		return err
	}

	if route.Meta.Ingress == types.EmptyString {
		route.Status.State = types.StateDestroyed
		route.Meta.Updated = time.Now()
	}

	return nil
}

func routeRemove(cs *ClusterState, route *types.Route) (err error) {

	vm := model.NewRouteModel(context.Background(), cs.storage)
	if err = routeManifestDel(cs.storage, route); err != nil {
		return err
	}

	if err = vm.Del(route); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func routeManifestAdd(stg storage.Storage, route *types.Route) error {

	log.V(logLevel).Debugf("%s: create route manifest for node: %s", logPrefixRoute, route.SelfLink().String())

	var mf = new(types.RouteManifest)
	mf.Set(route)
	rm := model.NewRouteModel(context.Background(), stg)
	if err := rm.ManifestAdd(route.Meta.Ingress, route.SelfLink().String(), mf); err != nil {
		log.Errorf("%s:> route manifest create err: %s", logPrefixRoute, err.Error())
		return err
	}

	return nil
}

func routeManifestSet(stg storage.Storage, route *types.Route) error {

	var (
		mf  *types.RouteManifest
		err error
	)

	rm := model.NewRouteModel(context.Background(), stg)

	mf, err = rm.ManifestGet(route.Meta.Ingress, route.Meta.SelfLink.String())
	if err != nil {
		return err
	}

	// Update manifest
	if mf == nil {

		if route.Status.State != types.StateDestroy && route.Status.State != types.StateDestroyed {

			log.V(logLevel).Debugf("%s: create route manifest for ingress: %s", logPrefixRoute, route.SelfLink().String())

			mf = new(types.RouteManifest)
			mf.Set(route)

			if err := rm.ManifestAdd(route.Meta.Ingress, route.SelfLink().String(), mf); err != nil {
				log.Errorf("%s:> route manifest create err: %s", logPrefixRoute, err.Error())
				return err
			}
		}

		return nil
	}

	mf.Set(route)
	if err := rm.ManifestSet(route.Meta.Ingress, route.SelfLink().String(), mf); err != nil {
		log.Errorf("can not update route manifest: %s", err.Error())
		return err
	}

	return nil
}

func routeManifestDel(stg storage.Storage, route *types.Route) error {

	if route.Meta.Ingress == types.EmptyString {
		return nil
	}

	// Remove manifest
	rm := model.NewRouteModel(context.Background(), stg)
	err := rm.ManifestDel(route.Meta.Ingress, route.SelfLink().String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			return err
		}
	}

	return nil
}

func routeManifestCheckEqual(mf *types.RouteManifest, route *types.Route) bool {

	if mf.Endpoint != route.Spec.Endpoint {
		return false
	}

	if mf.Port != route.Spec.Port {
		return false
	}

	if len(mf.Rules) != len(route.Spec.Rules) {
		return false
	}

	for _, mr := range mf.Rules {

		var f = false

		for _, rr := range route.Spec.Rules {

			if mr.Upstream != rr.Upstream {
				continue
			}

			f = true

			if mr.Port != rr.Port {
				return false
			}

			if mr.Path != rr.Path {
				return false
			}

			if mr.Upstream != rr.Upstream {
				return false
			}

		}

		if !f {
			return false
		}
	}

	return true
}
