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
	"github.com/satori/uuid"
)

const (
	DEFAULT_SERVICE_MEMORY   int64 = 128
	DEFAULT_SERVICE_REPLICAS int   = 1

	StepInitialized = "initialized"
	StepScheduled = "scheduled"
	StepPull = "pull"
	StepDestroyed = "destroyed"
	StepReady = "ready"

	StageInitialized = StepInitialized
	StageScheduled = StepScheduled
	StagePull = StepPull

	StageStarting = "starting"
	StageRunning = "running"
	StageStopped = "stopped"
	StageDestroyed = "destroyed"
	StageProvision = "provision"
	StageReady = "ready"
	StageCancel = "cancel"
	StageDestroy = "destroy"
	StageError = "error"
)

type Service struct {
	Meta        ServiceMeta            `json:"meta"`
	Status      ServiceStatus           `json:"status"`
	Spec        ServiceSpec            `json:"spec"`
	Deployments map[string]*Deployment `json:"deployments"`
}

type ServiceList map[string]*Service

type ServiceMeta struct {
	Meta
	Namespace string `json:"namespace"`
	SelfLink  string `json:"selflink"`
	Endpoint  string `json:"endpoint"`
}

type ServiceEndpoint struct {
	Name string `json:"name"`
	Main bool   `json:"main"`
}

type ServiceStatus struct {
	Stage string `json:"stage"`
	Message string `json:"message"`
}

type ServiceSpec struct {
	Meta     Meta          `json:"meta"`
	Replicas int           `json:"replicas"`
	State    SpecState     `json:"state"`
	Strategy SpecStrategy  `json:"strategy"`
	Triggers SpecTriggers  `json:"triggers"`
	Selector SpecSelector  `json:"selector"`
	Template SpecTemplate  `json:"template"`
	Quotas   ServiceQuotas `json:"quotas"`
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
	s.Meta.SetDefault()
	s.Meta.Name = uuid.NewV4().String()
	s.Replicas = int(1)
	s.Template.Volumes = make(SpecTemplateVolumes, 0)
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

	s.Meta.Name = uuid.NewV4().String()

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
	}

	if spec.Entrypoint != nil {
		c.Exec.Entrypoint = strings.Split(*spec.Entrypoint, " ")
	}

	if spec.Ports != nil {
		c.Ports = SpecTemplateContainerPorts{}
		for _, val := range *spec.Ports {
			c.Ports = append(c.Ports, SpecTemplateContainerPort{
				Protocol:      val.Protocol,
				ContainerPort: val.Internal,
			})
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
	}

	if spec.Memory != nil {
		c.Resources.Limits.RAM = *spec.Memory
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
		s.Meta.SelfLink = fmt.Sprintf("%s:%s", s.Meta.Namespace, s.Meta.Name)
	}
	return s.Meta.SelfLink
}

type ServiceCreateOptions struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Sources     *string             `json:"sources"`
	Replicas    *int                `json:"replicas"`
	Spec        *ServiceOptionsSpec `json:"spec"`
}

type ServiceUpdateOptions struct {
	Description *string             `json:"description"`
	Sources     *string             `json:"sources"`
	Spec        *ServiceOptionsSpec `json:"spec"`
}

type ServiceRemoveOptions struct {
	Force bool `json:"force"`
}

type ServiceOptionsSpec struct {
	Memory     *int64                    `json:"memory,omitempty"`
	Entrypoint *string                   `json:"entrypoint,omitempty"`
	Command    *string                   `json:"command,omitempty"`
	EnvVars    *[]string                 `json:"env,omitempty"`
	Ports      *[]ServiceOptionsSpecPort `json:"ports,omitempty"`
}

type ServiceOptionsSpecPort struct {
	Internal int    `json:"internal"`
	Protocol string `json:"protocol"`
}
