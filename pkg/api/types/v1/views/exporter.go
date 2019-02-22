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

package views

// Exporter - default node structure
// swagger:model views_ingress
type Exporter struct {
	Meta   ExporterMeta   `json:"meta"`
	Status ExporterStatus `json:"status"`
}

// ExporterList - node map list
// swagger:model views_ingress_list
type ExporterList map[string]*Exporter

// ExporterMeta - node metadata structure
// swagger:model views_ingress_meta
type ExporterMeta struct {
	Meta
	Version string `json:"version"`
}

// ExporterStatus - node state struct
// swagger:model views_ingress_status
type ExporterStatus struct {
	Ready bool `json:"ready"`
}

// swagger:model views_ingress_spec
type ExporterManifest struct {
	Meta ExporterManifestMeta `json:"meta"`
}

type ExporterManifestMeta struct {
	Initial bool `json:"initial"`
}
