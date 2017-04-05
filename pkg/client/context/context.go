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
	"github.com/lastbackend/lastbackend/pkg/client/storage"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/util/http"
)

var context Context

func Get() *Context {
	return &context
}

func Mock() *Context {
	context.mock = true
	context.Log = logger.New(true, 8)
	context.Storage = new(storage.DB)

	return &context
}

type Context struct {
	Token   string
	Log     *logger.Logger
	HTTP    *http.RawReq
	Storage *storage.DB
	mock    bool
	// Other info for HTTP handlers can be here, like user UUID
}
