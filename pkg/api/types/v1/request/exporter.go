//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

import (
	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

// swagger:model request_ingress_connect
type ExporterConnectOptions struct {
	Info    models.ExporterInfo   `json:"info"`
	Status  models.ExporterStatus `json:"status"`
	Network models.NetworkState   `json:"network"`
}

// swagger:ignore
// swagger:model request_ingress_remove
type ExporterRemoveOptions struct {
	Force bool `json:"force"`
}

type ExporterStatusOptions models.ExporterStatus
