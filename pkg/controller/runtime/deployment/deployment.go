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

package deployment

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

// Provision deployment
// Remove deployment or cancel if deployment is market for destroy
// Remove deployment if no active pods present and deployment is marked for destroy
func Provision(d *types.Deployment) error {

	var (
		stg      = envs.Get().GetStorage()
		replicas int
	)

	log.Debugf("controller:deployment:controller:provision: provision deployment: %s", d.SelfLink())

	dm := distribution.NewDeploymentModel(context.Background(), stg)
	if d, err := dm.Get(d.Meta.Namespace, d.Meta.Service, d.Meta.Name); d == nil || err != nil {
		if d == nil {
			return errors.New(store.ErrEntityNotFound)
		}
		log.Errorf("controller:deployment:controller:provision: get deployment error: %s", err.Error())
		return err
	}

	// Get all pods by service
	pm := distribution.NewPodModel(context.Background(), stg)
	pl, err := pm.ListByDeployment(d.Meta.Namespace, d.Meta.Service, d.Meta.Name)
	if err != nil {
		log.Errorf("controller:deployment:controller:provision: get pod list error: %s", err.Error())
		return err
	}

	// Check deployment is marked for destroy
	if d.Spec.State.Destroy {

		if len(pl) == 0 {
			log.Debugf("controller:deployment:controller:provision: remove deployment %s", d.SelfLink())
			if err := dm.Remove(d); err != nil {
				log.Errorf("controller:deployment:controller:provision: remove pod err: %s", err.Error())
				return nil
			}
		}

		for _, p := range pl {

			if p.Meta.Node == "" {
				// Mark pod for destroy
				log.Debugf("controller:deployment:controller:provision: remove pod %s", p.SelfLink())
				if err := pm.Remove(context.Background(), p); err != nil {
					log.Errorf("controller:deployment:controller:provision: remove pod err: %s", err.Error())
					continue
				}
			}

			// Mark pod for destroy
			log.Debugf("controller:deployment:controller:provision: mark pod for destroy %s", p.SelfLink())
			if err := pm.Destroy(context.Background(), p); err != nil {
				log.Errorf("controller:deployment:controller:provision: destroy pod err: %s", err.Error())
			}
		}
		return nil
	}

	// Replicas management
	for _, p := range pl {
		if !p.Spec.State.Destroy {
			replicas++
		}
	}

	// Create new replicas
	if replicas < d.Spec.Replicas {
		log.Debug("controller:deployment:controller:provision: create new pods")
		for i := 0; i < (d.Spec.Replicas - replicas); i++ {
			if _, err := pm.Create(d); err != nil {
				log.Errorf("controller:deployment:controller:provision: create new pod err: %s", err.Error())
			}
		}
	}

	// Remove unneeded replicas
	if replicas > d.Spec.Replicas {

		count := replicas - d.Spec.Replicas
		log.Debug("controller:deployment:controller:provision: remove unneeded pods")

		// Remove pods in error state
		for _, p := range pl {

			// check replicas needs to be destroyed
			if count <= 0 {
				break
			}

			if p.Status.Stage == types.StateError {
				if err := pm.Destroy(context.Background(), p); err != nil {
					log.Errorf("controller:service:controller:provision: remove pod err: %s", err.Error())
					continue
				}
				count--
			}

		}

		// Remove pods in provision state
		for _, p := range pl {

			// check replicas needs to be destroyed
			if count <= 0 {
				break
			}

			if p.Status.Stage == types.StateProvision {
				if err := pm.Destroy(context.Background(), p); err != nil {
					log.Errorf("controller:service:controller:provision: remove pod err: %s", err.Error())
					continue
				}
				count--
			}

		}

		// Remove ready pods
		for _, p := range pl {

			// check replicas needs to be destroyed
			if count <= 0 {
				break
			}

			if d.Status.State == types.StateReady {
				if err := pm.Destroy(context.Background(), p); err != nil {
					log.Errorf("controller:service:controller:provision: remove pod err: %s", err.Error())
					continue
				}
				count--
			}

		}
	}

	// Update deployment state
	d.Status.State = types.StateProvision
	if err := distribution.NewDeploymentModel(context.Background(), stg).SetStatus(d); err != nil {
		log.Errorf("controller:deployment:controller:provision: deployment set state err: %s", err.Error())
		return err
	}

	return nil
}

