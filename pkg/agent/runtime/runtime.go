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

package runtime

import (
	"github.com/docker/docker/client"
	"github.com/lastbackend/lastbackend/pkg/agent/config"
	"github.com/lastbackend/lastbackend/pkg/agent/manager"
)

type Runtime struct {
	Client *client.Client

	pManager *manager.PodManager
	iManager *manager.ImageManager
	eManager *manager.EventManager
}

func New(cfg *config.Runtime) *Runtime {
	var runtime = new(Runtime)
	return runtime.Init(cfg)
}

func (r *Runtime) Init(cfg *config.Runtime) *Runtime {
	r.Client, _ = client.NewEnvClient()

	r.pManager = manager.NewPodManager()
	r.iManager = manager.NewImageManager()
	r.eManager = manager.NewEventManager()

	return r
}

func (r *Runtime) Pods() *manager.PodManager {
	return r.pManager
}

func (r *Runtime) Images() *manager.ImageManager {
	return r.iManager
}

func (r *Runtime) Loop() {

}

func (r *Runtime) Sync() {

}
