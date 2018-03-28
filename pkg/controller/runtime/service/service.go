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
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

// Provision service
// Remove deployment or cancel if service is market for destroy
// Remove service if no active deployments present and service is marked for destroy
func Provision(svc *types.Service) error {

	var (
		stg = envs.Get().GetStorage()
		spc bool
		msg = "controller:service:controller:provision:"
	)

	sm := distribution.NewServiceModel(context.Background(), stg)
	if d, err := sm.Get(svc.Meta.Namespace, svc.Meta.Name); d == nil || err != nil {
		if d == nil {
			return errors.New(store.ErrEntityNotFound)
		}
		log.Errorf("%s> get deployment error: %s", msg, err.Error())
		return err
	}

	log.Debugf("%s> provision service: %s", msg, svc.SelfLink())

	// Get all deployments per service
	dm := distribution.NewDeploymentModel(context.Background(), stg)
	dl, err := dm.ListByService(svc.Meta.Namespace, svc.Meta.Name)
	if err != nil {
		log.Errorf("%s> get deployments list error: %s", msg, err.Error())
		return err
	}

	if len(dl) == 0 && svc.Status.State == types.StateDestroy {
		if err := sm.Remove(svc); err != nil {
			log.Errorf("%s> remove service err: %s", err.Error())
			return nil
		}
	}

	for _, d := range dl {

		if d.Spec.State.Destroy {
			continue
		}

		if !svc.Spec.State.Destroy && d.Spec.Meta.Name == svc.Spec.Meta.Name {
			spc = true
			continue
		}

		if svc.Spec.State.Destroy && !d.Spec.State.Destroy {
			if err := dm.Destroy(d); err != nil {
				log.Errorf("%s> remove deployment err: %s", err.Error())
				continue
			}
		}
	}

	// Check service is marked for destroy
	if spc || svc.Spec.State.Destroy {
		return nil
	}

	// Create new deployment
	if _, err := dm.Create(svc); err != nil {
		log.Errorf("%s> create deployment err: %s", msg, err.Error())
		svc.Status.State = types.StateError
		svc.Status.Message = err.Error()
	}

	// Update service state
	svc.Status.State = types.StateProvision
	if err := distribution.NewServiceModel(context.Background(), stg).SetStatus(svc); err != nil {
		log.Errorf("%s> service set state err: %s", msg, err.Error())
		return err
	}

	return nil
}

// HandleStatus handles status of service
func HandleStatus(svc *types.Service) error {

	var (
		stg     = envs.Get().GetStorage()
		msg     = "controller:deployment:service:status:"
		status  = make(map[string]int)
	)

	if svc == nil {
		log.Errorf("%s> service is nil", msg)
		return nil
	}

	log.Debugf("%s> handle service [%s] status", msg, svc.SelfLink())

	dm := distribution.NewDeploymentModel(context.Background(), stg)
	sm := distribution.NewServiceModel(context.Background(), stg)

	// Skip state handle
	if svc.Status.State == types.StateDestroy {
		log.Debugf("%s> skip service status [%s] handle: %s", msg, svc.SelfLink())
		return nil
	}

	dl, err := dm.ListByService(svc.Meta.Namespace, svc.Meta.Name)
	if err != nil {
		log.Errorf("%s> get pod list err: %s", msg, err.Error())
		return err
	}

	for _, di := range dl {
		switch di.Status.State {
		case types.StateDestroyed :
			status[types.StateDestroyed]+=1
			break
		}
	}

	if svc.Spec.State.Destroy && len(dl) == status[types.StateDestroyed] {
		log.Debugf("%s:> remove destroyed service: %s", msg, svc.SelfLink())
		if err := sm.Remove(svc); err != nil {
			log.Errorf("%s> remove destroyed service [%s] err: %s", msg, err.Error())
			return err
		}
		return nil
	}

	return nil
}
