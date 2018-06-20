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

import "github.com/lastbackend/lastbackend/pkg/distribution/types"

type Cache struct {
	Pods        *PodCache
	Deployments *DeploymentCache
	Services    *ServiceCache
}

func (c Cache) init() {
	c.Pods = new(PodCache)
	c.Deployments = new(DeploymentCache)
	c.Services = new(ServiceCache)
}

type PodCache struct {
	pods map[string]types.Pod

	ch chan types.Pod
}

func (pc PodCache) Subscribe() chan types.Pod {
	return pc.ch
}

type DeploymentCache struct {
	deployments map[string]types.Deployment

	ch chan types.Deployment
}

func (dc DeploymentCache) Subscribe() chan types.Deployment {
	return dc.ch
}

type ServiceCache struct {
	services map[string]types.Service

	ch chan types.Service
}

func (sc ServiceCache) Subscribe() chan types.Service {
	return sc.ch
}
