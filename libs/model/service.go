package model

import (
	"encoding/json"
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
	Detail *service.Service `json:"detail,omitempty" gorethink:"-"`
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

func (s *Service) DrawTable() {
	t := table.New([]string{"ID", "Project", "Name", "Created", "Updated"})
	t.AddRow(map[string]interface{}{
		"ID":      s.ID,
		"Project": s.Image,
		"Name":    s.Name,
		"Created": s.Created.String()[:10],
		"Updated": s.Updated.String()[:10],
	})
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

func (s *ServiceList) DrawTable() {
	table := table.New([]string{"ID", "Project", "Name", "Created", "Updated"})

	for _, service := range *s {
		table.AddRow(map[string]interface{}{
			"ID":      service.ID,
			"Project": service.Image,
			"Name":    service.Name,
			"Created": service.Created.String()[:10],
			"Updated": service.Updated.String()[:10],
		})
	}

	table.Print()
}

type ServiceConfig struct {
	Replicas   int32           `json:"scale,omitempty" yaml:"scale"`
	Containers ContainerConfig `json:"containers,omitempty" yaml:"containers"`
	Command    []string        `json:"command,omitempty" yaml:"command"`
	Args       []string        `json:"args,omitempty" yaml:"args"`
	WorkingDir string          `json:"workdir,omitempty" yaml:"workdir"`
	Ports      []Port          `json:"ports,omitempty" yaml:"ports"`
	Env        []EnvVar        `json:"env,omitempty" yaml:"env"`
}

type Port struct {
	Name          string `json:"name,omitempty" yaml:"name"`
	HostPort      int32  `json:"host,omitempty" yaml:"host"`
	ContainerPort int32  `json:"container" yaml:"container"`
	Protocol      string `json:"protocol,omitempty" yaml:"protocol"`
	HostIP        string `json:"ip,omitempty" yaml:"ip"`
}

type EnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value,omitempty" yaml:"value"`
}

type ContainerConfig struct {
}
