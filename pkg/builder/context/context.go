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
	"github.com/lastbackend/lastbackend/pkg/builder/config"
	_c "github.com/lastbackend/lastbackend/pkg/common/context"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

var _ctx ctx

type Context struct {
	context.Context
}

type ctx struct {
	_c.IContext

	logger  logger.ILogger
	storage storage.IStorage
	config  *config.Config
}

func Get() *ctx {
	return &_ctx
}

func (c *ctx) SetLogger(log logger.ILogger) {
	c.logger = log
}

func (c *ctx) GetLogger() logger.ILogger {
	return c.logger
}

func (c *ctx) SetConfig(cfg *config.Config) {
	c.config = cfg
}

func (c *ctx) GetConfig() *config.Config {
	return c.config
}

func (c *ctx) SetStorage(storage storage.IStorage) {
	c.storage = storage
}

func (c *ctx) GetStorage() storage.IStorage {
	return c.storage
}

func (c *ctx) Background() context.Context {
	return context.Background()
}