// Handler Deployment status
func HandleStatus(d *types.Deployment) error {

	var (
		stg    = envs.Get().GetStorage()
		msg    = "controller:deployment:controller:status:"
		status = make(map[string]int)
	)

	dm := distribution.NewDeploymentModel(context.Background(), stg)
	sm := distribution.NewServiceModel(context.Background(), stg)

	// Skip state handle
	if d.Status.State == types.StateDestroy {
		log.Debugf("%s> skip deployment status [%s] handle: %s", msg, d.Status.State, d.Meta.Name)
		return nil
	}

	svc, err := sm.Get(d.Meta.Namespace, d.Meta.Service)
	if err != nil {
		log.Errorf("%s> get service err: %s", msg, err.Error())
		return err
	}

	if svc == nil {
		log.Warnf("%s> service [%s:%s] not found", msg, d.Meta.Namespace, d.Meta.Service)
		if err := dm.Remove(d); err != nil {
			log.Errorf("%s> remove deployment err: %s", msg, err.Error())
			return err
		}

		return nil
	}

	dl, err := dm.ListByService(svc.Meta.Namespace, svc.Meta.Name)
	if err != nil {
		log.Errorf("%s> get deployment list err: %s", msg, err.Error())
		return err
	}

	for _, di := range dl {
		switch di.Status.State {
		case types.StateProvision:
			status[types.StateProvision] += 1
			break
		case types.StateRunning:
			status[types.StateRunning] += 1
			break
		case types.StateStopped:
			status[types.StateStopped] += 1
			break
		case types.StateDestroyed:
			status[types.StateDestroyed] += 1
			break
		}
	}

	var curr = svc.Spec.Meta.Name == d.Spec.Meta.Name

	switch true {
	case status[types.StateProvision] > 0:
		svc.Status.State = types.StateProvision
		svc.Status.Message = types.EmptyString
		break
	case status[types.StateDestroyed] == 1 && d.Status.State == types.StateDestroyed:
		svc.Status.State = types.StateDestroyed
		svc.Status.Message = types.EmptyString
		break
	case curr && d.Status.State == types.StateError:
		if (status[types.StateRunning] + status[types.StateStopped]) == 0 {
			svc.Status.State = types.StateError
			svc.Status.Message = d.Status.Message
			break
		}
		svc.Status.State = types.StateWarning
		svc.Status.Message = d.Status.Message
		break
	case curr && d.Status.State == types.StateRunning:
		svc.Status.State = types.StateRunning
		svc.Status.Message = types.EmptyString
		break
	case curr && d.Status.State == types.StateStopped:
		svc.Status.State = types.StateStopped
		svc.Status.Message = types.EmptyString
		break
	case curr && d.Status.State == types.StateWarning:
		svc.Status.State = types.StateWarning
		svc.Status.Message = d.Status.Message
		break
	}

	if (status[types.StateRunning] + status[types.StateStopped]) > 1 {
		for _, di := range dl {
			if svc.Spec.Meta.Name == di.Spec.Meta.Name {
				continue
			}

			if di.Status.State == types.StateRunning || di.Status.State == types.StateStopped {
				if err := dm.Destroy(di); err != nil {
					log.Errorf("%s> destroy deployment err: %s", msg, err.Error())
				}
			}

		}
	}

	// Remove destroyed deployment
	if d.Status.State == types.StateDestroyed {
		if err := dm.Remove(d); err != nil {
			log.Errorf("%s> remove deployment err: %s", msg, err.Error())
			return err
		}
	}

	if err := sm.SetStatus(svc); err != nil {
		log.Errorf("%s> set deployment status err: %s", msg, err.Error())
		return err
	}

	return nil
}
