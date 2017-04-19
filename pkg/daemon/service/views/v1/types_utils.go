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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/pod/views/v1"
	"strings"
)

func New(obj *types.Service) *Service {
	s := Service{}

	s.Meta.Name = obj.Meta.Name
	s.Meta.Description = obj.Meta.Description
	s.Meta.Namespace = obj.Meta.Namespace
	s.Meta.Region = obj.Meta.Region
	s.Meta.Updated = obj.Meta.Updated
	s.Meta.Created = obj.Meta.Created
	s.Meta.Replicas = obj.Meta.Replicas

	if len(obj.Spec) == 0 {
		s.Spec = make([]SpecInfo, 0)
		return &s
	}

	for _, spec := range obj.Spec {
		s.Spec = append(s.Spec, ToSpecInfo(spec))
	}

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

func (s *Service) DrawTable(projectName string) {
	//table.PrintHorizontal(map[string]interface{}{
	//	"ID":      s.ID,
	//	"NAME":    s.Name,
	//	"PROJECT": projectName,
	//	"PODS":    len(s.Spec.PodList),
	//})
	//
	//t := table.New([]string{" ", "NAME", "STATUS", "CONTAINERS"})
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
	//t.AddRow(map[string]interface{}{})
	//
	//t.Print()
}

func (obj *ServiceList) ToJson() ([]byte, error) {
	if obj == nil || len(*obj) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(obj)
}

func (s *ServiceList) DrawTable(projectName string) {
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

func ToSpecInfo(spec *types.ServiceSpec) SpecInfo {
	info := SpecInfo{
		Meta:    ToSpecMeta(spec.Meta),
		Memory:  spec.Memory,
		Command: strings.Join(spec.Command, " "),
		Image:   spec.Image,
		EnvVars: spec.EnvVars,
	}

	info.EnvVars = make([]string, len(spec.EnvVars))
	info.EnvVars = append(info.EnvVars, spec.EnvVars...)

	info.Ports = make([]Port, len(spec.Ports))
	for _, port := range spec.Ports {
		info.Ports = append(info.Ports, Port{
			External:  port.Host,
			Internal:  port.Container,
			Published: port.Published,
			Protocol:  port.Protocol,
		})
	}

	return info
}

func ToSpecMeta(meta types.SpecMeta) SpecMeta {
	m := SpecMeta{
		ID:      meta.ID,
		Labels:  meta.Labels,
		Created: meta.Created,
		Updated: meta.Updated,
	}

	if len(m.Labels) == 0 {
		m.Labels = make(map[string]string)
	}

	return m
}
