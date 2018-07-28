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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

func ServiceRemove(svc *types.Service) error {
	sm := distribution.NewServiceModel(context.Background(), envs.Get().GetStorage())
	return sm.Remove(svc)
}

func ServiceProvision(svc *types.Service) (*types.Deployment, error) {
	return DeploymentCreate(svc)
}

func ServiceDestroy(svc *types.Service, dl map[string]*types.Deployment) error {

	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
	for _, d := range dl {

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

func ServiceSync(svc *types.Service, d *types.Deployment) error {

	sm := distribution.NewServiceModel(context.Background(), envs.Get().GetStorage())

	if d == nil {
		svc.Status.State = types.StateWarning
		svc.Status.Message = "unknown state: no active or provision deployment"
		return sm.Set(svc)
	}

	switch d.Status.State {
	case types.StateReady:
	case types.StateWarning:
	case types.StateError:
		svc.Status.State = d.Status.State
		svc.Status.Message = d.Status.Message
	}

	return sm.Set(svc)
}
