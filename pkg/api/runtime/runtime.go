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

import "github.com/lastbackend/lastbackend/pkg/api/envs"

type Runtime struct {
}

func New() *Runtime {
	return new(Runtime)
}

func (r *Runtime) Run() {

	var (
		stg = envs.Get().GetStorage()
		c   = envs.Get().GetCache()
	)

	go c.Node().CachePods(stg.Node().EventPodSpec)
	go c.Node().CacheVolumes(stg.Node().EventVolumeSpec)
	go c.Node().CacheEndpoints(stg.Endpoint().EventSpec)
	go c.Node().Del(stg.Node().EventStatus)

	go c.Ingress().CacheRoutes(stg.Route().WatchSpecEvents)
	go c.Ingress().Status(stg.Ingress().WatchStatus)
}
