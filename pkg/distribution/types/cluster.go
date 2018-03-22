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

const (
	CentralUSRegions = "CU"
	WestEuropeRegion = "WE"
	EastAsiaRegion   = "EA"
)

type ClusterList []*Cluster

type Cluster struct {
	Meta   ClusterMeta   `json:"meta"`
	Status ClusterStatus `json:"status"`
	Quotas ClusterQuotas `json:"quotas"`
}

type ClusterMeta struct {
	Meta

	Region   string `json:"region"`
	Token    string `json:"token"`
	Provider string `json:"provider"`
	Shared   bool   `json:"shared"`
	Main     bool   `json:"main"`
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

type ClusterQuotas struct {
}
