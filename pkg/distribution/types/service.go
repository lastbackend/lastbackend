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
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
)

const (
	DEFAULT_SERVICE_MEMORY   int64 = 128
	DEFAULT_SERVICE_REPLICAS int   = 1
)

type Service struct {
	Meta        ServiceMeta            `json:"meta"`
	State       ServiceState           `json:"state"`
	Spec        ServiceSpec            `json:"spec"`
	Deployments map[string]*Deployment `json:"deployments"`
}

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

type ServiceState struct {
	Ready     bool `json:"ready"`
	Provision bool `json:"provision"`
	Destroy   bool `json:"destroy"`
}

type ServiceSpec struct {
	Replicas int           `json:"replicas"`
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

type ServiceSpecStrategyResources struct {
}

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
	s.Replicas = int(1)
	s.Template.Volumes = make(SpecTemplateVolumes, 0)
	s.Template.Containers = make(SpecTemplateContainers, 0)
	s.Triggers = make(SpecTriggers, 0)
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

type ServiceCreateOptions struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Sources     *string             `json:"sources"`
	Replicas    *int                `json:"replicas"`
	Spec        *ServiceOptionsSpec `json:"spec"`
}

type ServiceUpdateOptions struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Sources     *string             `json:"sources"`
	Spec        *ServiceOptionsSpec `json:"spec"`
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

func (s *ServiceCreateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Service: decode and validate data for creating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Service: decode and validate data for creating err: %s", err)
		return errors.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Service: convert struct from json err: %s", err)
		return errors.New("service").IncorrectJSON(err)
	}

	if s.Name == nil {
		log.V(logLevel).Error("Request: Service: parameter name can not be empty")
		return errors.New("service").BadParameter("name")
	}

	*s.Name = strings.ToLower(*s.Name)

	if len(*s.Name) < 4 && len(*s.Name) > 64 && !validator.IsServiceName(*s.Name) {
		log.V(logLevel).Error("Request: Service: parameter name not valid")
		return errors.New("service").BadParameter("name")
	}

	if s.Spec.Memory == nil {
		memory := int64(DEFAULT_SERVICE_MEMORY)
		s.Spec.Memory = &memory
	}

	if s.Replicas == nil {
		replicas := int(DEFAULT_SERVICE_REPLICAS)
		s.Replicas = &replicas
	}

	return nil
}

func (s *ServiceUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Service: decode and validate data for creating")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Service: decode and validate data for creating err: %s", err)
		return errors.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Service: convert struct from json err: %s", err)
		return errors.New("service").IncorrectJSON(err)
	}

	return nil
}
