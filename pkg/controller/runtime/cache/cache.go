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

package cache

const (
	EventCreate = "create"
	EventUpdate = "update"
	EventRemove = "remove"
)

type Cache struct {
	IsReady     bool
	Pods        *PodCache
	Deployments *DeploymentCache
	Services    *ServiceCache

	ready chan bool
}

func New() *Cache {
	c := new(Cache)
	c.ready = make(chan bool)
	c.Pods = NewPodCache()
	c.Deployments = NewDeploymentCache()
	c.Services = NewServiceCache()

	go func() {
		for {
			select {
			case <-c.Pods.Ready():
				c.checkReady()
			case <-c.Deployments.Ready():
				c.checkReady()
			case <-c.Services.Ready():
				c.checkReady()
			}
		}
	}()

	return c
}

func (c Cache) Ready() <-chan bool {
	return c.ready
}

func (c *Cache) checkReady() {
	if c.Services.IsReady && c.Deployments.IsReady && c.Pods.IsReady {
		c.IsReady = true
		c.ready <- true
	}
}
