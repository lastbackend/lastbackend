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

package docker

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

func (r *Runtime) Subscribe(ctx context.Context) (chan *models.Image, error) {
	//
	//log.Debug("Create new event listener subscribe")
	//var cs = make(chan *models.Image)
	//
	//go func() {
	//
	//	if _, err := r.client.Ping(ctx); err != nil {
	//		log.Errorf("Can not ping docker client")
	//		return
	//	}
	//
	//	es, errr := r.client.Events(ctx, d.EventsOptions{})
	//	for {
	//		select {
	//		case e := <-es:
	//
	//			if e.Type != events.ImageEventType {
	//				continue
	//			}
	//
	//			log.Debugf("Image %s", e.ID)
	//
	//			if e.Action == models.StateDestroy {
	//				c := new(models.Image)
	//				c.Meta.ID = e.ID
	//				c.Status.State = models.StateDestroyed
	//				cs <- c
	//				break
	//			}
	//
	//			c, err := r.Inspect(ctx, e.ID)
	//			if err != nil {
	//				log.Errorf("Container inspect err: %s", err.Error())
	//				continue
	//			}
	//			if c == nil {
	//				log.Errorf("Container: container not found")
	//				break
	//			}
	//			break
	//
	//		case err := <-errr:
	//			log.Errorf("Event listening error: %s", err)
	//		}
	//	}
	//}()

	return nil, nil

}
