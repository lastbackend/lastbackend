package model

import (
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/service"
	"github.com/lastbackend/lastbackend/pkg/service/resource/container"
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
	// Service description
	Description string `json:"description" gorethink:"description,omitempty"`
	// Service spec
	Spec *service.Service `json:"spec,omitempty" gorethink:"-"`
	// Service created time
	Created time.Time `json:"created" gorethink:"created,omitempty"`
	// Service updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}

func (s *Service) ToJson() ([]byte, error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *Service) GetConfig() *ServiceUpdateConfig {

	var config = new(ServiceUpdateConfig)

	config.Description = s.Description
	config.Replicas = s.Spec.Scale
	config.Containers = make([]ContainerConfig, len(s.Spec.Template.ContainerList))

	for index, val := range s.Spec.Template.ContainerList {
		cfg := ContainerConfig{}

		cfg.Name = val.Name
		cfg.Image = val.Image
		cfg.WorkingDir = val.WorkingDir
		cfg.Command = val.Command
		cfg.Args = val.Args

		for _, val := range val.PortList {
			cfg.Ports = append(cfg.Ports, Port{
				Name:      val.Name,
				Container: val.Container,
				Protocol:  string(val.Protocol),
			})
		}

		for _, val := range val.EnvList {
			cfg.Env = append(cfg.Env, EnvVar{
				Name:  val.Name,
				Value: val.Value,
			})
		}

		config.Containers[index] = cfg
	}

	return config
}

func (s *Service) DrawTable(projectName string) {
	table.PrintHorizontal(map[string]interface{}{
		"ID":      s.ID,
		"NAME":    s.Name,
		"PROJECT": projectName,
		"PODS":    len(s.Spec.PodList),
	})

	t := table.New([]string{" ", "NAME", "STATUS", "CONTAINERS"})
	t.VisibleHeader = true

	for _, pod := range s.Spec.PodList {
		t.AddRow(map[string]interface{}{
			" ":          "",
			"NAME":       pod.Name,
			"STATUS":     pod.Status,
			"CONTAINERS": len(pod.ContainerList),
		})
	}
	t.AddRow(map[string]interface{}{})

	t.Print()
}

func (s *ServiceList) ToJson() ([]byte, error) {

	if s == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *ServiceList) DrawTable(projectName string) {
	fmt.Print(" Project ", projectName+"\n\n")

	for _, s := range *s {

		t := make(map[string]interface{})
		t["ID"] = s.ID
		t["NAME"] = s.Name

		if s.Spec != nil {
			t["PODS"] = len(s.Spec.PodList)
		}

		table.PrintHorizontal(t)

		if s.Spec != nil {
			for _, pod := range s.Spec.PodList {
				tpods := table.New([]string{" ", "NAME", "STATUS", "CONTAINERS"})
				tpods.VisibleHeader = true

				tpods.AddRow(map[string]interface{}{
					" ":          "",
					"NAME":       pod.Name,
					"STATUS":     pod.Status,
					"CONTAINERS": len(pod.ContainerList),
				})
				tpods.Print()
			}
		}

		fmt.Print("\n\n")
	}
}

type ServiceUpdateConfig struct {
	Description string            `json:"description" yaml:"description"`
	Replicas    int32             `json:"scale" yaml:"scale"`
	Containers  []ContainerConfig `json:"containers" yaml:"containers"`
}

type ContainerConfig struct {
	Image      string   `json:"image" yaml:"image"`
	Name       string   `json:"name" yaml:"name"`
	WorkingDir string   `json:"workdir" yaml:"workdir"`
	Command    []string `json:"command" yaml:"command"`
	Args       []string `json:"args" yaml:"args"`
	Env        []EnvVar `json:"env" yaml:"env"`
	Ports      []Port   `json:"ports" yaml:"ports"`
}

type Port struct {
	Name      string `json:"name" yaml:"name"`
	Container int32  `json:"container" yaml:"container"`
	Protocol  string `json:"protocol" yaml:"protocol"`
}

type EnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

func (s ServiceUpdateConfig) CreateServiceConfig() *service.ServiceConfig {
	var cfg = new(service.ServiceConfig)

	cfg.Replicas = s.Replicas

	for _, val := range s.Containers {
		c := container.Container{}

		c.Name = val.Name
		c.Image = val.Image
		c.WorkingDir = val.WorkingDir
		c.Command = val.Command
		c.Args = val.Args

		for _, item := range val.Ports {
			c.Ports = append(c.Ports, container.Port{
				Name:          item.Name,
				ContainerPort: item.Container,
				Protocol:      item.Protocol,
			})

			for _, val := range val.Env {
				c.Env = append(c.Env, container.EnvVar{
					Name:  val.Name,
					Value: val.Value,
				})
			}

		}

		cfg.Containers = append(cfg.Containers, c)
	}

	return cfg
}
