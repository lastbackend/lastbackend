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
	"github.com/lastbackend/lastbackend/pkg/daemon/config"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/http"
)

var context Context

func Get() *Context {
	return &context
}

type Context struct {
	Log              *logger.Logger
	TemplateRegistry *http.RawReq
	Storage          storage.IStorage
}

func (c *Context) Init(cfg *config.Config) {
	var err error

	config.Set(cfg)

	c.Log = logger.New(cfg.Debug, 9)

	// Initializing database
	c.Log.Info("Initializing daemon context")

	c.Storage, err = storage.Get(cfg.GetEtcdDB())
	if err != nil {
		c.Log.Panic(err)
	}

	if cfg.HttpServer.Port == 0 {
		cfg.HttpServer.Port = 3000
	}
}
