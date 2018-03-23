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
	)

	sm := distribution.NewServiceModel(context.Background(), stg)
	if d, err := sm.Get(svc.Meta.Namespace, svc.Meta.Name); d == nil || err != nil {
		if d == nil {
			return errors.New(store.ErrEntityNotFound)
		}
		log.Errorf("controller:service:controller:provision: get deployment error: %s", err.Error())
		return err
	}

	log.Debugf("controller:service:controller:provision: provision service: %s/%s", svc.Meta.Namespace, svc.Meta.Name)

	// Get all deployments per service
	dm := distribution.NewDeploymentModel(context.Background(), stg)
	dl, err := dm.ListByService(svc.Meta.Namespace, svc.Meta.Name)
	if err != nil {
		log.Errorf("controller:service:controller:provision: get deployments list error: %s", err.Error())
		return err
	}

	// Destroy parent deployments
	for _, d := range dl {
		d.Spec.State.Destroy = true
		if err := dm.Destroy(d); err != nil {
			log.Errorf("controller:service:controller:provision: destroy deployment err: %s", err.Error())
		}
	}

	// Check service is marked for destroy
	if svc.Spec.State.Destroy {
		return nil
	}

	// Create new deployment
	if _, err := dm.Create(svc); err != nil {
		log.Errorf("controller:service:controller:provision: create deployment err: %s", err.Error())
		svc.Status.Stage = types.StageError
		svc.Status.Message = err.Error()
	}

	// Update service state
	svc.Status.Stage = types.StageProvision
	if err := distribution.NewServiceModel(context.Background(), stg).SetStatus(svc); err != nil {
		log.Errorf("controller:service:controller:provision: service set state err: %s", err.Error())
		return err
	}

	return nil
}

// HandleStatus handles status of service
func HandleStatus(svc *types.Service) error {

	return nil
}