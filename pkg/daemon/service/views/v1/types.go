//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package v1

import (
	"github.com/lastbackend/lastbackend/pkg/daemon/pod/views/v1"
	"time"
)

type Service struct {
	Meta  ServiceMeta  `json:"meta"`
	State ServiceState `json:"state"`
	Pods  []v1.PodInfo `json:"pods,omitempty"`
	Spec  []SpecInfo   `json:"spec,omitempty"`
}

type ServiceState struct {
	// Service state
	State string `json:"state"`
	// Service status
	Status string `json:"status"`
	// Service resources
	Resources ServiceResourcesState `json:"resources"`
	// Replicas state
	Replicas ServiceReplicasState `json:"replicas"`
}

type ServiceResourcesState struct {
	// Total containers
	Memory int `json:"memory"`
}

type ServiceReplicasState struct {
	// Total pods
	Total int `json:"total"`
	// Total pods provision
	Provision int `json:"provision"`
	// Total pods provision
	Ready int `json:"ready"`
	// Total running pods
	Running int `json:"running"`
	// Total created pods
	Created int `json:"created"`
	// Total stopped pods
	Stopped int `json:"stopped"`
	// Total errored pods
	Errored int `json:"errored"`
}

type ServiceMeta struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Region      string    `json:"region"`
	Replicas    int       `json:"replicas,omitempty"`
	Namespace   string    `json:"namespace"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type SpecInfo struct {
	Meta    SpecMeta `json:"meta"`
	Memory  int64    `json:"memory"`
	Command string   `json:"command"`
	Image   string   `json:"image"`
	EnvVars []string `json:"env"`
	Ports   []Port   `json:"ports"`
}

type Port struct {
	Protocol  string `json:"protocol"`
	External  int    `json:"external"`
	Internal  int    `json:"internal"`
	Published bool   `json:"published"`
}

type SpecMeta struct {
	// Meta id
	ID string `json:"id"`
	// Parent meta id
	Parent string `json:"parent"`
	// Revision version
	Revision int `json:"revision"`
	// Meta labels
	Labels map[string]string `json:"labels"`
	// Meta created time
	Created time.Time `json:"created"`
	// Meta updated time
	Updated time.Time `json:"updated"`
}

type ServiceList []*Service
