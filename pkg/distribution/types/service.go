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
	"strings"

	"fmt"

	"time"

	"github.com/lastbackend/lastbackend/pkg/util/network"
	"github.com/lastbackend/dynamic/pkg/log"
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
	Replicas int          `json:"replicas"`
	State    SpecState    `json:"state"`
	Network  SpecNetwork  `json:"network"`
	Strategy SpecStrategy `json:"strategy"`
	Triggers SpecTriggers `json:"triggers"`
	Selector SpecSelector `json:"selector"`
	Template SpecTemplate `json:"template"`
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

type ServiceQuotas struct {
	RAM *int64 `json:"ram"`
}

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
	s.Triggers = make(SpecTriggers, 0)
}

func (s *ServiceSpec) Update(spec *ServiceOptionsSpec) {

	if spec == nil {
		return
	}

	var (
		f bool
		i int
	)

	if spec.Replicas != nil {
		log.Infof("set replicas: %d", *spec.Replicas)
		s.Replicas = *spec.Replicas
	}

	c := SpecTemplateContainer{}
	for i, t := range s.Template.Containers {
		if t.Role == ContainerRolePrimary {
			c = s.Template.Containers[i]
			f = true
		}
	}

	if !f {
		c.SetDefault()
		c.Role = ContainerRolePrimary
	}

	if spec.Command != nil {
		c.Exec.Command = strings.Split(*spec.Command, " ")
		s.Template.Updated = time.Now()
	}

	if spec.Entrypoint != nil {
		c.Exec.Entrypoint = strings.Split(*spec.Entrypoint, " ")
		s.Template.Updated = time.Now()
	}

	// TODO: update for multi-container pod
	if spec.Ports != nil {

		s.Network.Ports = make(map[uint16]string, 0)
		c.Ports = SpecTemplateContainerPorts{}

		for pt, pm := range spec.Ports {
			port, proto, err := network.ParsePortMap(pm)
			if err != nil {
				continue
			}

			c.Ports = append(c.Ports, SpecTemplateContainerPort{
				Protocol:      proto,
				ContainerPort: port,
			})

			s.Network.Ports[pt] = fmt.Sprintf("%d/%s", port, proto)
			s.Network.Updated = time.Now()
		}

	}

	if spec.EnvVars != nil {
		c.EnvVars = SpecTemplateContainerEnvs{}

		for _, e := range *spec.EnvVars {
			match := strings.Split(e, "=")
			env := SpecTemplateContainerEnv{Name: match[0]}
			if len(match) == 2 {
				env.Value = match[1]
			}
			c.EnvVars = append(c.EnvVars, env)
		}
		s.Template.Updated = time.Now()
	}

	if spec.Memory != nil && *spec.Memory != c.Resources.Limits.RAM {
		c.Resources.Limits.RAM = *spec.Memory
		s.Template.Updated = time.Now()
	}

	if !f {
		s.Template.Containers = append(s.Template.Containers, c)
	} else {
		s.Template.Containers[i] = c
	}

}

type ServiceSources struct {
	// Image sources
	Image ServiceSourcesImage `json:"image"`
	// Deployment source lastbackend repo
	Repo ServiceSourcesRepo `json:"repo"`
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
	// Deployment source lastbackend repo ID
	ID string `json:"id"`
	// Branch info
	Tag string `json:"tag"`
	// Build sources info
	Build string `json:"build"`
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

type ServiceCreateOptions struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Image       *string             `json:"sources"`
	Spec        *ServiceOptionsSpec `json:"spec"`
}

type ServiceUpdateOptions struct {
	Description *string             `json:"description"`
	Spec        *ServiceOptionsSpec `json:"spec"`
}

type ServiceRemoveOptions struct {
	Force bool `json:"force"`
}

type ServiceOptionsSpec struct {
	Replicas   *int              `json:"replicas"`
	Memory     *int64            `json:"memory,omitempty"`
	Entrypoint *string           `json:"entrypoint,omitempty"`
	Command    *string           `json:"command,omitempty"`
	EnvVars    *[]string         `json:"env,omitempty"`
	Ports      map[uint16]string `json:"ports,omitempty"`
}


func NewServiceList () *ServiceList {
	dm := new(ServiceList)
	dm.Items = make([]*Service, 0)
	return dm
}

func NewServiceMap () *ServiceMap {
	dm := new(ServiceMap)
	dm.Items = make(map[string]*Service)
	return dm
}