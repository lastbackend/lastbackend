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

package types

type NodeStatusEvent struct {
	Event string `json:"event"`
	Node  string `json:"node"`
	Ready bool   `json:"ready"`
}

type PodManifestEvent struct {
	Event    string      `json:"event"`
	Node     string      `json:"node"`
	Name     string      `json:"name"`
	Manifest PodManifest `json:"manifest"`
}

type VolumeManifestEvent struct {
	Event    string         `json:"event"`
	Node     string         `json:"node"`
	Name     string         `json:"name"`
	Manifest VolumeManifest `json:"manifest"`
}

type NetworkManifestEvent struct {
	Event    string          `json:"event"`
	Manifest NetworkManifest `json:"manifest"`
}

type EndpointManifestEvent struct {
	Event    string           `json:"event"`
	Manifest EndpointManifest `json:"manifest"`
}

type IngressEvent struct {
	Event   string  `json:"event"`
	Name    string  `json:"name"`
	Ingress Ingress `json:"ingress"`
	Ready   bool    `json:"ready"`
}

type IngressRouteEvent struct {
	Event string    `json:"event"`
	Name  string    `json:"name"`
	Route RouteSpec `json:"spec"`
}
