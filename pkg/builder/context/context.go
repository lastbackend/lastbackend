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
	"github.com/lastbackend/lastbackend/pkg/cri"
	"github.com/lastbackend/lastbackend/pkg/logger"
	_c "github.com/lastbackend/lastbackend/pkg/context"
)

var _ctx ctx

func Get() *ctx {
	return &_ctx
}

type ctx struct {
	_c.Context
	cri     cri.CRI
	logger  *logger.Logger
	config  *config.Config
}

func (c *ctx) SetLogger(log *logger.Logger) {
	c.logger = log
}

func (c *ctx) GetLogger() *logger.Logger {
	return c.logger
}

func (c *ctx) SetConfig(cfg *config.Config) {
	c.config = cfg
}

func (c *ctx) GetConfig() *config.Config {
	return c.config
}

func (c *ctx) SetCri(_cri cri.CRI) {
	c.cri = _cri
}

func (c *ctx) GetCri() cri.CRI {
	return c.cri
}

func (c *ctx) Background() context.Context {
	return context.Background()
}