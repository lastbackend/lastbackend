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

// Cluster represents cluster model for api
//
// swagger:model views_cluster
type Cluster struct {
	Meta
	Status ClusterStatus `json:"status"`
}

// ClusterStatus represents status info of cluster model for api
//
// swagger:model views_cluster_status
type ClusterStatus struct {
	// cluster nodes info
	Nodes struct {
		// total number of nodes
		Total int `json:"total"`

		// number of nodes online
		Online int `json:"online"`

		// number of nodes offline
		Offline int `json:"offline"`
	} `json:"nodes"`

	Namespace []Namespace `json:"namespace"`
	Service   []Service   `json:"service"`
	Route     []Route     `json:"route"`
	Task      []Task      `json:"task"`

	Capacity ClusterResources `json:"capacity"`

	Allocated ClusterResources `json:"allocated"`
}

// swagger:ignore
// ClusterList is a list of cluster models for api
//
// swagger:model views_cluster_list
type ClusterList []*Cluster

// ClusterResources represents quantity of cluster resources
//
// swagger:model views_cluster_resources
type ClusterResources struct {
	// number of containers
	Containers int `json:"containers"`
	// number of pods
	Pods int `json:"pods"`
	// ram size
	RAM string `json:"ram"`
	// cpu  size
	Cpu string `json:"cpu"`
	// storage size
	Storage string `json:"storage"`
}
