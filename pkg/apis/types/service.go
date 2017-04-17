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
	"errors"
)

type ServiceList []*Service

type Service struct {
	// Service Meta
	Meta ServiceMeta `json:"meta"`
	// Service custom domains
	Domains []string `json:"domains"`
	// Service source info
	Source ServiceSource `json:"source"`
	// Service config info
	Config ServiceConfig `json:"config"`
	// Pods list
	Pods map[string]*Pod `json:"pods"`
}

type ServiceMeta struct {
	Meta
	// Service replicas
	Replicas int `json:"replicas"`
	// Service state
	State ServiceState `json:"state"`
	// Service namespace
	Namespace string `json:"namespace"`
	// Service region
	Region string `json:"region,omitempty"`
}

type ServiceState struct {
	// Service state
	State string `json:"state"`
	// Service status
	Status string `json:"status"`
	// Service pods state
	Pods ServicePodsState `json:"pods"`
}

type ServicePodsState struct {
	// Total containers
	Total int `json:"total"`
	// Total running containers
	Running int `json:"running"`
	// Total created containers
	Created int `json:"created"`
	// Total stopped containers
	Stopped int `json:"stopped"`
	// Total errored containers
	Errored int `json:"errored"`
}

type ServiceSource struct {
	Hub    string `json:"hub"`
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
}

type ServiceConfig struct {
	Replicas   int      `json:"replicas"`
	Memory     int64    `json:"memory"`
	Entrypoint []string `json:"entrypoint"`
	Image      string   `json:"image"`
	Command    []string `json:"command"`
	EnvVars    []string `json:"env"`
	Ports      []Port   `json:"ports"`
}

func (c *ServiceConfig) Update(patch *ServiceConfig) error {

	if patch.Replicas < 0 {
		return errors.New("The value of the `replicas` parameter must be at least 1")
	}
	c.Replicas = patch.Replicas

	if patch.Memory < 32 {
		return errors.New("The value of the `memory` parameter must be at least 32")
	}
	c.Memory = patch.Memory

	c.Entrypoint = patch.Entrypoint
	c.Image = patch.Image
	c.Command = patch.Command

	c.Ports = patch.Ports

	// TODO: Check valid format env params
	c.EnvVars = patch.EnvVars

	return nil
}

func (ServiceConfig) GetDefault() ServiceConfig {
	var config = ServiceConfig{}
	config.Replicas = 1
	config.Memory = 256
	return config
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
	Replicas    *int32             `json:"replicas,omitempty" yaml:"replicas,omitempty"`
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
