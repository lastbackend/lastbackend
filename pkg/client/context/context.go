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
	"github.com/lastbackend/lastbackend/pkg/client/config"
	"github.com/lastbackend/lastbackend/pkg/client/storage"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/util/http"
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
	logger  *logger.Logger
	http    *http.RawReq
	storage *storage.DB
	config  *config.Config
	mock    bool
}

func (c *ctx) SetLogger(log *logger.Logger) {
	c.logger = log
}

func (c *ctx) GetLogger() *logger.Logger {
	return c.logger
}

func (c *ctx) SetHttpClient(http *http.RawReq) {
	c.http = http
}

func (c *ctx) GetHttpClient() *http.RawReq {
	return c.http
}

func (c *ctx) SetStorage(storage *storage.DB) {
	c.storage = storage
}

func (c *ctx) GetStorage() *storage.DB {
	return c.storage
}

func (c *ctx) SetConfig(cfg *config.Config) {
	c.config = cfg
}

func (c *ctx) GetConfig() *config.Config {
	return c.config
}
