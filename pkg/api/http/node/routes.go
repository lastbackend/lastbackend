//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package node

import (
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/lastbackend/lastbackend/pkg/util/http/middleware"
)

var Routes = []http.Route{
	{Path: "/cluster/node", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeListH},
	{Path: "/cluster/node/{node}", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeInfoH},
	{Path: "/cluster/node/{node}/spec", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeGetSpecH},
	{Path: "/cluster/node/{node}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeRemoveH},
	{Path: "/cluster/node/{node}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeConnectH},
	{Path: "/cluster/node/{node}/meta", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeSetMetaH},
	{Path: "/cluster/node/{node}/status", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeSetStatusH},
	{Path: "/cluster/node/{node}/status/pod/{pod}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeSetPodStatusH},
	{Path: "/cluster/node/{node}/status/volume/{pod}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeSetVolumeStatusH},
	{Path: "/cluster/node/{node}/status/route/{pod}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeSetRouteStatusH},
}
