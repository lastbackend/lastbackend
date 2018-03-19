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

package context

import (
	"github.com/lastbackend/lastbackend/pkg/api/client/http/v1"
	"github.com/lastbackend/lastbackend/pkg/cli/config"
	"github.com/lastbackend/lastbackend/pkg/cli/storage"
)

var _ctx ctx

func Get() *ctx {
	return &_ctx
}

func Mock() *ctx {
	_ctx.mock = true
	return &_ctx
}

type ctx struct {
	storage storage.IStorage
	config  *config.Config
	token   *string
	mock    bool
	client  *v1.Client
}

func (c *ctx) SetStorage(storage storage.IStorage) {
	c.storage = storage
}

func (c *ctx) SetClient(client *v1.Client) {
	c.client = client
}

func (c *ctx) GetClient() *v1.Client {
	return c.client
}

func (c *ctx) GetStorage() storage.IStorage {
	return c.storage
}

func (c *ctx) SetSessionToken(token string) {
	c.token = &token
}

func (c *ctx) GetSessionToken() *string {
	return c.token
}

func (c *ctx) SetConfig(cfg *config.Config) {
	c.config = cfg
}

func (c *ctx) GetConfig() *config.Config {
	return c.config
}
