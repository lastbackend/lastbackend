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
	"time"

	"github.com/ararog/timeago"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"github.com/lastbackend/lastbackend/pkg/util/table"
)

type ServiceList []*Service
type Service struct {
	Meta        ServiceMeta    `json:"meta"`
	Spec        ServiceSpec    `json:"spec"`
	Status      ServiceStatus  `json:"status"`
	Sources     ServiceSources `json:"sources"`
	Deployments DeploymentMap  `json:"deployments"`
}

type ServiceMeta struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Endpoint    string            `json:"endpoint"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels"`
	Created     time.Time         `json:"created"`
	Updated     time.Time         `json:"updated"`
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
type DeploymentMap map[string]*Deployment
type Deployment struct {
	Name  string             `json:"name"`
	State string             `json:"state"`
	Spec  DeploymentSpec     `json:"spec"`
	Pods  map[string]PodView `json:"pods"`
}

type PodView struct {
	Name    string     `json:"name"`
	Created time.Time  `json:"created"`
	Status  *PodStatus `json:"status"`
}

type PodStatus struct {
	State      string        `json:"state"`
	Message    string        `json:"message"`
	Containers PodContainers `json:"containers"`
}

type PodContainers []PodContainer
type PodContainer struct {
	ID      string `json:"id"`
	Ready   bool   `json:"ready"`
	Restart int    `json:"restared"`
}

type DeploymentSpec struct {
	Replicas int `json:"replicas"`
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
		//data["REPLICAS"] = s.Deployments.Replicas()

		t.AddRow(data)
	}
	println()
	t.Print()
	println()
}

func (s *Service) Print() {

	var data = map[string]interface{}{}

	data["Name"] = s.Meta.Name
	data["Namespace"] = s.Meta.Namespace
	data["Endpoint"] = s.Meta.Endpoint
	data["Status"] = s.Status.State
	data["Message"] = s.Status.Message
	data["Created"] = s.Meta.Created
	data["Updated"] = s.Meta.Updated

	//var labelList string
	//for key, l := range s.Meta.Labels {
	//	labelList += key + "=" + l + " "
	//}
	//data["LABELS"] = labelList

	//if s.Deployments != nil {
	//	data["REPLICAS"] = s.Deployments.Replicas()
	//}

	println()
	table.PrintHorizontal(data)
	println()
	if s.Deployments != nil {
		fmt.Println("Deployments:")
		println()
		s.Deployments.Print()
	}
	println()
}

//func (dl *DeploymentList) Replicas() int {
//	for _, d := range *dl {
//		if d.State.Active {
//			return d.Spec.Replicas
//		}
//	}
//	return 0
//}
//
func (dl *DeploymentMap) Print() {
	for _, d := range *dl {
		var data = map[string]interface{}{}

		data["Name"] = d.Name
		data["Status"] = d.State
		//data["REPLICAS"] = string(d.Spec.Replicas) + " updated | " + string(d.Spec.Replicas) + " total | " + string(d.Spec.Replicas) + " availible | 0 unavailible"

		table.PrintHorizontal(data)
		fmt.Println("\n Pods:")

		podTable := table.New([]string{"Name", "Ready", "Status", "Message", "Restarts", "Age"})
		podTable.VisibleHeader = true

		for _, p := range d.Pods {
			var ready, restarts int
			for _, c := range p.Status.Containers {
				if c.Ready {
					ready++
					restarts += c.Restart
				}
			}
			var podRow = map[string]interface{}{}
			got, _ := timeago.TimeAgoWithTime(time.Now(), p.Created)
			podRow["Name"] = p.Name
			podRow["Ready"] = string(converter.IntToString(ready) + "/" + converter.IntToString(len(p.Status.Containers)))
			podRow["Status"] = p.Status.State
			podRow["Message"] = p.Status.Message
			podRow["Restarts"] = restarts
			podRow["Age"] = got
			podTable.AddRow(podRow)
		}

		podTable.Print()
		println()
	}
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
	if service == nil {
		return nil
	}
	item.Meta.Name = service.Meta.Name
	item.Meta.Description = service.Meta.Description
	item.Meta.Namespace = service.Meta.Namespace
	item.Meta.Endpoint = service.Meta.Endpoint
	item.Meta.Labels = service.Meta.Labels
	item.Meta.Created = service.Meta.Created
	item.Meta.Updated = service.Meta.Updated

	item.Status.State = service.Status.State
	item.Status.Message = service.Status.Message

	item.Deployments = make(map[string]*Deployment, 0)

	for i, d := range service.Deployments {
		var itd Deployment

		itd.State = d.Status.State
		itd.Name = d.Meta.Name
		itd.Pods = make(map[string]PodView, 0)

		for j, p := range d.Pods {
			var pd = PodView{Status: &PodStatus{p.Status.State, p.Status.Message, PodContainers{}}}

			pd.Name = p.Meta.Name
			pd.Created = p.Meta.Created

			if pd.Status.Containers == nil {
				pd.Status.Containers = make(PodContainers, 0)
			}

			for _, c := range p.Status.Containers {
				var cn PodContainer

				cn.Ready = c.Ready
				cn.Restart = c.Restart

				pd.Status.Containers = append(pd.Status.Containers, cn)
			}

			itd.Pods[j] = pd
		}

		item.Deployments[i] = &itd
	}

	return item
}

func FromApiServiceListView(services *views.ServiceList) *ServiceList {
	var items = make(ServiceList, 0)
	for _, service := range *services {
		items = append(items, FromApiServiceView(service))
	}
	return &items
}
