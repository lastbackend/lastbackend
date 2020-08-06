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
//
//import (
//	"context"
//	"errors"
//	"time"
//
//	"github.com/lastbackend/lastbackend/internal/pkg/models"
//	"github.com/lastbackend/lastbackend/tools/log"
//)
//
//const (
//	logServicePrefix = "state:observer:service"
//)
//
//// serviceObserve manage handlers based on service state
//func serviceObserve(ss *ServiceState, s *models.Service) error {
//
//	switch s.Status.State {
//
//	// Check service created state triggers
//	case models.StateCreated:
//		if err := handleServiceStateCreated(ss, s); err != nil {
//			log.Debugf("%s:observe:serviceStateCreated err:> %s", logPrefix, err.Error())
//			return err
//		}
//		break
//
//	// Check service provision state triggers
//	case models.StateProvision:
//		if err := handleServiceStateProvision(ss, s); err != nil {
//			log.Debugf("%s:observe:serviceStateProvision err:> %s", logPrefix, err.Error())
//			return err
//		}
//		break
//
//	// Check service ready state triggers
//	case models.StateReady:
//		if err := handleServiceStateReady(ss, s); err != nil {
//			log.Debugf("%s:observe:serviceStateReady err:> %s", logPrefix, err.Error())
//			return err
//		}
//		break
//
//	// Check service error state triggers
//	case models.StateError:
//		if err := handleServiceStateError(ss, s); err != nil {
//			log.Debugf("%s:observe:serviceStateError err:> %s", logPrefix, err.Error())
//			return err
//		}
//		break
//
//		// Check service error state triggers
//	case models.StateDegradation:
//		if err := handleServiceStateDegradation(ss, s); err != nil {
//			log.Debugf("%s:observe:serviceStateDegradation err:> %s", logPrefix, err.Error())
//			return err
//		}
//		break
//
//	// Run service destroy process
//	case models.StateDestroy:
//		if err := handleServiceStateDestroy(ss, s); err != nil {
//			log.Debugf("%s:observe:serviceStateDestroy err:> %s", logPrefix, err.Error())
//			return err
//		}
//		break
//
//	// Remove service from storage if it is already destroyed
//	case models.StateDestroyed:
//		if err := handleServiceStateDestroyed(ss, s); err != nil {
//			log.Debugf("%s:observe:serviceStateDestroyed err:> %s", logPrefix, err.Error())
//			return err
//		}
//		break
//	}
//
//	if ss.service == nil {
//		return nil
//	}
//
//	log.Debugf("%s:> observe handle: %s > %s", logServicePrefix, s.SelfLink(), s.Status.State)
//
//	ss.service = s
//	if err := serviceStatusState(ss); err != nil {
//		return err
//	}
//
//	log.Debugf("%s:> observe finish: %s > %s", logServicePrefix, s.SelfLink(), s.Status.State)
//
//	return nil
//}
//
//// handleServiceStateCreated handles service created state
//func handleServiceStateCreated(ss *ServiceState, svc *models.Service) error {
//
//	log.Debugf("%s:> handleServiceStateCreated: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)
//
//	// Endpoint provision call
//	if err := serviceEndpointProvision(ss, svc); err != nil {
//		log.Errorf("%s:> endpoint provision err: %s", logServicePrefix, err.Error())
//		return err
//	}
//
//	// Deployment provision call
//	if err := serviceDeploymentProvision(ss, svc); err != nil {
//		log.Errorf("%s:> deployment provision err: %s", logServicePrefix, err.Error())
//		return err
//	}
//
//	return nil
//}
//
//// handleServiceStateProvision handles service provision state
//func handleServiceStateProvision(ss *ServiceState, svc *models.Service) error {
//
//	log.Debugf("%s:> handleServiceStateProvision: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)
//
//	// Endpoint provision call
//	if err := serviceEndpointProvision(ss, svc); err != nil {
//		log.Errorf("%s:> endpoint provision err: %s", logServicePrefix, err.Error())
//		return err
//	}
//
//	// Deployment provision call
//	if err := serviceDeploymentProvision(ss, svc); err != nil {
//		log.Errorf("%s:> deployment provision err: %s", logServicePrefix, err.Error())
//		return err
//	}
//
//	return nil
//}
//
//// handleServiceStateReady handles service ready state
//func handleServiceStateReady(ss *ServiceState, svc *models.Service) error {
//
//	log.Debugf("%s:> handleServiceStateReady: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)
//
//	return nil
//}
//
//// handleServiceStateError handles service error state
//func handleServiceStateError(ss *ServiceState, svc *models.Service) error {
//
//	log.Debugf("%s:> handleServiceStateError: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)
//
//	return nil
//}
//
//// handleServiceStateDegradation handles service degradation state
//func handleServiceStateDegradation(ss *ServiceState, svc *models.Service) error {
//
//	log.Debugf("%s:> handleServiceStateDegradation: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)
//
//	return nil
//}
//
//// handleServiceStateDestroy handles service destroy state
//func handleServiceStateDestroy(ss *ServiceState, svc *models.Service) (err error) {
//
//	log.Debugf("%s:> handleServiceStateDestroy: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)
//
//	if ss.endpoint.endpoint != nil {
//		if err = endpointDel(ss); err != nil {
//			log.Errorf("%s:> endpoint remove err: %s", logServicePrefix, err.Error())
//			return err
//		}
//	}
//
//	dm := service.NewDeploymentModel(context.Background(), ss.storage)
//
//	for _, d := range ss.deployment.list {
//
//		if d.Status.State == models.StateDestroyed {
//			continue
//		}
//
//		if d.Status.State != models.StateDestroy {
//			if err := dm.Destroy(d); err != nil {
//				return err
//			}
//		}
//	}
//
//	if len(ss.deployment.list) == 0 {
//		svc.Status.State = models.StateDestroyed
//		svc.Meta.Updated = time.Now()
//	}
//
//	return nil
//}
//
//// handleServiceStateDestroyed handles service destroyed state
//func handleServiceStateDestroyed(ss *ServiceState, svc *models.Service) (err error) {
//
//	log.Debugf("%s:> handleServiceStateDestroyed: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)
//
//	if err = endpointDel(ss); err != nil {
//		log.Errorf("%s:> endpoint remove err: %s", logServicePrefix, err.Error())
//		return err
//	}
//
//	svc.Status.State = models.StateDestroy
//	svc.Meta.Updated = time.Now()
//
//	if len(ss.deployment.list) > 0 {
//		dm := service.NewDeploymentModel(context.Background(), ss.storage)
//		for _, d := range ss.deployment.list {
//
//			if d.Status.State == models.StateDestroyed {
//				if err = dm.Remove(d); err != nil {
//					return err
//				}
//			}
//
//			if d.Status.State != models.StateDestroy {
//				if err = dm.Destroy(d); err != nil {
//					return err
//				}
//			}
//
//		}
//
//		svc.Status.State = models.StateDestroy
//		svc.Meta.Updated = time.Now()
//
//		return nil
//	}
//
//	sm := service.NewServiceModel(context.Background(), ss.storage)
//	nm := service.NewNamespaceModel(context.Background(), ss.storage)
//
//	ns, err := nm.Get(svc.Meta.Namespace)
//	if err != nil {
//		log.Errorf("%s:> namespece fetch err: %s", logServicePrefix, err.Error())
//	}
//
//	if ns != nil {
//		ns.ReleaseResources(svc.Spec.GetResourceRequest())
//
//		if err := nm.Update(ns); err != nil {
//			log.Errorf("%s:> namespece update err: %s", logServicePrefix, err.Error())
//		}
//	}
//
//	if err = sm.Remove(svc); err != nil {
//		log.Errorf("%s:> service remove err: %s", logServicePrefix, err.Error())
//		return err
//	}
//
//	ss.service = nil
//	return nil
//}
//
//// serviceEndpointProvision function handles all cases for endpoint management
//func serviceEndpointProvision(ss *ServiceState, svc *models.Service) error {
//
//	// check endpoint is needed
//	if svc == nil {
//		log.Errorf("%s:> need pointer as svc argument", logServicePrefix)
//		return errors.New("svc is nil")
//	}
//
//	if err := endpointProvision(ss, svc); err != nil {
//		log.Errorf("%s:> endpoint provision err: %s", logServicePrefix, err.Error())
//		return err
//	}
//
//	if ss.endpoint.endpoint != nil {
//		svc.Meta.Endpoint = ss.endpoint.endpoint.Spec.Domain
//		svc.Meta.IP = ss.endpoint.endpoint.Spec.IP
//	}
//
//	return nil
//}
//
//// serviceDeploymentProvision function handles all cases when deployment needs to be created or updated
//func serviceDeploymentProvision(ss *ServiceState, svc *models.Service) error {
//
//	var (
//		d *models.Deployment
//	)
//
//	// select deployment for provision
//	switch true {
//
//	// check provision deployment exists and is current for service
//	case ss.deployment.provision != nil:
//		d = ss.deployment.provision
//		break
//
//	// check active deployment exists and is current for service
//	case ss.deployment.active != nil:
//		d = ss.deployment.active
//		break
//	}
//
//	// if deployment found for provision: check and update replicas
//	if d != nil {
//		if d.Spec.Replicas != svc.Spec.Replicas {
//			if err := deploymentScale(ss.storage, d, svc.Spec.Replicas); err != nil {
//				log.Errorf("%s:> deployment scale err: %s", logServicePrefix, err.Error())
//				return err
//			}
//		}
//	}
//
//	// create deployment if needed
//	if d == nil {
//
//		if len(svc.Spec.Template.Containers) == 0 {
//			svc.Status.State = models.StateReady
//			return nil
//		}
//
//		ss.deployment.index++
//
//		d, err := deploymentCreate(ss.storage, svc, ss.deployment.index)
//		if err != nil {
//			ss.deployment.index--
//			log.Errorf("%s:> deployment create err: %s", logServicePrefix, err.Error())
//			return err
//		}
//
//		for _, od := range ss.deployment.list {
//
//			if ss.deployment.active != nil {
//				if ss.deployment.active.SelfLink().String() == od.SelfLink().String() && od.Status.State == models.StateReady {
//					continue
//				}
//			}
//
//			if od.Status.State != models.StateDestroy && od.Status.State != models.StateDestroyed {
//				if err := deploymentDestroy(ss, od); err != nil {
//					log.Errorf("%s:> deployment cancel err: %s", logServicePrefix, err.Error())
//					return err
//				}
//			}
//		}
//
//		ss.deployment.list[d.SelfLink().String()] = d
//		ss.deployment.provision = d
//	}
//
//	return nil
//}
//
//// serviceStatusState calculates current service status based on deployments
//func serviceStatusState(ss *ServiceState) (err error) {
//
//	status := ss.service.Status
//
//	defer func() {
//
//		if status.State == ss.service.Status.State && status.Message == ss.service.Status.Message {
//			return
//		}
//
//		ss.service.Meta.Updated = time.Now()
//		sm := service.NewServiceModel(context.Background(), ss.storage)
//		if err := sm.Set(ss.service); err != nil {
//			log.Errorf("%s", err.Error())
//			return
//		}
//
//		return
//	}()
//
//	if ss.service.Status.State == models.StateDestroyed {
//		sm := service.NewServiceModel(context.Background(), ss.storage)
//		if err := sm.Set(ss.service); err != nil {
//			log.Errorf("%s", err.Error())
//			return err
//		}
//		return nil
//	}
//
//	if ss.service.Status.State == models.StateDestroy {
//
//		if len(ss.deployment.list) == 0 {
//			ss.service.Status.State = models.StateDestroyed
//		}
//
//		return nil
//	}
//
//	if ss.service.Status.State == models.StateProvision || ss.service.Status.State == models.StateCreated {
//
//		if ss.deployment.provision == nil && ss.deployment.active != nil {
//
//			ss.service.Status.State = ss.deployment.active.Status.State
//			ss.service.Status.Message = ss.deployment.active.Status.Message
//		}
//
//		if ss.deployment.provision != nil {
//
//			ss.service.Status.State = ss.deployment.provision.Status.State
//			ss.service.Status.Message = ss.deployment.provision.Status.Message
//
//		}
//
//		return nil
//	}
//
//	if ss.deployment.active != nil {
//
//		ss.service.Status.State = ss.deployment.active.Status.State
//		ss.service.Status.Message = ss.deployment.active.Status.Message
//
//		if ss.deployment.active.Status.State == models.StateCreated {
//
//			ss.service.Status.State = models.StateProvision
//			ss.service.Status.Message = models.EmptyString
//		}
//
//		if ss.deployment.provision != nil && ss.deployment.provision.Status.State == models.StateProvision {
//			ss.service.Status.State = ss.deployment.provision.Status.State
//			ss.service.Status.Message = ss.deployment.provision.Status.Message
//		}
//	}
//
//	if ss.deployment.provision != nil {
//		ss.service.Status.State = ss.deployment.provision.Status.State
//		ss.service.Status.Message = ss.deployment.provision.Status.Message
//	}
//
//	return nil
//}
