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
	Meta    ServiceMeta  `json:"meta"`
	Pods    []v1.PodInfo `json:"pods,omitempty"`
	Spec    []SpecInfo   `json:"spec,omitempty"`
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
	// Meta labels
	Labels map[string]string `json:"labels"`
	// Meta created time
	Created time.Time `json:"created"`
	// Meta updated time
	Updated time.Time `json:"updated"`
}

type ServiceList []*Service
