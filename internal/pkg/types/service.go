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

package types

const (
	DefaultServiceMemory   int64 = 128
	DefaultServiceReplicas int   = 1
)

type Service struct {
	System
	Meta   ServiceMeta   `json:"meta"`
	Status ServiceStatus `json:"status"`
	Spec   ServiceSpec   `json:"spec"`
}

type ServiceMap struct {
	System
	Items map[string]*Service
}

type ServiceList struct {
	System
	Items []*Service
}

type ServiceMeta struct {
	Meta
	Namespace string          `json:"namespace"`
	SelfLink  ServiceSelfLink `json:"self_link"`
	Endpoint  string          `json:"endpoint"`
	IP        string          `json:"ip"`
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
	s.Replicas = DefaultServiceReplicas
	s.Template.Volumes = make(SpecTemplateVolumeList, 0)
	s.Template.Containers = make(SpecTemplateContainers, 0)
}

func (s *Service) SelfLink() *ServiceSelfLink {
	return &s.Meta.SelfLink
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

func (s *ServiceSpec) GetResourceRequest() ResourceRequest {
	rr := ResourceRequest{}

	var (
		limitsRAM int64
		limitsCPU int64

		requestRAM int64
		requestCPU int64
	)

	for _, c := range s.Template.Containers {

		limitsCPU += c.Resources.Limits.CPU
		limitsRAM += c.Resources.Limits.RAM

		requestCPU += c.Resources.Request.CPU
		requestRAM += c.Resources.Request.RAM
	}

	if requestRAM > 0 {
		requestRAM = int64(s.Replicas) * requestRAM
		rr.Request.RAM = requestRAM
	}

	if requestCPU > 0 {
		requestCPU = int64(s.Replicas) * requestCPU
		rr.Request.CPU = requestCPU
	}

	if limitsRAM > 0 {
		limitsRAM = int64(s.Replicas) * limitsRAM
		rr.Limits.RAM = limitsRAM
	}

	if limitsCPU > 0 {
		limitsCPU = int64(s.Replicas) * limitsCPU
		rr.Limits.CPU = limitsCPU
	}

	return rr
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
