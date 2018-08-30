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

import (
	"fmt"
)

const (
	DEFAULT_SERVICE_MEMORY   int64 = 128
	DEFAULT_SERVICE_REPLICAS int   = 1
)

type Service struct {
	Runtime
	Meta   ServiceMeta   `json:"meta"`
	Status ServiceStatus `json:"status"`
	Spec   ServiceSpec   `json:"spec"`
}

type ServiceMap struct {
	Runtime
	Items map[string]*Service
}

type ServiceList struct {
	Runtime
	Items []*Service
}

type ServiceMeta struct {
	Meta
	Namespace string `json:"namespace"`
	SelfLink  string `json:"self_link"`
	Endpoint  string `json:"endpoint"`
	IP        string `json:"ip"`
}

type ServiceEndpoint struct {
	Name string `json:"name"`
	Main bool   `json:"main"`
}

type ServiceStatus struct {
	State   string               `json:"state"`
	Message string               `json:"message"`
	Network ServiceStatusNetwork `json:"network"`
}

type ServiceSpec struct {
	Replicas int          `json:"replicas" yaml:"replicas"`
	State    SpecState    `json:"state" yaml:"state"`
	Network  SpecNetwork  `json:"network" yaml:"network" `
	Strategy SpecStrategy `json:"strategy" yaml:"strategy"`
	Selector SpecSelector `json:"selector" yaml:"selector"`
	Template SpecTemplate `json:"template" yaml:"template"`
}

type ServiceStatusNetwork struct {
	IP string `json:"ip"`
}

type ServiceSpecStrategy struct {
	Type           string                            `json:"type"` // Rolling
	RollingOptions ServiceSpecStrategyRollingOptions `json:"rollingOptions"`
	Resources      ServiceSpecStrategyResources      `json:"resources"`
	Deadline       int                               `json:"deadline"`
}

type ServiceSpecStrategyResources struct{}

type ServiceSpecStrategyRollingOptions struct {
	PeriodUpdate   int `json:"period_update"`
	Interval       int `json:"interval"`
	Timeout        int `json:"timeout"`
	MaxUnavailable int `json:"max_unavailable"`
	MaxSurge       int `json:"max_surge"`
}

type ServiceReplicas struct {
	Total     int `json:"total"`
	Provision int `json:"provision"`
	Created   int `json:"created"`
	Started   int `json:"started"`
	Stopped   int `json:"stopped"`
	Errored   int `json:"errored"`
}

func (s *ServiceSpec) SetDefault() {
	s.Replicas = DEFAULT_SERVICE_REPLICAS
	s.Template.Volumes = make(SpecTemplateVolumeList, 0)
	s.Template.Containers = make(SpecTemplateContainers, 0)
}

func (s *Service) SelfLink() string {
	if s.Meta.SelfLink == "" {
		s.Meta.SelfLink = s.CreateSelfLink(s.Meta.Namespace, s.Meta.Name)
	}
	return s.Meta.SelfLink
}

func (s *Service) CreateSelfLink(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

type ServiceManifest struct {
	Meta ServiceMeta `json:"meta"`
}

type ServiceRemoveOptions struct {
	Force bool `json:"force"`
}

type ServiceImageSpec struct {
	Name   *string `json:"image"`
	Secret *string `json:"secret"`
}

type ServiceOptionsSpec struct {
	Replicas   *int                `json:"replicas"`
	Memory     *int64              `json:"memory,omitempty"`
	Entrypoint *string             `json:"entrypoint,omitempty"`
	Command    *string             `json:"command,omitempty"`
	EnvVars    *[]ServiceEnvOption `json:"env,omitempty"`
	Ports      map[uint16]string   `json:"ports,omitempty"`
}

type ServiceEnvOption struct {
	Name  string               `json:"name"`
	Value string               `json:"value"`
	From  ServiceEnvFromOption `json:"from"`
}

type ServiceEnvFromOption struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func NewServiceList() *ServiceList {
	dm := new(ServiceList)
	dm.Items = make([]*Service, 0)
	return dm
}

func NewServiceMap() *ServiceMap {
	dm := new(ServiceMap)
	dm.Items = make(map[string]*Service)
	return dm
}
