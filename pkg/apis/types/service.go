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

package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/util/table"
)

type ServiceList []Service

type Service struct {
	serviceMeta

	// Service username
	User string `json:"user"`
	// Service project
	Project string `json:"project"`

	// Service custom domains
	Domains []string `json:"domains"`
	// Service source info
	Source *ServiceSource `json:"source,omitempty"`
	// Service config info
	Config *ServiceConfig `json:"config,omitempty"`
}

const (
	SourceGitType      = "git"
	SourceDockerType   = "docker"
	SourceTemplateType = "template"
)

type serviceMeta struct{ ServiceMeta }
type ServiceMeta struct {
	meta

	// Add fields to expand the meta data
	// Example:
	// Note string `json:"note,omitempty"`
	// Uptime time.Time `json:"uptime"

	// Service image
	Image string `json:"image"`
}

type ServiceSource struct {
	Hub    string `json:"hub"`
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
}

type ServiceConfig struct {
	Replicas   int      `json:"scale"`
	Memory     int      `json:"memory"`
	Region     string   `json:"region"`
	WorkingDir string   `json:"workdir"`
	Entrypoint string   `json:"entrypoint"`
	Image      string   `json:"image"`
	Command    []string `json:"command"`
	Args       []string `json:"args"`
	EnvVars    []string `json:"env"`
	Ports      []Port   `json:"ports"`
}

func (c *ServiceConfig) Update(patch *ServiceConfig) error {

	if patch.Replicas < 0 {
		return errors.New("The value of the `scale` parameter must be at least 1")
	}
	c.Replicas = patch.Replicas

	if patch.Memory < 32 {
		return errors.New("The value of the `memory` parameter must be at least 32")
	}
	c.Memory = patch.Memory

	c.WorkingDir = patch.WorkingDir
	c.Entrypoint = patch.Entrypoint
	c.Image = patch.Image
	c.Command = patch.Command
	c.Args = patch.Args

	c.Ports = patch.Ports

	// TODO: Check valid format env params
	c.EnvVars = patch.EnvVars

	return nil
}

func (ServiceConfig) GetDefault() *ServiceConfig {
	var config = new(ServiceConfig)
	config.Memory = 256
	return config
}

type Port struct {
	Name      string `json:"name"`
	Protocol  string `json:"protocol"`
	Container int32  `json:"container"`
	Published bool   `json:"published"`
}

func (s *Service) ToJson() ([]byte, error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *Service) DrawTable(projectName string) {
	//table.PrintHorizontal(map[string]interface{}{
	//	"ID":      s.ID,
	//	"NAME":    s.Name,
	//	"PROJECT": projectName,
	//	"PODS":    len(s.Spec.PodList),
	//})
	//
	t := table.New([]string{" ", "NAME", "STATUS", "CONTAINERS"})
	//t.VisibleHeader = true
	//
	//for _, pod := range s.Spec.PodList {
	//	t.AddRow(map[string]interface{}{
	//		" ":          "",
	//		"NAME":       pod.Name,
	//		"STATUS":     pod.Status,
	//		"CONTAINERS": len(pod.ContainerList),
	//	})
	//}
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

	//for _, s := range *s {
	//
	//	t := make(map[string]interface{})
	//	t["ID"] = s.ID
	//	t["NAME"] = s.Name
	//
	//	if s.Spec != nil {
	//		t["PODS"] = len(s.Spec.PodList)
	//	}
	//
	//	table.PrintHorizontal(t)
	//
	//	if s.Spec != nil {
	//		for _, pod := range s.Spec.PodList {
	//			tpods := table.New([]string{" ", "NAME", "STATUS", "CONTAINERS"})
	//			tpods.VisibleHeader = true
	//
	//			tpods.AddRow(map[string]interface{}{
	//				" ":          "",
	//				"NAME":       pod.Name,
	//				"STATUS":     pod.Status,
	//				"CONTAINERS": len(pod.ContainerList),
	//			})
	//			tpods.Print()
	//		}
	//	}
	//
	//	fmt.Print("\n\n")
	//}
}

type ServiceUpdateConfig struct {
	Name        *string            `json:"name,omitempty" yaml:"name,omitempty"`
	Description *string            `json:"description,omitempty" yaml:"description,omitempty"`
	Replicas    *int32             `json:"scale,omitempty" yaml:"scale,omitempty"`
	Containers  *[]ContainerConfig `json:"containers,omitempty" yaml:"containers,omitempty"`
}

type ContainerConfig struct {
	Image      string   `json:"image" yaml:"image"`
	Name       string   `json:"name" yaml:"name"`
	WorkingDir string   `json:"workdir" yaml:"workdir"`
	Command    []string `json:"command" yaml:"command"`
	Args       []string `json:"args" yaml:"args"`
	EnvVars    []string `json:"env" yaml:"env"`
	Ports      []Port   `json:"ports" yaml:"ports"`
}
