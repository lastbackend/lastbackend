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

package request

import "github.com/lastbackend/lastbackend/pkg/distribution/types"

// swagger:model request_ingress_connect
type IngressConnectOptions struct {
	Info    types.IngressInfo   `json:"info"`
	Status  types.IngressStatus `json:"status"`
	Network types.NetworkState  `json:"network"`
}

// swagger:ignore
// swagger:model request_ingress_remove
type IngressRemoveOptions struct {
	Force bool `json:"force"`
}

type IngressStatusOptions struct {
	Status types.IngressStatus           `json:"status"`
	Routes map[string]*types.RouteStatus `json:"routes"`
}
