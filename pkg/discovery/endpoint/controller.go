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
)

type EndpointController struct {
	context   *context.Context
	endpoints chan string
	cache     *cache.EndpointCache

	active bool
}

func (ec *EndpointController) Watch() {
	var (
		log = ec.context.GetLogger()
		stg = ec.context.GetStorage()
	)

	log.Debug("EndpointController: start watch")
	go func() {
		for {
			select {
			case e := <-ec.endpoints:
				{

					if !ec.active {
						log.Debug("EndpointController: skip management cause it is in slave mode")
						continue
					}

					i, err := stg.Endpoint().Get(context.Get().Background(), e)
					if err != nil {
						if err.Error() != store.ErrKeyNotFound {
							if err = ec.cache.Del(e); err != nil {
								log.Errorf("Endpoint: remove ips from cache error %s", err.Error())
							}
						} else {
							log.Errorf("Endpoint: get ips for domain error %s", err.Error())
						}
						continue
					}

					ips, err := util.ConvertStringIPToNetIP(i)
					if err != nil {
						log.Errorf("Endpoint: convert ips to net ips error %s", err.Error())
						continue
					}

					if err = ec.cache.Set(e, ips); err != nil {
						log.Errorf("Endpoint: save ips to cache error %s", err.Error())
						continue
					}
				}
			}
		}
	}()

	stg.Endpoint().Watch(ec.context.Background(), ec.endpoints)
}

func (ec *EndpointController) Pause() {
	ec.active = false
}

func (ec *EndpointController) Resume() {
	ec.active = true
}

func NewEndpointController(ctx *context.Context) *EndpointController {
	sc := new(EndpointController)
	sc.context = ctx
	sc.active = false
	sc.endpoints = make(chan string)
	sc.cache = ctx.GetCache().EndpointCache
	return sc
}
