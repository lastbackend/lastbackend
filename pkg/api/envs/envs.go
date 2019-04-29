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
	"github.com/lastbackend/lastbackend/pkg/api/cache"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/monitor"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

var e Env

type Env struct {
	storage storage.Storage
	cache   *cache.Cache
	monitor *monitor.Monitor
	vault   *types.Vault

	clusterName string
	clusterDesc string

	internalDomain string
	externalDomain string

	accessToken string
}

func Get() *Env {
	return &e
}

func (c *Env) SetStorage(storage storage.Storage) {
	c.storage = storage
}

func (c *Env) GetStorage() storage.Storage {
	return c.storage
}

func (c *Env) SetVault(vault *types.Vault) {
	c.vault = vault
}

func (c *Env) GetVault() *types.Vault {
	return c.vault
}

func (c *Env) SetCache(ch *cache.Cache) {
	c.cache = ch
}

func (c *Env) GetCache() *cache.Cache {
	return c.cache
}

func (c *Env) SetMonitor(monitor *monitor.Monitor) {
	c.monitor = monitor
}

func (c *Env) GetMonitor() *monitor.Monitor {
	return c.monitor
}

func (c *Env) SetClusterInfo(name, desc string) {
	c.clusterName = name
	c.clusterDesc = desc
}

func (c *Env) GetClusterInfo() (string, string) {
	return c.clusterName, c.clusterDesc
}

func (c *Env) SetDomain(internal, external string) {
	c.internalDomain = internal
	c.externalDomain = external
}

func (c *Env) GetDomain() (string, string) {
	return c.internalDomain, c.externalDomain
}

func (c *Env) SetAccessToken(token string) {
	c.accessToken = token
}

func (c *Env) GetAccessToken() string {
	return c.accessToken
}
