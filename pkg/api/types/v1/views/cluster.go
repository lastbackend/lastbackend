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

import (
	"time"
)

type Cluster struct {
	ID    string       `json:"id"`
	Meta  ClusterMeta  `json:"meta"`
	State ClusterState `json:"state"`
}

type ClusterMeta struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Region      string            `json:"region"`
	Provider    string            `json:"provider"`
	Labels      map[string]string `json:"labels"`
	Created     time.Time         `json:"created"`
	Updated     time.Time         `json:"updated"`
}

type ClusterState struct {
	Nodes struct {
		Total   int `json:"total"`
		Online  int `json:"online"`
		Offline int `json:"offline"`
	} `json:"nodes"`
	Capacity  ClusterResources `json:"capacity"`
	Allocated ClusterResources `json:"allocated"`
	Deleted   bool             `json:"deleted"`
}
type ClusterList []*Cluster

type ClusterResources struct {
	Containers int   `json:"containers"`
	Pods       int   `json:"pods"`
	Memory     int64 `json:"memory"`
	Cpu        int   `json:"cpu"`
	Storage    int   `json:"storage"`
}
