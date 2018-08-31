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

package service

import (
	"context"
	"errors"

	"time"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	logServicePrefix = "state:observer:service"
)

// serviceObserve manage handlers based on service state
func serviceObserve(ss *ServiceState, s *types.Service) error {

	log.V(logLevel).Debugf("%s:> observe start: %s > %s", logServicePrefix, s.SelfLink(), s.Status.State)

	switch s.Status.State {

	// Check service created state triggers
	case types.StateCreated:
		if err := handleServiceStateCreated(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:serviceStateCreated err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Check service provision state triggers
	case types.StateProvision:
		if err := handleServiceStateProvision(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:serviceStateProvision err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Check service ready state triggers
	case types.StateReady:
		if err := handleServiceStateReady(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:serviceStateReady err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Check service error state triggers
	case types.StateError:
		if err := handleServiceStateError(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:serviceStateError err:> %s", logPrefix, err.Error())
			return err
		}
		break

		// Check service error state triggers
	case types.StateDegradation:
		if err := handleServiceStateDegradation(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:serviceStateDegradation err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Run service destroy process
	case types.StateDestroy:
		if err := handleServiceStateDestroy(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:serviceStateDestroy err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Remove service from storage if it is already destroyed
	case types.StateDestroyed:
		if err := handleServiceStateDestroyed(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:serviceStateDestroyed err:> %s", logPrefix, err.Error())
			return err
		}
		break
	}

	if ss.service == nil {
		return nil
	}

	ss.service = s
	if err := serviceStatusState(ss); err != nil {
		return err
	}

	log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logServicePrefix, s.SelfLink(), s.Status.State)

	return nil
}

// handleServiceStateCreated handles service created state
func handleServiceStateCreated(ss *ServiceState, svc *types.Service) error {

	log.V(logLevel).Debugf("%s:> handleServiceStateCreated: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)

	// Endpoint provision call
	if err := serviceEndpointProvision(ss, svc); err != nil {
		log.Errorf("%s:> endpoint provision err: %s", logServicePrefix, err.Error())
		return err
	}

	// Deployment provision call
	if err := serviceDeploymentProvision(ss, svc); err != nil {
		log.Errorf("%s:> deployment provision err: %s", logServicePrefix, err.Error())
		return err
	}

	return nil
}

// handleServiceStateProvision handles service provision state
func handleServiceStateProvision(ss *ServiceState, svc *types.Service) error {

	log.V(logLevel).Debugf("%s:> handleServiceStateProvision: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)

	// Endpoint provision call
	if err := serviceEndpointProvision(ss, svc); err != nil {
		log.Errorf("%s:> endpoint provision err: %s", logServicePrefix, err.Error())
		return err
	}

	// Deployment provision call
	if err := serviceDeploymentProvision(ss, svc); err != nil {
		log.Errorf("%s:> deployment provision err: %s", logServicePrefix, err.Error())
		return err
	}

	return nil
}

// handleServiceStateReady handles service ready state
func handleServiceStateReady(ss *ServiceState, svc *types.Service) error {

	log.V(logLevel).Debugf("%s:> handleServiceStateReady: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)

	return nil
}

// handleServiceStateError handles service error state
func handleServiceStateError(ss *ServiceState, svc *types.Service) error {

	log.V(logLevel).Debugf("%s:> handleServiceStateError: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)

	return nil
}

// handleServiceStateDegradation handles service degradation state
func handleServiceStateDegradation(ss *ServiceState, svc *types.Service) error {

	log.V(logLevel).Debugf("%s:> handleServiceStateDegradation: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)

	return nil
}

// handleServiceStateDestroy handles service destroy state
func handleServiceStateDestroy(ss *ServiceState, svc *types.Service) (err error) {

	log.V(logLevel).Debugf("%s:> handleServiceStateDestroy: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)

	if ss.endpoint.endpoint != nil {
		if err = endpointDel(ss); err != nil {
			log.Errorf("%s:> endpoint remove err: %s", logServicePrefix, err.Error())
			return err
		}
	}

	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
	if len(ss.deployment.list) == 0 {
		sm := distribution.NewServiceModel(context.Background(), envs.Get().GetStorage())
		if err = sm.Remove(svc); err != nil {
			log.Errorf("%s:> service remove err: %s", logServicePrefix, err.Error())
			return err
		}

		ss.service = nil
		return nil
	}

	for _, d := range ss.deployment.list {

		if d.Status.State == types.StateDestroyed {
			continue
		}

		if d.Status.State != types.StateDestroy {
			if err := dm.Destroy(d); err != nil {
				return err
			}
		}
	}

	if len(ss.deployment.list) == 0 {
		svc.Status.State = types.StateDestroyed
		svc.Meta.Updated = time.Now()
	}

	return nil
}

// handleServiceStateDestroyed handles service destroyed state
func handleServiceStateDestroyed(ss *ServiceState, svc *types.Service) (err error) {

	log.V(logLevel).Debugf("%s:> handleServiceStateDestroyed: %s > %s", logServicePrefix, svc.SelfLink(), svc.Status.State)

	if err = endpointDel(ss); err != nil {
		log.Errorf("%s:> endpoint remove err: %s", logServicePrefix, err.Error())
		return err
	}

	svc.Status.State = types.StateDestroy
	svc.Meta.Updated = time.Now()

	if len(ss.deployment.list) > 0 {
		dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
		for _, d := range ss.deployment.list {

			if d.Status.State == types.StateDestroyed {
				if err = dm.Remove(d); err != nil {
					return err
				}
			}

			if d.Status.State != types.StateDestroy {
				if err = dm.Destroy(d); err != nil {
					return err
				}
			}

		}

		svc.Status.State = types.StateDestroy
		svc.Meta.Updated = time.Now()

		return nil
	}

	sm := distribution.NewServiceModel(context.Background(), envs.Get().GetStorage())
	if err = sm.Remove(svc); err != nil {
		log.Errorf("%s:> service remove err: %s", logServicePrefix, err.Error())
		return err
	}

	ss.service = nil
	return nil
}

// serviceEndpointProvision function handles all cases for endpoint management
func serviceEndpointProvision(ss *ServiceState, svc *types.Service) error {

	// check endpoint is needed
	if svc == nil {
		log.Errorf("%s:> need pointer as svc argument", logServicePrefix)
		return errors.New("svc is nil")
	}

	if err := endpointProvision(ss, svc); err != nil {
		log.Errorf("%s:> endpoint provision err: %s", logServicePrefix, err.Error())
		return err
	}

	if err := endpointManifestProvision(ss); err != nil {
		log.Errorf("%s:> endpoint provision err: %s", logServicePrefix, err.Error())
		return err
	}

	if ss.endpoint.endpoint != nil {
		svc.Meta.Endpoint = ss.endpoint.endpoint.Spec.Domain
		svc.Meta.IP = ss.endpoint.endpoint.Spec.IP
	}

	return nil
}

// serviceDeploymentProvision function handles all cases when deployment needs to be created or updated
func serviceDeploymentProvision(ss *ServiceState, svc *types.Service) error {

	var (
		d *types.Deployment
	)

	// select deployment for provision
	switch true {

	// check provision deployment exists and is current for service
	case ss.deployment.provision != nil && deploymentSpecValidate(ss.deployment.provision, svc.Spec.Template):
		d = ss.deployment.provision
		break

	// check active deployment exists and is current for service
	case ss.deployment.active != nil && deploymentSpecValidate(ss.deployment.active, svc.Spec.Template):
		d = ss.deployment.active
		break
	}

	// if deployment found for provision: check and update replicas
	if d != nil {
		if d.Spec.Replicas != svc.Spec.Replicas {
			if err := deploymentScale(d, svc.Spec.Replicas); err != nil {
				log.Errorf("%s:> deployment scale err: %s", logServicePrefix, err.Error())
				return err
			}
		}
	}

	// create deployment if needed
	if d == nil {

		d, err := deploymentCreate(svc)
		if err != nil {
			log.Errorf("%s:> deployment create err: %s", logServicePrefix, err.Error())
			return err
		}

		for _, od := range ss.deployment.list {

			if ss.deployment.active != nil {
				if ss.deployment.active.SelfLink() == od.SelfLink() && od.Status.State == types.StateReady {
					continue
				}
			}

			if od.Status.State != types.StateDestroy && od.Status.State != types.StateDestroyed {
				if err := deploymentDestroy(ss, od); err != nil {
					log.Errorf("%s:> deployment cancel err: %s", logServicePrefix, err.Error())
					return err
				}
			}
		}

		ss.deployment.list[d.SelfLink()] = d
		ss.deployment.provision = d
	}

	return nil
}

// serviceStatusState calculates current service status based on deployments
func serviceStatusState(ss *ServiceState) (err error) {

	status := ss.service.Status

	defer func() error {
		if status.State == ss.service.Status.State && status.Message == ss.service.Status.Message {
			return nil
		}

		ss.service.Meta.Updated = time.Now()
		sm := distribution.NewServiceModel(context.Background(), envs.Get().GetStorage())
		if err = sm.Set(ss.service); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}

		return nil
	}()

	if ss.service.Status.State == types.StateProvision || ss.service.Status.State == types.StateCreated {

		if ss.deployment.active != nil {
			if deploymentSpecValidate(ss.deployment.active, ss.service.Spec.Template) &&
				ss.deployment.active.Spec.Replicas == ss.service.Spec.Replicas {
				ss.service.Status.State = ss.deployment.active.Status.State
				ss.service.Status.Message = ss.deployment.active.Status.Message
				if ss.deployment.active.Status.State == types.StateCreated {
					ss.service.Status.State = types.StateProvision
					ss.service.Status.Message = types.EmptyString
				}
			}
		}

		if ss.deployment.provision != nil && ss.deployment.active == nil {
			if deploymentSpecValidate(ss.deployment.provision, ss.service.Spec.Template) &&
				ss.deployment.provision.Spec.Replicas == ss.service.Spec.Replicas {
				ss.service.Status.State = ss.deployment.provision.Status.State
				ss.service.Status.Message = ss.deployment.provision.Status.Message
			}
		}

		return nil
	}

	if ss.service.Status.State == types.StateDestroy {

		if len(ss.deployment.list) == 0 {
			ss.service.Status.State = types.StateDestroyed
		}

		return nil
	}

	if ss.deployment.provision != nil {
		ss.service.Status.State = ss.deployment.provision.Status.State
		ss.service.Status.Message = ss.deployment.provision.Status.Message
	}

	if ss.deployment.active != nil {

		ss.service.Status.State = ss.deployment.active.Status.State
		ss.service.Status.Message = ss.deployment.active.Status.Message

		if ss.deployment.active.Status.State == types.StateCreated {
			ss.service.Status.State = types.StateProvision
			ss.service.Status.Message = types.EmptyString
		}

		if ss.deployment.provision != nil && ss.deployment.provision.Status.State == types.StateProvision {
			ss.service.Status.State = ss.deployment.provision.Status.State
			ss.service.Status.Message = ss.deployment.provision.Status.Message
		}
	}

	return nil
}
