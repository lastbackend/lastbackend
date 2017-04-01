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
	"github.com/lastbackend/lastbackend/pkg/agent/config"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"golang.org/x/net/context"
)

var _ctx ctx

func Get() *ctx {
	return &_ctx
}

type ctx struct {
	Log    *logger.Logger
	Config *config.Config
}

func (c *ctx) New(config *config.Config) {
	c.Config = config
	c.Log = logger.New(c.Config.Debug)
}

func Background() context.Context {
	return context.Background()
}
