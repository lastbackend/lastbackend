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

package types

const (
	KindCluster    = "cluster"
	KindAPIServer  = "api"
	KindController = "controller"
	KindDiscovery  = "discovery"
	KindIngress    = "ingress"
	KindExporter   = "exporter"
	KindNode       = "node"
)

type Process struct {
	System
	// Process Meta
	Meta ProcessMeta `json:"meta"`
	// Process status
	Status ProcessStatus `json:"status"`
}

type ProcessMeta struct {
	// Include default Meta struct
	Meta `json:"id" `

	ID string `json:"id" `

	// Process PID
	PID int `json:"pid" `

	// Process Master state
	Lead bool `json:"lead" `
	// Process Slave state
	Slave bool `json:"slave" `

	// Process registered type
	Kind string `json:"kind" `
	// Process registered hostname
	Hostname string `json:"hostname" `
}

type ProcessStatus struct{}
