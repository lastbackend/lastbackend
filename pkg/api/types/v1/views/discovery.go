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

import "github.com/lastbackend/lastbackend/internal/pkg/models"

// Discovery - default node structure
// swagger:model views_ingress
type Discovery struct {
	Meta   DiscoveryMeta   `json:"meta"`
	Status DiscoveryStatus `json:"status"`
}

// DiscoveryList - node map list
// swagger:model views_ingress_list
type DiscoveryList map[string]*Discovery

// DiscoveryMeta - node metadata structure
// swagger:model views_ingress_meta
type DiscoveryMeta struct {
	Meta
	Version string `json:"version"`
}

// DiscoveryStatus - node state struct
// swagger:model views_ingress_status
type DiscoveryStatus struct {
	Ready bool `json:"ready"`
}

// swagger:model views_ingress_spec
type DiscoveryManifest struct {
	Meta    DiscoveryManifestMeta             `json:"meta"`
	Subnets map[string]*models.SubnetManifest `json:"subnets"`
}

type DiscoveryManifestMeta struct {
	Initial bool `json:"initial"`
}
