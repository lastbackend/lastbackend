package model

import (
	"encoding/json"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/service"
	"github.com/lastbackend/lastbackend/pkg/util/table"
	"time"
	"fmt"
)

type ServiceList []Service

type Service struct {
	// Service uuid, incremented automatically
	ID string `json:"id" gorethink:"id,omitempty"`
	// Service user
	User string `json:"user" gorethink:"user,omitempty"`
	// Service project
	Project string `json:"project" gorethink:"project,omitempty"`
	// Service image
	Image string `json:"image" gorethink:"image,omitempty"`
	// Service name
	Name string `json:"name" gorethink:"name,omitempty"`
	// Service spec
	Spec *service.Service `json:"spec,omitempty" gorethink:"-"`
	// Service created time
	Created time.Time `json:"created" gorethink:"created,omitempty"`
	// Service updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}

func (s *Service) ToJson() ([]byte, *e.Err) {
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, e.New("service").Unknown(err)
	}

	return buf, nil
}

func (s *Service) DrawTable(projectName string) {
	table.PrintHorizontal(map[string]interface{}{
		"ID":      s.ID,
		"NAME":    s.Name,
		"PROJECT": projectName,
		"PODS":    s.Spec.PodList.ListMeta.Total,
	})

	t := table.New([]string{" ", "NAME", "STATUS", "RESTARTS", "CONTAINERS"})
	t.VisibleHeader = true

	for _, pod := range s.Spec.PodList.Pods {
		t.AddRow(map[string]interface{}{
			" ":          "",
			"NAME":       pod.ObjectMeta.Name,
			"STATUS":     pod.PodStatus.PodPhase,
			"RESTARTS":   pod.RestartCount,
			"CONTAINERS": pod.Containers.ListMeta.Total,
		})
	}
	t.AddRow(map[string]interface{}{})

	t.Print()
}

func (s *ServiceList) ToJson() ([]byte, *e.Err) {

	if s == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(s)
	if err != nil {
		return nil, e.New("service").Unknown(err)
	}

	return buf, nil
}

func (s *ServiceList) DrawTable(projectName string) {
	fmt.Print(" Project ", projectName + "\n\n")

	for _, s := range *s {
		//tservice :=  table.New([]string{"ID", "NAME", "PODS"})
		//tservice.VisibleHeader = true
		//
		//tservice.AddRow(map[string]interface{}{
		//	"ID": s.ID,
		//	"NAME": s.Name,
		//	"PODS": s.Spec.PodList.ListMeta.Total,
		//})
		//tservice.Print()

		table.PrintHorizontal(map[string]interface{}{
			"ID":      s.ID,
			"NAME":    s.Name,
			"PODS":    s.Spec.PodList.ListMeta.Total,
		})

		for _, pod := range s.Spec.PodList.Pods {
			tpods := table.New([]string{" ", "NAME", "STATUS", "RESTARTS", "CONTAINERS"})
			tpods.VisibleHeader = true

			tpods.AddRow(map[string]interface{}{
				" ":          "",
				"NAME":       pod.ObjectMeta.Name,
				"STATUS":     pod.PodStatus.PodPhase,
				"RESTARTS":   pod.RestartCount,
				"CONTAINERS": pod.Containers.ListMeta.Total,
			})
			tpods.Print()
		}

		fmt.Print("\n\n")
	}
}
