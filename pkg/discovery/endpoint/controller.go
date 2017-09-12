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

package endpoint

import (
	"github.com/lastbackend/lastbackend/pkg/cache"
	"github.com/lastbackend/lastbackend/pkg/discovery/context"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util"
	"github.com/lastbackend/lastbackend/pkg/log"
)

type EndpointController struct {
	context   *context.Context
	endpoints chan string
	cache     *cache.EndpointCache

	active bool
}

func (ec *EndpointController) Watch() {
	var (
		stg = ec.context.GetStorage()
	)

	log.V(logLevel).Debug("EndpointController: start watch")

	go func() {
		for {
			select {
			case name := <-ec.endpoints:
				{

					if !ec.active {
						log.V(logLevel).Debug("EndpointController: skip management cause it is in slave mode")
						continue
					}

					i, err := stg.Endpoint().Get(context.Get().Background(), name)
					if err != nil {
						if err.Error() == store.ErrKeyNotFound {
							if err = ec.cache.Del(name); err != nil {
								log.V(logLevel).Debugf("EndpointController: remove endpoint `%s` ips from cache", name)
							}
						} else {
							log.V(logLevel).Errorf("EndpointController: get endpoint `%s` ips for domain err: %s", name, err.Error())
						}
						continue
					}

					ips, err := util.ConvertStringIPToNetIP(i)
					if err != nil {
						log.V(logLevel).Errorf("EndpointController: convert endpoint `%s` ips to net ips err: %s", name, err.Error())
						continue
					}

					if err = ec.cache.Set(name, ips); err != nil {
						log.V(logLevel).Errorf("EndpointController: save endpoint `%s` ips to cache err: %s", name, err.Error())
						continue
					}
				}
			}
		}
	}()

	stg.Endpoint().Watch(ec.context.Background(), ec.endpoints)
}

func (ec *EndpointController) Pause() {
	log.Debugf("EndpointController: pause")
	ec.active = false
}

func (ec *EndpointController) Resume() {
	log.Debugf("EndpointController: resume")
	ec.active = true
}

func NewEndpointController(ctx *context.Context) *EndpointController {
	sc := new(EndpointController)
	sc.context = ctx
	sc.active = false
	sc.endpoints = make(chan string)
	sc.cache = ctx.GetCache().Endpoints()
	return sc
}
