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

package endpoint

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
)

const logEndpointPrefix = "runtime:endpoint:>"

func Manage(ctx context.Context, key string, spec *types.EndpointSpec) error {
	log.Debugf("%s manage: %s", logEndpointPrefix, key)

	ep := envs.Get().GetState().Endpoints().GetEndpoint(key)
	if ep != nil {
		if check(spec, ep) {
			Update(ctx, key, ep, spec)
		}
	}

	status, err := Create(ctx, key, spec)
	if err != nil {
		log.Errorf("%s create error: %s", logEndpointPrefix, err.Error())
		return err
	}

	envs.Get().GetState().Endpoints().SetEndpoint(key, status)
	return nil
}

func Restore(ctx context.Context) error {
	log.Debugf("%s restore", logEndpointPrefix)
	cpi := envs.Get().GetCPI()
	endpoints, err := cpi.Info(ctx)
	if err != nil {
		log.Errorf("%s restore error: %s", err.Error())
		return err
	}
	envs.Get().GetState().Endpoints().SetEndpoints(endpoints)
	return nil
}

func Create(ctx context.Context, key string, spec *types.EndpointSpec) (*types.EndpointStatus, error) {
	log.Debugf("%s create %s", logEndpointPrefix)
	cpi := envs.Get().GetCPI()
	return cpi.Create(ctx, spec)
}

func Update(ctx context.Context, endpoint string, status *types.EndpointStatus, spec *types.EndpointSpec) (*types.EndpointStatus, error) {
	log.Debugf("%s update %s", logEndpointPrefix, endpoint)
	cpi := envs.Get().GetCPI()
	return cpi.Update(ctx, status, spec)
}

func Destroy(ctx context.Context, endpoint string, status *types.EndpointStatus) error {
	log.Debugf("%s destroy", logEndpointPrefix, endpoint)
	cpi := envs.Get().GetCPI()
	return cpi.Destroy(ctx, status)
}

func check(spec *types.EndpointSpec, status *types.EndpointStatus) bool {
	return false
}
