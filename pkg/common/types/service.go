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

package types

import (
	"encoding/json"
)

type ServiceList []*Service

type Service struct {
	// Service Meta
	Meta ServiceMeta `json:"meta"`
	// Service state
	State ServiceState `json:"state"`
	// Service source info
	Source ServiceSource `json:"source"`
	// Service config info
	Spec map[string]*ServiceSpec `json:"spec"`
	// Pods list
	Pods map[string]*Pod `json:"pods"`
	//Service DNS
	DNS ServiceDNS `json:"dns"`
}

type ServiceCreateSpec struct {
	// Service Meta
	Meta ServiceMeta `json:"meta"`
	// Service source info
	Source ServiceSource `json:"source"`
	// Service config info
	Config ServiceSpec `json:"config"`
}

type ServiceUpdateSpec struct {
	// Service Meta
	Meta ServiceMeta `json:"meta"`
	// Service source info
	Source ServiceSource `json:"source"`
	// Service config info
	Config ServiceSpec `json:"config"`
}

type ServiceMeta struct {
	Meta
	// Service replicas
	Replicas int `json:"replicas"`
	// Service namespace
	Namespace string `json:"namespace"`
	// Service region
	Region string `json:"region,omitempty"`
	//Service hook
	Hook string `json:"hook"`
}

type ServiceState struct {
	// Service state
	State string `json:"state"`
	// Service status
	Status string `json:"status,omitempty"`
	// Service resources
	Resources ServiceResourcesState `json:"resources"`
	// Replicas state
	Replicas ServiceReplicasState `json:"replicas"`
}

type ServiceDNS struct {
	// Service primary dns
	Primary string `json:"primary"`
	// Service secondary dns
	Secondary string `json:"secondary,omitempty"`
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

type ServiceSource struct {
	Hub    string `json:"hub"`
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
}

type SpecMeta struct {
	Meta
	ID       string `json:"id"`
	Parent   string `json:"parent"`
	Revision int    `json:"revision"`
}

type ServiceSpec struct {
	Meta       SpecMeta `json:"meta"`
	Memory     int64    `json:"memory"`
	Entrypoint []string `json:"entrypoint"`
	Image      string   `json:"image"`
	Command    []string `json:"command"`
	EnvVars    []string `json:"env"`
	Ports      []Port   `json:"ports"`
}

type Port struct {
	Protocol  string `json:"protocol"`
	Container int    `json:"internal"`
	Host      int    `json:"external"`
	Published bool   `json:"published"`
}

func (s *Service) ToJson() ([]byte, error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *ServiceList) ToJson() ([]byte, error) {

	if s == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

type ServiceUpdateConfig struct {
	Name        *string            `json:"name,omitempty" yaml:"name,omitempty"`
	Description *string            `json:"description,omitempty" yaml:"description,omitempty"`
	Replicas    *int               `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	Containers  *[]ContainerConfig `json:"containers,omitempty" yaml:"containers,omitempty"`
}

type ContainerConfig struct {
	Image      string   `json:"image" yaml:"image"`
	Name       string   `json:"name" yaml:"name"`
	WorkingDir string   `json:"workdir" yaml:"workdir"`
	Command    []string `json:"command" yaml:"command"`
	Args       []string `json:"args" yaml:"args"`
	EnvVars    []string `json:"env" yaml:"env"`
	Ports      []Port   `json:"ports" yaml:"ports"`
}
