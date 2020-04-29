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

package endpoint

import (
	"context"

	"github.com/lastbackend/lastbackend/internal/discovery/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logLevel  = 3
	logPrefix = "runtime:endpoint"
)

func Watch(ctx context.Context) {

	log.Debugf("%s:restore:> watch change endpoint start", logPrefix)

	var (
		em    = service.NewEndpointModel(ctx, envs.Get().GetStorage())
		cache = envs.Get().GetCache().Endpoint()
		event = make(chan models.EndpointEvent)
	)

	go func() {
		for {
			select {
			case e := <-event:
				{

					if e.Data == nil {
						continue
					}

					endpoint := e.Data

					switch e.Action {
					case models.EventActionCreate:
						fallthrough
					case models.EventActionUpdate:
						cache.Del(endpoint.Spec.Domain)
						envs.Get().GetCache().Endpoint().Set(endpoint.Spec.Domain, []string{endpoint.Spec.IP})
						continue
					case models.EventActionDelete:
						cache.Del(endpoint.Spec.Domain)
						continue
					}

				}
			}
		}
	}()

	go em.Watch(event, nil)
}
