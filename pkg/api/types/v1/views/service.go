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

import "time"

type ServiceList []*Service

// ***************************************************
// SERVICE INFO MODEL
// ***************************************************

type Service struct {
	Meta        ServiceMeta            `json:"meta"`
	Stats       ServiceStats           `json:"stats"`
	Status      ServiceStatus          `json:"status"`
	Spec        ServiceSpec            `json:"spec"`
	Sources     ServiceSources         `json:"sources"`
	Deployments map[string]*Deployment `json:"deployments"`
}

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

type ServiceSources struct {
	// Image sources
	Image *ServiceSourcesImage `json:"image,omitempty"`
	// Deployment source lastbackend repo
	Repo *ServiceSourcesRepo `json:"repo,omitempty"`
}

type ServiceSourcesImage struct {
	// Image namespace name
	Namespace string `json:"namespace"`
	// Image tag
	Tag string `json:"tag"`
	// Hash
	Hash string `json:"hash"`
}

type ServiceSourcesRepo struct {
}

type ServiceStats struct {
	Memory  int64 `json:"memory"`
	Cpu     int64 `json:"cpu"`
	Network int64 `json:"network"`
}

type ServiceStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

type ServiceDeployment struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Pods    map[string]Pod
	Started time.Time `json:"started"`
}

type ServiceSpec struct {
	Replicas   int                `json:"replicas"`
	Meta       ServiceSpecMeta    `json:"meta"`
	Memory     int64              `json:"memory"`
	Image      string             `json:"image"`
	Entrypoint string             `json:"entrypoint"`
	Command    string             `json:"command"`
	EnvVars    []string           `json:"env"`
	Ports      []*ServiceSpecPort `json:"ports"`
}

type ServiceSpecMeta struct {
	ID        string `json:"id,omitempty"`
	ServiceID string `json:"service_id,omitempty"`
	Parent    string `json:"parent,omitempty"`
	Active    bool   `json:"active,omitempty"`
	Revision  int    `json:"revision,omitempty"`
}

type ServiceSpecPort struct {
	Protocol  string `json:"protocol"`
	Container int    `json:"internal"`
	Host      int    `json:"external"`
	Published bool   `json:"published"`
}

type ServiceEndpoint struct {
	Name      string `json:"name"`
	Technical bool   `json:"technical"`
	Main      bool   `json:"main"`
}
