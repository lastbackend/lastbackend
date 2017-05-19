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

package pod

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/discovery/context"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"strings"
)

type PodController struct {
	context *context.Context
	pods    chan *types.Pod

	active bool
}

func (pc *PodController) Watch(services chan *types.Service) {
	var (
		log = pc.context.GetLogger()
		stg = pc.context.GetStorage()
	)

	log.Debug("PodController: start watch")
	go func() {
		for {
			select {
			case p := <-pc.pods:
				{

					if !pc.active {
						log.Debug("PodController: skip management cause it is in slave mode")
						continue
					}

					if p == nil || p.Meta.Hostname == "" {
						continue
					}

					endpoint := fmt.Sprintf("%s.%s", p.Meta.Name, *context.Get().GetConfig().SystemDomain)
					endpoint = strings.Replace(endpoint, ":", "-", -1)

					srv, err := stg.Service().GetByPodName(context.Get().Background(), p.Meta.Name)
					if err != nil {
						if err.Error() == store.ErrKeyNotFound {
							if err := stg.Endpoint().Remove(context.Get().Background(), endpoint); err != nil {
								log.Errorf("Endpoint: remove endpoint error %s", err.Error())
							}
						} else {
							log.Errorf("Endpoint: get service error %s", err.Error())
						}
						continue
					}

					node, err := stg.Node().Get(context.Get().Background(), p.Meta.Hostname)
					if err != nil {
						log.Errorf("Endpoint: get node error %s", err.Error())
						break
					}

					serviceEndpoint := fmt.Sprintf("%s-%s.%s", srv.Meta.Name, srv.Meta.Namespace, *context.Get().GetConfig().SystemDomain)
					serviceEndpoint = strings.Replace(serviceEndpoint, ":", "-", -1)

					if p.Spec.State == types.StateDestroyed {
						if err := stg.Endpoint().Remove(context.Get().Background(), endpoint); err != nil {
							log.Errorf("Endpoint: remove endpoint error %s", err.Error())
						}
						continue
					}

					if err := stg.Endpoint().Upsert(context.Get().Background(), endpoint, []string{node.Meta.IP}); err != nil {
						log.Errorf("Endpoint: upsert endpoint error %s", err.Error())
						continue
					}

					services <- srv
				}
			}
		}
	}()

	stg.Pod().Watch(pc.context.Background(), pc.pods)
}

func (pc *PodController) Pause() {
	pc.active = false
}

func (pc *PodController) Resume() {
	pc.active = true
}

func NewPodController(ctx *context.Context) *PodController {
	sc := new(PodController)
	sc.context = ctx
	sc.active = false
	sc.pods = make(chan *types.Pod)
	return sc
}
