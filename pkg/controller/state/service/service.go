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

	"github.com/lastbackend/dynamic/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	)

const (
	logServicePrefix = "state:observer:service"
)

// serviceObserve manage handlers based on service state
func serviceObserve(ss *ServiceState, s *types.Service) error {
	switch s.Status.State {

	// Check service created state triggers
	case types.StateCreated:
		if err := handleServiceStateCreated(ss, s); err != nil {
			log.Debugf("%s:observe:serviceStateCreated err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Check service provision state triggers
	case types.StateProvision:
		if err := handleServiceStateProvision(ss, s); err != nil {
			log.Debugf("%s:observe:serviceStateProvision err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Check service ready state triggers
	case types.StateReady:
		if err := handleServiceStateReady(ss, s); err != nil {
			log.Debugf("%s:observe:serviceStateReady err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Check service error state triggers
	case types.StateError:
		if err := handleServiceStateError(ss, s); err != nil {
			log.Debugf("%s:observe:serviceStateError err:> %s", logPrefix, err.Error())
			return err
		}
		break

		// Check service error state triggers
	case types.StateDegradation:
		if err := handleServiceStateDegradation(ss, s); err != nil {
			log.Debugf("%s:observe:serviceStateDegradation err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Run service destroy process
	case types.StateDestroy:
		if err := handleServiceStateDestroy(ss, s); err != nil {
			log.Debugf("%s:observe:serviceStateDestroy err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Remove service from storage if it is already destroyed
	case types.StateDestroyed:
		if err := handleServiceStateDestroyed(ss, s); err != nil {
			log.Debugf("%s:observe:serviceStateDestroyed err:> %s", logPrefix, err.Error())
			return err
		}
		break
	}

	if ss.service == nil {
		return nil
	}

	// update service state
	status := ss.service.Status

	ss.service = s
	serviceStatusState(ss)

	if status.State != ss.service.Status.State || status.Message != ss.service.Status.Message {
		log.Debug(status.State, status.Message, ss.service.Status.State, ss.service.Status.Message)
		sm := distribution.NewServiceModel(context.Background(), envs.Get().GetStorage())
		sm.Set(ss.service)
	}

	return nil
}

// handleServiceStateCreated handles service created state
func handleServiceStateCreated(ss *ServiceState, svc *types.Service) error {

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

	// Endpoint provision call
	if err := serviceEndpointProvision(ss, svc); err != nil {
		log.Errorf("%s:> endpoint provision err:", logServicePrefix, err.Error())
		return err
	}

	// Deployment provision call
	if err := serviceDeploymentProvision(ss, svc); err != nil {
		log.Errorf("%s:> deployment provision err:", logServicePrefix, err.Error())
		return err
	}

	return nil
}

// handleServiceStateReady handles service ready state
func handleServiceStateReady(ss *ServiceState, svc *types.Service) error {
	return nil
}

// handleServiceStateError handles service error state
func handleServiceStateError(ss *ServiceState, svc *types.Service) error {
	return nil
}

// handleServiceStateDegradation handles service degradation state
func handleServiceStateDegradation(ss *ServiceState, svc *types.Service) error {
	return nil
}

// handleServiceStateDestroy handles service destroy state
func handleServiceStateDestroy(ss *ServiceState, svc *types.Service) error {

	if ss.endpoint != nil  {
		if err := EndpointRemove(ss.endpoint); err != nil {
			log.Errorf("%s:> endpoint remove err:", logServicePrefix, err.Error())
			return err
		}
		ss.endpoint = nil
	}

	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
	if len(ss.deployment.list) == 0  {
		sm := distribution.NewServiceModel(context.Background(), envs.Get().GetStorage())
		if err := sm.Remove(svc); err != nil {
			log.Errorf("%s:> service remove err:", logServicePrefix, err.Error())
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


	return nil
}

// handleServiceStateDestroyed handles service destroyed state
func handleServiceStateDestroyed(ss *ServiceState, svc *types.Service) error {

	if ss.endpoint != nil {
		if err := EndpointRemove(ss.endpoint); err != nil {
			log.Errorf("%s:> endpoint remove err:", logServicePrefix, err.Error())
			return err
		}
		ss.endpoint = nil
	}

	if len(ss.deployment.list) > 0 {
		dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
		for _, d := range ss.deployment.list {

			if d.Status.State == types.StateDestroyed {
				if err := dm.Remove(d); err != nil {
					return err
				}
			}

			if d.Status.State != types.StateDestroy {
				if err := dm.Destroy(d); err != nil {
					return err
				}
			}
		}

		return nil
	}

	sm := distribution.NewServiceModel(context.Background(), envs.Get().GetStorage())
	if err := sm.Remove(svc); err != nil {
		log.Errorf("%s:> service remove err:", logServicePrefix, err.Error())
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

	ef := len(svc.Spec.Network.Ports) > 0
	switch true {
	// remove endpoint if not needed but exists
	case !ef && ss.endpoint != nil:
		if err := EndpointRemove(ss.endpoint); err != nil {
			log.Errorf("%s:> endpoint remove err:", logServicePrefix, err.Error())
			return err
		}
		ss.endpoint = nil
		break

		// create endpoint if needed and not exists
	case ef && ss.endpoint == nil:
		e, err := EndpointCreate(svc.Meta.Namespace, svc.Meta.Name, svc.Meta.Endpoint, svc.Spec.Network)
		if err != nil {
			log.Errorf("%s:> endpoint create err:", logServicePrefix, err.Error())
			return err
		}
		ss.endpoint = e
		break

		// update endpoint if spec different
	case ef && ss.endpoint != nil:
		if EndpointValidate(ss.endpoint, svc.Spec.Network) {
			break
		}

		e, err := EndpointUpdate(ss.endpoint, svc.Spec.Network)
		if err != nil {
			log.Errorf("%s:> endpoint update err:", logServicePrefix, err.Error())
			return err
		}
		ss.endpoint = e

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
	if d != nil && d.Spec.Replicas != svc.Spec.Replicas {
		if err := deploymentScale(d, svc.Spec.Replicas); err != nil {
			log.Errorf("%s:> deployment scale err:", logServicePrefix, err.Error())
			return err
		}
		return nil
	}

	// create deployment if needed
	if d == nil {

		d, err := deploymentCreate(svc)
		if err != nil {
			log.Errorf("%s:> deployment create err:", logServicePrefix, err.Error())
			return err
		}

		ss.deployment.list[d.SelfLink()] = d

		if ss.deployment.provision != nil {
			if err := deploymentDestroy(ss, ss.deployment.provision); err != nil {
				log.Errorf("%s:> deployment cancel err:", logServicePrefix, err.Error())
				return err
			}
		}

		ss.deployment.provision = d
	}



	return nil
}

// serviceStatusState calculates current service status based on deployments
func serviceStatusState(ss *ServiceState)  {

	if ss.service.Status.State == types.StateDestroy {

		if len(ss.deployment.list) == 0 {
			ss.service.Status.State = types.StateDestroyed
		}

		return
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

	return
}
