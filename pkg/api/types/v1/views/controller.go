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

package views

// Controller - default node structure
// swagger:model views_ingress
type Controller struct {
	Meta   ControllerMeta   `json:"meta"`
	Status ControllerStatus `json:"status"`
}

// ControllerList - node map list
// swagger:model views_ingress_list
type ControllerList map[string]*Controller

// ControllerMeta - node metadata structure
// swagger:model views_ingress_meta
type ControllerMeta struct {
	Meta
	Version string `json:"version"`
}

// ControllerStatus - node state struct
// swagger:model views_ingress_status
type ControllerStatus struct {
	Ready bool `json:"ready"`
}
