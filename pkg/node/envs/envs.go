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

package envs

import (
	"github.com/lastbackend/lastbackend/pkg/cache"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cni"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/node/state"
)

var e Env

func Get() *Env {
	return &e
}

type Env struct {
	cri   cri.CRI
	cni   cni.CNI
	cache *cache.Cache
	state *state.State
}

func (c *Env) SetCri(cri cri.CRI) {
	c.cri = cri
}

func (c *Env) GetCri() cri.CRI {
	return c.cri
}

func (c *Env) SetCNI(n cni.CNI) {
	c.cni = n
}

func (c *Env) GetCNI() cni.CNI {
	return c.cni
}

func (c *Env) SetCache(s *cache.Cache) {
	c.cache = s
}

func (c *Env) GetCache() *cache.Cache {
	return c.cache
}

func (c *Env) SetState(s *state.State) {
	c.state = s
}

func (c *Env) GetState() *state.State {
	return c.state
}
