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
	"github.com/lastbackend/lastbackend/pkg/agent/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/agent/storage"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"golang.org/x/net/context"
	"github.com/lastbackend/lastbackend/pkg/agent/events/listener"
)

var _ctx ctx

func Get() *ctx {
	return &_ctx
}

type ctx struct {
	cri     cri.CRI
	logger  *logger.Logger
	config  *config.Config
	storage *storage.Storage
	http    *http.RawReq
	event   *listener.EventListener
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

func (c *ctx) SetStorage(s *storage.Storage) {
	c.storage = s
}

func (c *ctx) GetStorage() *storage.Storage {
	return c.storage
}

func (c *ctx) SetHttpClient(http *http.RawReq) {
	c.http = http
}

func (c *ctx) GetHttpClient() *http.RawReq {
	return c.http
}

func (c *ctx) SetEventListener (el *listener.EventListener) {
	c.event = el
}

func (c *ctx) GetEventListener () *listener.EventListener {
	return c.event
}

func Background() context.Context {
	return context.Background()
}
