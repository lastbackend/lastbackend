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

package view

import (
	"fmt"

	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/util/table"
)

type ServiceList []*Service
type Service struct {
	Meta        ServiceMeta    `json:"meta"`
	Spec        ServiceSpec    `json:"spec"`
	Status      ServiceStatus  `json:"status"`
	Sources     ServiceSources `json:"sources"`
	Deployments DeploymentList `json:"deployments"`
}

type ServiceMeta struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Endpoint    string `json:"endpoint"`
}

type ServiceStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

type ServiceSources struct {
	Image *ServiceSourcesImage `json:"image,omitempty"`
}

type ServiceSourcesImage struct {
	Namespace string `json:"namespace"`
	Tag       string `json:"tag"`
}

type ServiceSpec struct {
	Memory  int64  `json:"memory"`
	Command string `json:"command"`
}

type DeploymentList []*Deployment
type Deployment struct {
	Name  string              `json:"name"`
	State DeploymentStateInfo `json:"state"`
	Spec  DeploymentSpec      `json:"spec"`
	Pods  []*PodView          `json:"pods"`
}

type PodView struct {
	Name   string     `json:"name"`
	Status *PodStatus `json:"status"`
}

type PodStatus struct {
	Containers *PodContainers `json:"containers"`
}

type PodContainers []*PodContainer
type PodContainer struct {
	ID string `json:"id"`
}

type DeploymentSpec struct {
	Replicas int `json:"replicas"`
}

type DeploymentStateInfo struct {
	Active bool `json:"active"`
}

func (sl *ServiceList) Print() {

	t := table.New([]string{"NAME", "SOURCES", "ENDPOINT", "STATUS", "CMD", "MEMORY", "REPLICAS"})
	t.VisibleHeader = true

	for _, s := range *sl {

		var data = map[string]interface{}{}

		data["NAME"] = s.Meta.Name
		data["SOURCES"] = s.Sources.String()
		data["ENDPOINT"] = s.Meta.Endpoint
		data["STATUS"] = s.Status.State
		data["CMD"] = s.Spec.Command
		data["MEMORY"] = s.Spec.Memory
		data["REPLICAS"] = s.Deployments.Replicas()

		t.AddRow(data)
	}
	println()
	t.Print()
	println()
}

func (s *Service) Print() {

	var data = map[string]interface{}{}

	data["NAME"] = s.Meta.Name
	data["SOURCES"] = s.Sources.String()
	data["ENDPOINT"] = s.Meta.Endpoint
	data["STATUS"] = s.Status.State

	if s.Status.Message != "" {
		data["STATUS"] = fmt.Sprintf("%s: %s", s.Status.State, s.Status.Message)
	}

	if s.Spec.Command != "" {
		data["CMD"] = s.Spec.Command
	}
	data["MEMORY"] = s.Spec.Memory

	if s.Deployments != nil {
		data["REPLICAS"] = s.Deployments.Replicas()
	}

	println()
	table.PrintHorizontal(data)
	println()
}

func (dl *DeploymentList) Replicas() int {
	for _, d := range *dl {
		if d.State.Active {
			return d.Spec.Replicas
		}
	}
	return 0
}

func (s *ServiceSources) String() string {
	if s.Image != nil && s.Image.Namespace != "" {
		return fmt.Sprintf("%s:%s",
			s.Image.Namespace, s.Image.Tag)
	}
	return ""
}

func FromApiServiceView(service *views.Service) *Service {
	var item = new(Service)
	item.Meta.Name = service.Meta.Name
	item.Meta.Description = service.Meta.Description
	item.Meta.Endpoint = service.Meta.Endpoint
	item.Status.State = service.Status.State
	item.Status.Message = service.Status.Message
	return item
}

func FromApiServiceListView(services *views.ServiceList) *ServiceList {
	var items = make(ServiceList, 0)
	for _, service := range *services {
		items = append(items, FromApiServiceView(service))
	}
	return &items
}
