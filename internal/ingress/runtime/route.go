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

package runtime

import (
	"github.com/lastbackend/lastbackend/internal/ingress/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/log"
)

func (r Runtime) RouteManage(name string, route *models.RouteManifest) (err error) {

	log.Debugf("route manage: %s", name)

	var status = new(models.RouteStatus)

	defer func() {
		if err = r.config.Sync(); err != nil {
			status.State = models.StateError
			status.Message = err.Error()
			envs.Get().GetState().Routes().SetRouteStatus(name, status)
			return
		}

		if status.State == models.StateDestroy {
			envs.Get().GetState().Routes().DelRoute(name)
			return
		}

		envs.Get().GetState().Routes().SetRouteStatus(name, status)
	}()

	if route.State == models.StateDestroyed {
		status.State = models.StateDestroyed
		envs.Get().GetState().Routes().DelRoute(name)
		return nil
	}

	if route.State == models.StateDestroy {
		status.State = models.StateDestroyed
		envs.Get().GetState().Routes().DelRouteManifests(name)
		return nil
	}

	envs.Get().GetState().Routes().SetRouteManifest(name, route)
	status.State = models.StateProvision

	return nil
}
