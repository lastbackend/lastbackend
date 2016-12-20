package model

import (
	"encoding/json"
	"fmt"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/service"
	"github.com/lastbackend/lastbackend/pkg/util/table"
	"time"
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

	t := table.New([]string{" ", "NAME", "STATUS", "RESTARTS"})
	t.VisibleHeader = true
	for _, pod := range s.Spec.PodList.Pods {
		t.AddRow(map[string]interface{}{
			" ":        "",
			"NAME":     pod.ObjectMeta.Name,
			"STATUS":   pod.PodStatus.PodPhase,
			"RESTARTS": pod.RestartCount,
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
	fmt.Println("Project: " + projectName)

	t := table.New([]string{"", "ID", "Name", "Created", "Updated"})
	t.VisibleHeader = true

	for _, s := range *s {
		t.AddRow(map[string]interface{}{
			"":        "",
			"ID":      s.ID,
			"Name":    s.Name,
			"Created": s.Created.String()[:10],
			"Updated": s.Updated.String()[:10],
		})
	}

	t.Print()
}
