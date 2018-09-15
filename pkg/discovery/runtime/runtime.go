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
	"github.com/lastbackend/lastbackend/pkg/discovery/envs"
	"github.com/lastbackend/lastbackend/pkg/discovery/runtime/endpoint"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/system"
)

const (
	logPrefix = "discovery:runtime"
	logLevel  = 3
)

type Runtime struct {
	ctx context.Context
}

func (r *Runtime) Restore() {
	log.V(logLevel).Debugf("%s:restore:> restore init", logPrefix)
}

func (r *Runtime) Loop() error {

	log.V(logLevel).Debugf("%s:loop:> update current discovery service info", logPrefix)

	hostname, err := system.GetHostname()
	if err != nil {
		log.Errorf(" can not get discovery hostname:%s", err.Error())
		return err
	}

	ip, err := system.GetNodeIP()
	if err != nil {
		log.Errorf(" can not get discovery ip:%s", err.Error())
		return err
	}

	discovery := new(types.Discovery)
	discovery.Meta.Name = hostname
	discovery.Status.IP = ip

	dm := distribution.NewDiscoveryModel(context.Background(), envs.Get().GetStorage())
	if err := dm.Set(discovery); err != nil {
		log.Errorf(" can not get discovery data from storage:%s", err.Error())
		return err
	}

	log.V(logLevel).Debugf("%s:loop:> watch endpoint start", logPrefix)
	endpoint.Watch(r.ctx)

	return nil
}

func NewRuntime(ctx context.Context) *Runtime {
	return &Runtime{ctx: ctx}
}
