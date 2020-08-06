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

import (
	"time"
)

// swagger:model views_service_list
type ServiceList []*Service

// ***************************************************
// SERVICE INFO MODEL
// ***************************************************

// swagger:model views_service
type Service struct {
	Meta        ServiceMeta            `json:"meta"`
	Status      ServiceStatus          `json:"status"`
	Spec        ServiceSpec            `json:"spec"`
	Deployments map[string]*Deployment `json:"deployments,omitempty"`
}

// swagger:model views_service_meta
type ServiceMeta struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Description string            `json:"description"`
	SelfLink    string            `json:"self_link"`
	Endpoint    string            `json:"endpoint"`
	Replicas    int               `json:"replicas"`
	Labels      map[string]string `json:"labels"`
	Created     time.Time         `json:"created"`
	Updated     time.Time         `json:"updated"`
}

// swagger:ignore
type ServiceImage struct {
	// Name namespace name
	Namespace string `json:"namespace"`
	// Name tag
	Tag string `json:"tag"`
	// Hash
	Hash string `json:"hash"`
}

// swagger:ignore
type ServiceSourcesRepo struct {
}

// swagger:model views_service_stats
type ServiceStats struct {
	Memory  int64 `json:"memory"`
	Cpu     int64 `json:"cpu"`
	Network int64 `json:"network"`
}

// swagger:model views_service_status
type ServiceStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

// swagger:ignore
type ServiceDeployment struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Pods    map[string]Pod
	Started time.Time `json:"started"`
}

// swagger:model views_service_spec
type ServiceSpec struct {
	Selector ManifestSpecSelector `json:"selector,omitempty" yaml:"selector,omitempty"`
	Replicas int                  `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	Network  ManifestSpecNetwork  `json:"network,omitempty" yaml:"network,omitempty"`
	Strategy ManifestSpecStrategy `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Template ManifestSpecTemplate `json:"template,omitempty" yaml:"template,omitempty"`
}

type ServiceTemplateSpec struct {
}

// swagger:model views_service_spec_meta
type ServiceSpecMeta struct {
	ID        string `json:"id,omitempty"`
	ServiceID string `json:"service_id,omitempty"`
	Parent    string `json:"parent,omitempty"`
	Active    bool   `json:"active,omitempty"`
	Revision  int    `json:"revision,omitempty"`
}

// swagger:model views_service_spec_port
type ServiceSpecPort struct {
	Protocol  string `json:"protocol"`
	Container uint16 `json:"internal"`
	Host      uint16 `json:"external"`
	Published bool   `json:"published"`
}

// swagger:ignore
type ServiceEndpoint struct {
	Name      string `json:"name"`
	Technical bool   `json:"technical"`
	Main      bool   `json:"main"`
}
