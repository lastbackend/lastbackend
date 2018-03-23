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
		// Mark pod for destroy
		for _, p := range pl {
			if err := pm.Destroy(context.Background(), p); err != nil {
				log.Errorf("controller:deployment:controller:provision: destroy deployment err: %s", err.Error())
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

			if p.Status.Stage == types.StageError {
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

			if p.Status.Stage == types.StageProvision {
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

			if d.Status.Stage == types.StageReady {
				if err := pm.Destroy(context.Background(), p); err != nil {
					log.Errorf("controller:service:controller:provision: remove pod err: %s", err.Error())
					continue
				}
				count--
			}

		}
	}

	// Update deployment state
	d.Status.Stage = types.StageProvision
	if err := distribution.NewDeploymentModel(context.Background(), stg).SetStatus(d); err != nil {
		log.Errorf("controller:deployment:controller:provision: deployment set state err: %s", err.Error())
		return err
	}

	return nil
}

// Handler Deployment status
func HandleStatus(d *types.Deployment) error {
	return nil
}