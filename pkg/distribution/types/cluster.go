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
	Runtime
	Items []*Cluster
}
type ClusterMap struct {
	Runtime
	Items map[string]*Cluster
}

type Cluster struct {
	Runtime
	Meta   Meta          `json:"meta"`
	Status ClusterStatus `json:"status"`
	Spec   ClusterSpec   `json:"spec"`
}

type ClusterStatus struct {
	Nodes     ClusterStatusNodes `json:"nodes"`
	Capacity  ClusterResources   `json:"capacity"`
	Allocated ClusterResources   `json:"allocated"`
	Deleted   bool               `json:"deleted"`
}

type ClusterStatusNodes struct {
	Total   int `json:"total"`
	Online  int `json:"online"`
	Offline int `json:"offline"`
}

type ClusterResources struct {
	Containers int   `json:"containers"`
	Pods       int   `json:"pods"`
	Memory     int64 `json:"memory"`
	Cpu        int   `json:"cpu"`
	Storage    int   `json:"storage"`
}

type ClusterSpec struct {
}

// swagger:ignore
type ClusterCreateOptions struct {
	Description string                  `json:"description"`
	Quotas      *NamespaceQuotasOptions `json:"quotas"`
}
