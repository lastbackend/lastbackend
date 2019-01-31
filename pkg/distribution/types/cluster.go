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

type ClusterList struct {
	System
	Items []*Cluster
}
type ClusterMap struct {
	System
	Items map[string]*Cluster
}

type Cluster struct {
	System
	Meta   Meta          `json:"meta"`
	Status ClusterStatus `json:"status"`
	Spec   ClusterSpec   `json:"spec"`
}

type ClusterStatus struct {
	Nodes     ClusterStatusNodes     `json:"nodes"`
	Discovery ClusterStatusDiscovery `json:"discovery"`
	Ingress   ClusterStatusIngress   `json:"ingress"`
	Capacity  ClusterResources       `json:"capacity"`
	Allocated ClusterResources       `json:"allocated"`
	Deleted   bool                   `json:"deleted"`
}

type ClusterStatusNodes struct {
	Total   int `json:"total"`
	Online  int `json:"online"`
	Offline int `json:"offline"`
}

type ClusterStatusIngress struct {
	Total   int `json:"total"`
	Online  int `json:"online"`
	Offline int `json:"offline"`
}

type ClusterStatusDiscovery struct {
	Total   int `json:"total"`
	Online  int `json:"online"`
	Offline int `json:"offline"`
}

type ClusterResources struct {
	Containers int   `json:"containers"`
	Pods       int   `json:"pods"`
	RAM        int64 `json:"ram"`
	CPU        int   `json:"cpu"`
	Storage    int64 `json:"storage"`
}

type ClusterSpec struct {
}

// swagger:ignore
type ClusterCreateOptions struct {
	Description string `json:"description"`
}
