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
	"github.com/ararog/timeago"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"sort"
	"time"

	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/util/table"
)

type ServiceList []*Service
type Service struct {
	Meta        ServiceMeta   `json:"meta"`
	Spec        ServiceSpec   `json:"spec"`
	Status      ServiceStatus `json:"status"`
	Deployments DeploymentMap `json:"deployments"`
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
	Meta   DeploymentMeta     `json:"meta"`
	Status DeploymentStatus   `json:"status"`
	Spec   DeploymentSpec     `json:"spec"`
	Pods   map[string]PodView `json:"pods"`
}

type DeploymentMeta struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Created     time.Time         `json:"created"`
	Updated     time.Time         `json:"updated"`
}

type DeploymentStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
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

	t := table.New([]string{"NAME", "ENDPOINT", "STATUS", "REPLICAS"})
	t.VisibleHeader = true

	for _, s := range *sl {

		var data = map[string]interface{}{}

		data["NAME"] = s.Meta.Name
		data["ENDPOINT"] = s.Meta.Endpoint
		data["STATUS"] = s.Status.State

		t.AddRow(data)
	}
	println()
	t.Print()
	println()
}

func (s *Service) Print() {

	fmt.Printf("Name:\t\t%s/%s\n", s.Meta.Namespace, s.Meta.Name)
	if s.Meta.Description != types.EmptyString {
		fmt.Printf(" Description:\t%s\n", s.Meta.Description)
	}

	fmt.Printf("State:\t\t%s\n", s.Status.State)
	if s.Status.Message != types.EmptyString {
		fmt.Printf("Message:\t\t%s\n", s.Status.Message)
	}
	if s.Meta.Endpoint != types.EmptyString {
		fmt.Printf("Endpoint:\t%s\n", s.Meta.Endpoint)
	}

	created, _ := timeago.TimeAgoWithTime(time.Now(), s.Meta.Created)
	updated, _ := timeago.TimeAgoWithTime(time.Now(), s.Meta.Updated)

	fmt.Printf("Created:\t%s\n", created)
	fmt.Printf("Updated:\t%s\n", updated)

	var (
		labels = make([]string, 0, len(s.Meta.Labels))
		out string
	)

	for key := range s.Meta.Labels {
		labels = append(labels, key)
	}

	sort.Strings(labels) //sort by key
	for _, key := range labels {
		out += key + "=" + s.Meta.Labels[key] + " "
	}


	fmt.Printf("Labels:\t\t%s\n", out)
	println()

	println()
	if len(s.Deployments) > 0 {

		var states = make(map[string]int, 0)
		states[types.StateReady] = 0
		states[types.StateProvision] = 0
		states[types.EmptyString] = 0

		for _, d := range s.Deployments {
			switch d.Status.State {
			case types.StateReady:
				states[types.StateReady]++
				break
			case types.StateProvision:
				states[types.StateProvision]++
				break
			default:
				states[types.EmptyString]++
				break
			}
		}

		if states[types.StateReady] > 0 {

			fmt.Println("Active deployments:")
			println()

			for _, d := range s.Deployments {
				if d.Status.State == types.StateReady {
					d.Print()
				}
			}
		}

		if states[types.StateProvision] > 0 {
			println()
			fmt.Println("Provision deployments:")
			println()
			for _, d := range s.Deployments {
				if d.Status.State == types.StateProvision {
					d.Print()
					println()
				}
			}
		}

		if states[types.EmptyString] > 0 {
			println()
			fmt.Println("Inactive deployments:")
			println()
			for _, d := range s.Deployments {
				if d.Status.State != types.StateProvision && d.Status.State != types.StateReady {
					d.Print()
					println()
				}
			}
		}







	}
	println()
}

func (d *Deployment) Print() {

	fmt.Printf(" Name:\t\t%s\n", d.Meta.Name)
	if d.Meta.Description != types.EmptyString {
		fmt.Printf(" Description:\t%s\n", d.Meta.Description)
	}
	fmt.Printf(" State:\t\t%s\n", d.Status.State)
	if d.Status.Message != types.EmptyString {
		fmt.Printf(" Message:\t%s\n", d.Status.Message)
	}
	created, _ := timeago.TimeAgoWithTime(time.Now(), d.Meta.Created)
	updated, _ := timeago.TimeAgoWithTime(time.Now(), d.Meta.Updated)

	fmt.Printf(" Created:\t%s\n", created)
	fmt.Printf(" Updated:\t%s\n", updated)
	println()
	fmt.Printf(" Pods:\n")
	println()
	podTable := table.New([]string{"Name", "Ready", "Status", "Restarts", "Age"})
	podTable.VisibleHeader = true

	var (
		ids = make([]string, 0, len(d.Pods))

	)
	for key := range d.Pods {
		ids = append(ids, key)
	}

	sort.Strings(ids) //sort by key

	for _, id := range ids {
		p := d.Pods[id]

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
		podRow["Restarts"] = restarts
		podRow["Age"] = got
		podTable.AddRow(podRow)
	}

	podTable.Print()
	println()
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

		itd.Meta.Name = d.Meta.Name
		itd.Meta.Description = d.Meta.Description
		itd.Meta.Created = d.Meta.Created
		itd.Meta.Updated = d.Meta.Updated

		itd.Status.State = d.Status.State
		itd.Status.Message = d.Status.Message
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
