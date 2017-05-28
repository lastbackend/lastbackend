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

package context

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/api/config"
	_c "github.com/lastbackend/lastbackend/pkg/common/context"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/sockets"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/http"
)

var _ctx Context

type Context struct {
	_c.IContext

	logger               logger.ILogger
	storage              storage.IStorage
	config               *config.Config
	httpTemplateRegistry *http.RawReq
	wssHub               *sockets.Hub
}

func Get() *Context {
	return &_ctx
}

func (c *Context) SetLogger(log logger.ILogger) {
	c.logger = log
}

func (c *Context) GetLogger() logger.ILogger {
	return c.logger
}

func (c *Context) SetConfig(cfg *config.Config) {
	c.config = cfg
}

func (c *Context) GetConfig() *config.Config {
	return c.config
}

func (c *Context) SetStorage(storage storage.IStorage) {
	c.storage = storage
}

func (c *Context) GetStorage() storage.IStorage {
	return c.storage
}

func (c *Context) SetHttpTemplateRegistry(http *http.RawReq) {
	c.httpTemplateRegistry = http
}

func (c *Context) GetHttpTemplateRegistry() *http.RawReq {
	return c.httpTemplateRegistry
}

func (c *Context) SetWssHub(hub *sockets.Hub) {
	c.wssHub = hub
}

func (c *Context) GetWssHub() *sockets.Hub {
	return c.wssHub
}

func (c *Context) Background() context.Context {
	return context.Background()
}
