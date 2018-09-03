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

package runtime

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
)

func RouteManage(ctx context.Context, name string, route *types.RouteManifest) error {

	defer envs.Get().GetIngress().Update(ctx)

	log.Debugf("route manage: %s", name)


	log.Debugf("total routes: %d", len(envs.Get().GetState().Routes().GetRoutes()))

	if route.State == types.StateDestroyed {
		envs.Get().GetState().Routes().DelRoute(name)
		return nil
	}

	envs.Get().GetState().Routes().SetRoute(name, route)
	return nil
}
