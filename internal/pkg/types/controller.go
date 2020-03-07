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

package types

// swagger:ignore
type Controller struct {
	System
	Meta   ControllerMeta   `json:"meta"`
	Status ControllerStatus `json:"status"`
	Spec   ControllerSpec   `json:"spec"`
}

type ControllerList struct {
	System
	Items []*Controller
}

type ControllerMap struct {
	System
	Items map[string]*Controller
}

// swagger:ignore
type ControllerMeta struct {
	Meta
	SelfLink ControllerSelfLink `json:"self_link"`
	Node     string             `json:"node"`
}

// swagger:model types_discovery_info
type ControllerInfo struct {
	Version      string `json:"version"`
	Hostname     string `json:"hostname"`
	Architecture string `json:"architecture"`

	OSName string `json:"os_name"`
	OSType string `json:"os_type"`

	// RewriteIP - need to set true if you want to use an external ip
	ExternalIP string `json:"external_ip"`
	InternalIP string `json:"internal_ip"`
}

// swagger:model types_discovery_status
type ControllerStatus struct {
	IP     string `json:"ip"`
	Port   uint16 `json:"port"`
	Ready  bool   `json:"ready"`
	Online bool   `json:"online"`
}

// swagger:ignore
type ControllerSpec struct {
}

func (n *Controller) SelfLink() *ControllerSelfLink {
	return &n.Meta.SelfLink
}

func NewControllerList() *ControllerList {
	dm := new(ControllerList)
	dm.Items = make([]*Controller, 0)
	return dm
}

func NewControllerMap() *ControllerMap {
	dm := new(ControllerMap)
	dm.Items = make(map[string]*Controller)
	return dm
}
