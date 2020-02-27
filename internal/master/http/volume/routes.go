//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package volume

import (
	"github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/internal/util/http/middleware"
)

var Routes = []http.Route{
	// Route handlers
	{Path: "/namespace/{namespace}/volume", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Authenticate}, Handler: VolumeCreateH},
	{Path: "/namespace/{namespace}/volume", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: VolumeListH},
	{Path: "/namespace/{namespace}/volume/{volume}", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: VolumeInfoH},
	{Path: "/namespace/{namespace}/volume/{volume}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: VolumeUpdateH},
	{Path: "/namespace/{namespace}/volume/{volume}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Authenticate}, Handler: VolumeRemoveH},
}
