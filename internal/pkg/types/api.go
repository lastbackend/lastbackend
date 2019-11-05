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

// swagger:ignore
type API struct {
	System
	Meta   APIMeta   `json:"meta"`
	Status APIStatus `json:"status"`
	Spec   APISpec   `json:"spec"`
}

type APIList struct {
	System
	Items []*API
}

type APIMap struct {
	System
	Items map[string]*API
}

// swagger:ignore
type APIMeta struct {
	Meta
	SelfLink APISelfLink `json:"self_link"`
	Node     string      `json:"node"`
}

// swagger:model types_discovery_info
type APIInfo struct {
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
type APIStatus struct {
	IP     string `json:"ip"`
	Port   uint16 `json:"port"`
	Ready  bool   `json:"ready"`
	Online bool   `json:"online"`
}

// swagger:ignore
type APISpec struct {
}

func (n *API) SelfLink() *APISelfLink {
	return &n.Meta.SelfLink
}

func NewAPIList() *APIList {
	dm := new(APIList)
	dm.Items = make([]*API, 0)
	return dm
}

func NewAPIMap() *APIMap {
	dm := new(APIMap)
	dm.Items = make(map[string]*API)
	return dm
}
