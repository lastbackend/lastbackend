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

package docker

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"

	d "github.com/docker/docker/api/types"
)

func (r *Runtime) Subscribe(ctx context.Context, container chan *types.Container) error {

	log.Debugf("%s:subscribe:> create new event listener subscribe", logPrefix)

	if _, err := r.client.Ping(ctx); err != nil {
		log.Errorf("%s:subscribe:> can not ping docker client err: %v", logPrefix, err)
		return err
	}

	event, err := r.client.Events(ctx, d.EventsOptions{})

	for {
		select {
		case e := <-event:
			log.Debugf("%s:subscribe:> event type: %s action: %v", logPrefix, e.Type, e.Action)

			if len(e.ID) == 0 {
				continue
			}

			if e.Status == types.StateDestroy {
				c := new(types.Container)
				c.ID = e.ID
				c.State = types.StateDestroyed
				container <- c
				continue
			}

			c, err := r.Inspect(ctx, e.ID)
			if err != nil {
				log.Errorf("%s:subscribe:> container inspect err: %v", logPrefix, err)
				continue
			}

			if c == nil {
				continue
			}

			container <- c

		case err := <-err:
			if err == context.Canceled {
				log.Warnf("%s:subscribe:> context canceled err: %v", logPrefix, err)
				return nil
			}
			log.Errorf("%s:subscribe:> event listening err: %v", logPrefix, err)
			return err
		}
	}
}
