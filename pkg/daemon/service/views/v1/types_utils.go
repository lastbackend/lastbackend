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

package v1

import (
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/pod/views/v1"
	"github.com/lastbackend/lastbackend/pkg/util/table"
	"strings"
)

func New(obj *types.Service) *Service {
	s := Service{}

	s.Meta.Name = obj.Meta.Name
	s.Meta.Description = obj.Meta.Description
	s.Meta.Region = obj.Meta.Region
	s.Meta.Updated = obj.Meta.Updated
	s.Meta.Created = obj.Meta.Created
	s.Meta.Replicas = obj.Meta.Replicas

	s.Config.Memory = obj.Config.Memory
	s.Config.Command = strings.Join(obj.Config.Command, " ")
	s.Config.Image = obj.Config.Image

	if len(obj.Pods) == 0 {
		s.Pods = make([]v1.PodInfo, 0)
		return &s
	}

	for _, pod := range obj.Pods {
		s.Pods = append(s.Pods, v1.ToPodInfo(pod))
	}

	return &s
}

func (obj *Service) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func NewList(obj types.ServiceList) *ServiceList {
	s := ServiceList{}
	if obj == nil {
		return nil
	}
	for _, v := range obj {
		s = append(s, New(v))
	}
	return &s
}

func (s *Service) DrawTable(namespaceName string) {
	serviceTable := table.New([]string{"NAME", "DESCRIPTION", "NAMESPACE",
		"REPLICAS", "MEMORY", "IMAGE", "CREATED", "UPDATED"})
	podsTable := table.New([]string{"ID", "STATE", "STATUS", "TOTAL",
		"RUNNING", "CREATED",
		"STOPPED", "ERRORED", "CREATED POD", "UPDATED POD"})
	containersTable := table.New([]string{"ID", "IMAGE", "STATE",
		"STATUS", "CREATE", "UPDATED"})

	serviceTable.VisibleHeader = true
	podsTable.VisibleHeader = true
	containersTable.VisibleHeader = true

	serviceTable.AddRow(map[string]interface{}{
		"NAME":        s.Meta.Name,
		"DESCRIPTION": s.Meta.Description,
		"NAMESPACE":   namespaceName,
		"REPLICAS":    s.Meta.Replicas,
		"MEMORY":      s.Config.Memory,
		"IMAGE":       s.Config.Image,
		"CREATED":     s.Meta.Created.String()[:10],
		"UPDATED":     s.Meta.Updated.String()[:10],
	})
	serviceTable.Print()

	fmt.Println("\n\nPODS")
	for _, pod := range s.Pods {
		podsTable.AddRow(map[string]interface{}{
			"ID":          pod.Meta.ID,
			"STATE":       pod.Meta.State.State,
			"STATUS":      pod.Meta.State.Status,
			"TOTAL":       pod.Meta.State.Containers.Total,
			"RUNNING":     pod.Meta.State.Containers.Running,
			"CREATED":     pod.Meta.State.Containers.Created,
			"STOPPED":     pod.Meta.State.Containers.Stopped,
			"ERRORED":     pod.Meta.State.Containers.Errored,
			"CREATED POD": pod.Meta.Created.String()[:10],
			"UPDATED POD": pod.Meta.Updated.String()[:10],
		})
		podsTable.Print()

		fmt.Println("CONTAINERS")
		for _, container := range pod.Containers {
			containersTable.AddRow(map[string]interface{}{
				"ID":      container.ID[:12],
				"IMAGE":   container.Image,
				"STATE":   container.State,
				"STATUS":  container.Status,
				"CREATED": container.Created.String()[:10],
				"STARTED": container.Started.String()[:10],
			})
		}
		containersTable.Print()
	}
}

func (obj *ServiceList) ToJson() ([]byte, error) {
	if obj == nil || len(*obj) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(obj)
}

func (sl *ServiceList) DrawTable(namespaceName string) {
	t := table.New([]string{"NAME", "DESCRIPTION", "REPLICAS", "CREATED", "UPDATED"})
	t.VisibleHeader = true

	fmt.Println("NAMESPACE: ", namespaceName)
	for _, s := range *sl {
		t.AddRow(map[string]interface{}{
			"NAME":        s.Meta.Name,
			"DESCRIPTION": s.Meta.Description,
			"REPLICAS":    s.Meta.Replicas,
			"CREATED":     s.Meta.Created.String()[:10],
			"UPDATED":     s.Meta.Updated.String()[:10],
		})
	}

	t.AddRow(map[string]interface{}{})

	t.Print()
}
