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
	_c "github.com/lastbackend/lastbackend/pkg/common/context"
	"github.com/lastbackend/lastbackend/pkg/controller/config"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

var _ctx Context

type Context struct {
	_c.IContext

	logger  logger.ILogger
	storage *storage.Storage
	config  *config.Config
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

func (c *Context) SetStorage(storage *storage.Storage) {
	c.storage = storage
}

func (c *Context) GetStorage() *storage.Storage {
	return c.storage
}

func (c *Context) Background() context.Context {
	return context.Background()
}
