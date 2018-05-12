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

package views

// Cluster represents cluster model for api
//
// swagger:model views_cluster
type Cluster struct {
	Meta   ClusterMeta   `json:"meta"`

	Status ClusterStatus `json:"status"`
}

// ClusterMeta represents meta of cluster model for api
//
// swagger:model views_cluster_meta
type ClusterMeta struct {
	// name of the cluster
	// example: cluster name
	Name        string            `json:"name"`

	// cluster description
	// example: this is cluster
	Description string            `json:"description"`

	// labels of the cluster
	Labels      map[string]string `json:"labels"`
}

// ClusterStatus represents status info of cluster model for api
//
// swagger:model views_cluster_status
type ClusterStatus struct {
	// cluster nodes info
	Nodes struct {
		// total number of nodes
		Total   int `json:"total"`

		// number of nodes online
		Online  int `json:"online"`

		// number of nodes offline
		Offline int `json:"offline"`
	} `json:"nodes"`

	Capacity  ClusterResources `json:"capacity"`

	Allocated ClusterResources `json:"allocated"`

	// is this cluster deleted
	// default: false
	Deleted   bool             `json:"deleted"`
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
	Containers int   `json:"containers"`

	// number of pods
	Pods       int   `json:"pods"`

	// memory volume
	Memory     int64 `json:"memory"`

	// number of cpu
	Cpu        int   `json:"cpu"`

	// storage volume
	Storage    int   `json:"storage"`
}
