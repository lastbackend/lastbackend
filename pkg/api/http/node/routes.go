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
	{Path: "/cluster/node/{node}", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeGetH},
	{Path: "/cluster/node/{node}/spec", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeGetSpecH},
	{Path: "/cluster/node/{node}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeRemoveH},
	{Path: "/cluster/node/{node}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeUpdateH},
	{Path: "/cluster/node/{node}/info", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeSetInfoH},
	{Path: "/cluster/node/{node}/state", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeSetStateH},
	{Path: "/cluster/node/{node}/state/pod/{pod}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeSetPodStateH},
	{Path: "/cluster/node/{node}/state/volume/{pod}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeSetVolumeStateH},
	{Path: "/cluster/node/{node}/state/route/{pod}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: NodeSetRouteStateH},
}
