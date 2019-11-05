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
type Exporter struct {
	System
	Meta   ExporterMeta   `json:"meta"`
	Status ExporterStatus `json:"status"`
	Spec   ExporterSpec   `json:"spec"`
}

type ExporterList struct {
	System
	Items []*Exporter
}

type ExporterMap struct {
	System
	Items map[string]*Exporter
}

// swagger:ignore
type ExporterMeta struct {
	Meta
	SelfLink ExporterSelfLink `json:"self_link"`
	Node     string           `json:"node"`
}

type ExporterInfo struct {
	Type         string `json:"type"`
	Version      string `json:"version"`
	Hostname     string `json:"hostname"`
	Architecture string `json:"architecture"`

	OSName string `json:"os_name"`
	OSType string `json:"os_type"`

	// RewriteIP - need to set true if you want to use an external ip
	ExternalIP string `json:"external_ip"`
	InternalIP string `json:"internal_ip"`
}

// swagger:model types_ingress_status
type ExporterStatus struct {
	Master   bool                   `json:"master"`
	Ready    bool                   `json:"ready"`
	Online   bool                   `json:"online"`
	Http     ExporterStatusHttp     `json:"http"`
	Listener ExporterStatusListener `json:"listener"`
}

type ExporterStatusHttp struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

type ExporterStatusListener struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

type ExporterSpec struct {
}

// swagger:ignore
func (n *Exporter) SelfLink() *ExporterSelfLink {
	return &n.Meta.SelfLink
}

func NewExporterList() *ExporterList {
	dm := new(ExporterList)
	dm.Items = make([]*Exporter, 0)
	return dm
}

func NewExporterMap() *ExporterMap {
	dm := new(ExporterMap)
	dm.Items = make(map[string]*Exporter)
	return dm
}
