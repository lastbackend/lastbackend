//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

// Provision service
// Remove deployment or cancel if service is market for destroy
// Remove service if no active deployments present and service is marked for destroy
func Provision(svc *types.Service) error {

	var (
		stg = envs.Get().GetStorage()
	)

	log.Debugf("Service Controller: provision service: %s/%s", svc.Meta.Namespace, svc.Meta.Name)

	// Check service is marked for destroy

	// Get all deployments per service

	dm := distribution.NewDeploymentModel(context.Background(), stg)
	dl, err := dm.ListByService(svc.Meta.Namespace, svc.Meta.Name)
	if err != nil {
		log.Errorf("Controller: service controller: get deployments list error: %s", err.Error())
		return err
	}

	for _, d := range dl {
		d.Spec.Template.Termination = 1
	}

	return nil
}
