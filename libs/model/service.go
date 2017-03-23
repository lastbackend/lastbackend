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

package model

import (
	"encoding/json"
	"fmt"
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
	Spec interface{} `json:"spec,omitempty" gorethink:"-"`
	// Service created time
	Source *Source `json:"source,omitempty" gorethink:"source,omitempty"`
	// Service created time
	Created time.Time `json:"created" gorethink:"created,omitempty"`
	// Service updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}

type Source struct {
	Hub      string `json:"hub" gorethink:"hub,omitempty"`
	Username string `json:"username" gorethink:"username,omitempty"`
	Repo     string `json:"repo" gorethink:"repo,omitempty"`
	Branch   string `json:"branch" gorethink:"branch,omitempty"`
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