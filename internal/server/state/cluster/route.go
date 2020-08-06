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
//
//import (
//	"context"
//	"github.com/lastbackend/lastbackend/internal/pkg/service"
//	"github.com/lastbackend/lastbackend/tools/logger"
//	"time"
//
//	"github.com/lastbackend/lastbackend/internal/pkg/errors"
//	"github.com/lastbackend/lastbackend/internal/pkg/models"
//	"github.com/lastbackend/lastbackend/internal/pkg/storage"
//)
//
//const (
//	logPrefixRoute = "observer:cluster:route"
//)
//
//func routeObserve(ss *ClusterState, d *models.Route) error {
//	log := logger.WithContext(context.Background())
//	log.Debugf("%s:> observe start: %s > %s", logPrefixRoute, d.SelfLink().String(), d.Status.State)
//
//	switch d.Status.State {
//	case models.StateCreated:
//		if err := handleRouteStateCreated(ss, d); err != nil {
//			log.Errorf("%s:> handle route state create err: %s", logPrefixRoute, err.Error())
//			return err
//		}
//		break
//	case models.StateProvision:
//		if err := handleRouteStateProvision(ss, d); err != nil {
//			log.Errorf("%s:> handle route state provision err: %s", logPrefixRoute, err.Error())
//			return err
//		}
//		break
//	case models.StateReady:
//		if err := handleRouteStateReady(ss, d); err != nil {
//			log.Errorf("%s:> handle route state ready err: %s", logPrefixRoute, err.Error())
//			return err
//		}
//		break
//	case models.StateError:
//		if err := handleRouteStateError(ss, d); err != nil {
//			log.Errorf("%s:> handle route state error err: %s", logPrefixRoute, err.Error())
//			return err
//		}
//		break
//	case models.StateDestroy:
//		if err := handleRouteStateDestroy(ss, d); err != nil {
//			log.Errorf("%s:> handle route state destroy err: %s", logPrefixRoute, err.Error())
//			return err
//		}
//		break
//	case models.StateDestroyed:
//		if err := handleRouteStateDestroyed(ss, d); err != nil {
//			log.Errorf("%s:> handle route state destroyed err: %s", logPrefixRoute, err.Error())
//			return err
//		}
//		break
//	}
//
//	log.Debugf("%s:> observe finish: %s > %s", logPrefixRoute, d.SelfLink().String(), d.Status.State)
//
//	return nil
//}
//
//func handleRouteStateCreated(cs *ClusterState, v *models.Route) error {
//	log := logger.WithContext(context.Background())
//	log.Debugf("%s:> handleRouteStateCreated: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)
//
//	if err := routeProvision(cs, v); err != nil {
//		return err
//	}
//	return nil
//}
//
//func handleRouteStateProvision(cs *ClusterState, v *models.Route) error {
//	log := logger.WithContext(context.Background())
//	log.Debugf("%s:> handleRouteStateProvision: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)
//
//	if err := routeProvision(cs, v); err != nil {
//		return err
//	}
//	return nil
//}
//
//func handleRouteStateReady(cs *ClusterState, v *models.Route) error {
//	log := logger.WithContext(context.Background())
//	log.Debugf("%s:> handleRouteStateReady: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)
//	return nil
//}
//
//func handleRouteStateError(cs *ClusterState, v *models.Route) error {
//	log := logger.WithContext(context.Background())
//	log.Debugf("%s:> handleRouteStateError: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)
//	return nil
//}
//
//func handleRouteStateDestroy(cs *ClusterState, v *models.Route) error {
//	log := logger.WithContext(context.Background())
//	log.Debugf("%s:> handleRouteStateDestroy: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)
//
//	if err := routeDestroy(cs, v); err != nil {
//		log.Errorf("%s", err.Error())
//		return err
//	}
//
//	return nil
//}
//
//func handleRouteStateDestroyed(cs *ClusterState, v *models.Route) error {
//	log := logger.WithContext(context.Background())
//	log.Debugf("%s:> handleRouteStateDestroyed: %s > %s", logPrefixRoute, v.SelfLink().String(), v.Status.State)
//
//	if err := routeRemove(cs, v); err != nil {
//		log.Errorf("%s", err.Error())
//		return err
//	}
//
//	return nil
//}
//
//func routeUpdate(stg storage.IStorage, v *models.Route, timestamp time.Time) error {
//	log := logger.WithContext(context.Background())
//	if timestamp.Before(v.Meta.Updated) {
//		vm := service.NewRouteModel(context.Background(), stg)
//		if _, err := vm.Set(v); err != nil {
//			log.Errorf("%s", err.Error())
//			return err
//		}
//	}
//
//	return nil
//}
//
//func routeProvision(cs *ClusterState, route *models.Route) (err error) {
//	log := logger.WithContext(context.Background())
//	t := route.Meta.Updated
//
//	defer func() {
//		if err == nil {
//			err = routeUpdate(cs.storage, route, t)
//		}
//	}()
//
//	rm := service.NewRouteModel(context.Background(), cs.storage)
//
//	if route.Meta.Ingress != models.EmptyString {
//		log.Debugf("%s:> route manifest provision: %s", logPrefixRoute, route.SelfLink().String())
//
//		mf, err := rm.ManifestGet(route.Meta.Ingress, route.SelfLink().String())
//		if err != nil {
//			log.Errorf("%s:> route manifest create err: %s", logPrefixRoute, err.Error())
//			return err
//		}
//
//		if mf != nil {
//			if !routeManifestCheckEqual(mf, route) {
//				if err := routeManifestSet(cs.storage, route); err != nil {
//					log.Errorf("%s:> route manifest set err: %s", logPrefixRoute, err.Error())
//					return err
//				}
//			}
//		} else {
//			if err := routeManifestAdd(cs.storage, route); err != nil {
//				log.Errorf("%s:> route manifest add err: %s", logPrefixRoute, err.Error())
//				return err
//			}
//		}
//
//		if route.Status.State != models.StateProvision {
//			route.Status.State = models.StateProvision
//			route.Meta.Updated = time.Now()
//		}
//
//		return nil
//	}
//
//	log.Debugf("%s:> route provision > find ingress server for route: %s", logPrefixRoute, route.SelfLink().String())
//
//	if len(cs.ingress.list) == 0 {
//		route.Status.State = models.StateError
//		route.Status.Message = errors.IngressNotFound
//		return nil
//	}
//
//	var ing string
//
//	for k, i := range cs.ingress.list {
//
//		if ing == models.EmptyString {
//			ing = k
//		} else {
//			if cs.route.ingress[ing] > cs.route.ingress[k] {
//				ing = k
//			}
//		}
//
//		if !i.Status.Ready {
//			continue
//		}
//
//		if route.Spec.Selector.Ingress == k {
//			route.Meta.Ingress = k
//			break
//		}
//	}
//
//	if route.Meta.Ingress == models.EmptyString && ing != models.EmptyString {
//		route.Meta.Ingress = ing
//	} else {
//		route.Status.State = models.StateError
//		route.Status.Message = errors.IngressNotFound
//		return nil
//	}
//
//	mf, err := rm.ManifestGet(route.Meta.Ingress, route.SelfLink().String())
//	if err != nil {
//		log.Errorf("%s:> route manifest create err: %s", logPrefixRoute, err.Error())
//		return err
//	}
//
//	if mf != nil {
//		if !routeManifestCheckEqual(mf, route) {
//			if err := routeManifestSet(cs.storage, route); err != nil {
//				log.Errorf("%s:> route manifest set err: %s", logPrefixRoute, err.Error())
//				return err
//			}
//		}
//	} else {
//		if err := routeManifestAdd(cs.storage, route); err != nil {
//			log.Errorf("%s:> route manifest add err: %s", logPrefixRoute, err.Error())
//			return err
//		}
//	}
//
//	if route.Status.State != models.StateProvision {
//		route.Status.State = models.StateProvision
//		route.Meta.Updated = time.Now()
//	}
//
//	return nil
//}
//
//func routeDestroy(cs *ClusterState, route *models.Route) (err error) {
//
//	t := route.Meta.Updated
//
//	defer func() {
//		if err == nil {
//			err = routeUpdate(cs.storage, route, t)
//		}
//	}()
//
//	if route.Spec.State == models.StateDestroy {
//		if route.Meta.Ingress == models.EmptyString {
//			route.Status.State = models.StateDestroyed
//			route.Meta.Updated = time.Now()
//			return nil
//		}
//	} else {
//		route.Spec.State = models.StateDestroy
//		route.Status.State = models.StateDestroy
//		route.Meta.Updated = time.Now()
//	}
//
//	if route.Status.State != models.StateDestroy {
//		route.Status.State = models.StateDestroy
//		route.Meta.Updated = time.Now()
//	}
//
//	if err = routeManifestSet(cs.storage, route); err != nil {
//		if errors.Storage().IsErrEntityNotFound(err) {
//			if route.Meta.Ingress != models.EmptyString {
//				return nil
//			}
//
//			route.Status.State = models.StateDestroyed
//			route.Meta.Updated = time.Now()
//			return nil
//		}
//
//		return err
//	}
//
//	if route.Meta.Ingress == models.EmptyString {
//		route.Status.State = models.StateDestroyed
//		route.Meta.Updated = time.Now()
//	}
//
//	return nil
//}
//
//func routeRemove(cs *ClusterState, route *models.Route) (err error) {
//	log := logger.WithContext(context.Background())
//	vm := service.NewRouteModel(context.Background(), cs.storage)
//	if err = routeManifestDel(cs.storage, route); err != nil {
//		return err
//	}
//
//	if err = vm.Del(route); err != nil {
//		log.Errorf("%s", err.Error())
//		return err
//	}
//
//	return nil
//}
//
//func routeManifestAdd(stg storage.IStorage, route *models.Route) error {
//	log := logger.WithContext(context.Background())
//	log.Debugf("%s: create route manifest for node: %s", logPrefixRoute, route.SelfLink().String())
//
//	var mf = new(models.RouteManifest)
//	mf.Set(route)
//	rm := service.NewRouteModel(context.Background(), stg)
//	if err := rm.ManifestAdd(route.Meta.Ingress, route.SelfLink().String(), mf); err != nil {
//		log.Errorf("%s:> route manifest create err: %s", logPrefixRoute, err.Error())
//		return err
//	}
//
//	return nil
//}
//
//func routeManifestSet(stg storage.IStorage, route *models.Route) error {
//	log := logger.WithContext(context.Background())
//	var (
//		mf  *models.RouteManifest
//		err error
//	)
//
//	rm := service.NewRouteModel(context.Background(), stg)
//
//	mf, err = rm.ManifestGet(route.Meta.Ingress, route.Meta.SelfLink.String())
//	if err != nil {
//		return err
//	}
//
//	// Update manifest
//	if mf == nil {
//
//		if route.Status.State != models.StateDestroy && route.Status.State != models.StateDestroyed {
//
//			log.Debugf("%s: create route manifest for ingress: %s", logPrefixRoute, route.SelfLink().String())
//
//			mf = new(models.RouteManifest)
//			mf.Set(route)
//
//			if err := rm.ManifestAdd(route.Meta.Ingress, route.SelfLink().String(), mf); err != nil {
//				log.Errorf("%s:> route manifest create err: %s", logPrefixRoute, err.Error())
//				return err
//			}
//		}
//
//		return nil
//	}
//
//	mf.Set(route)
//	if err := rm.ManifestSet(route.Meta.Ingress, route.SelfLink().String(), mf); err != nil {
//		log.Errorf("can not update route manifest: %s", err.Error())
//		return err
//	}
//
//	return nil
//}
//
//func routeManifestDel(stg storage.IStorage, route *models.Route) error {
//
//	if route.Meta.Ingress == models.EmptyString {
//		return nil
//	}
//
//	// Remove manifest
//	rm := service.NewRouteModel(context.Background(), stg)
//	err := rm.ManifestDel(route.Meta.Ingress, route.SelfLink().String())
//	if err != nil {
//		if !errors.Storage().IsErrEntityNotFound(err) {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func routeManifestCheckEqual(mf *models.RouteManifest, route *models.Route) bool {
//
//	if mf.Endpoint != route.Spec.Endpoint {
//		return false
//	}
//
//	if mf.Port != route.Spec.Port {
//		return false
//	}
//
//	if len(mf.Rules) != len(route.Spec.Rules) {
//		return false
//	}
//
//	for _, mr := range mf.Rules {
//
//		var f = false
//
//		for _, rr := range route.Spec.Rules {
//
//			if mr.Upstream != rr.Upstream {
//				continue
//			}
//
//			f = true
//
//			if mr.Port != rr.Port {
//				return false
//			}
//
//			if mr.Path != rr.Path {
//				return false
//			}
//
//			if mr.Upstream != rr.Upstream {
//				return false
//			}
//
//		}
//
//		if !f {
//			return false
//		}
//	}
//
//	return true
//}
