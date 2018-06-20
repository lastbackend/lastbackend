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

package observe

import (
	"golang.org/x/net/context"

	"github.com/lastbackend/lastbackend/pkg/controller/runtime/pod"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/deployment"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/service"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/cache"
)

type Controller interface {
	Observe(ctx context.Context, cache *cache.Cache)
	Pause()
	Resume()
}

type Controllers struct {
	pc Controller
	dc Controller
	sc Controller
}

type Observer struct {
	cache *cache.Cache
	ctrl  Controllers
}

func New() *Observer {
	o := new(Observer)

	o.cache = new(cache.Cache)

	o.ctrl.pc = pod.NewPodController(context.Background())
	o.ctrl.dc = deployment.NewDeploymentController(context.Background())
	o.ctrl.sc = service.NewServiceController(context.Background())

	return o
}

func (o Observer) Run() {
	go o.ctrl.pc.Observe(context.Background(), o.cache)
	go o.ctrl.dc.Observe(context.Background(), o.cache)
	go o.ctrl.sc.Observe(context.Background(), o.cache)
}
