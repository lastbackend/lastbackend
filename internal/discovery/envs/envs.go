//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/internal/discovery/cache"
	"github.com/lastbackend/lastbackend/internal/discovery/state"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/client/types"
)

var _env Env

type Env struct {
	client  types.DiscoveryClientV1
	storage storage.Storage
	cache   *cache.Cache
	state   *state.State
}

func Get() *Env {
	return &_env
}

func (c *Env) SetStorage(storage storage.Storage) {
	c.storage = storage
}

func (c *Env) GetStorage() storage.Storage {
	return c.storage
}

func (c *Env) SetCache(cache *cache.Cache) {
	c.cache = cache
}

func (c *Env) GetCache() *cache.Cache {
	return c.cache
}

func (c *Env) SetState(st *state.State) {
	c.state = st
}

func (c *Env) GetState() *state.State {
	return c.state
}

func (c *Env) SetClient(client types.DiscoveryClientV1) {
	c.client = client
}

func (c *Env) GetClient() types.DiscoveryClientV1 {
	return c.client
}
