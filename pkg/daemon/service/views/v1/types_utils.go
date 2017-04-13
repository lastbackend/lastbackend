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
	"github.com/lastbackend/lastbackend/pkg/util/table"
)

func New(obj *types.Service) *Service {
	s := new(Service)

	s.Name = obj.Meta.Name
	s.Description = obj.Meta.Description
	s.Region = obj.Meta.Region
	s.Updated = obj.Meta.Updated
	s.Created = obj.Meta.Created

	s.Config.Memory = obj.Config.Memory
	s.Config.Replicas = obj.Config.Replicas
	s.Config.Image = obj.Config.Image

	return s
}

func (obj *Service) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func NewList(obj *types.ServiceList) *ServiceList {
	s := new(ServiceList)
	if obj == nil {
		return nil
	}
	for _, v := range *obj {
		*s = append(*s, *New(&v))
	}
	return s
}

func (s *Service) DrawTable(projectName string) {
	t := table.New([]string{"Name", "Description", "Created", "Updated"})
	t.VisibleHeader = true

	t.AddRow(map[string]interface{}{
		"Name":        s.Name,
		"Description": s.Description,
		"Created":     s.Created.String()[:10],
		"Updated":     s.Updated.String()[:10],
	})

	t.AddRow(map[string]interface{}{})

	t.Print()
}

func (obj *ServiceList) ToJson() ([]byte, error) {
	if obj == nil || len(*obj) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(obj)
}

func (s *ServiceList) DrawTable(projectName string) {
	t := table.New([]string{"Name", "Description", "Created", "Updated"})
	t.VisibleHeader = true

	for _, ss := range *s {
		t.AddRow(map[string]interface{}{
			"Name":        ss.Name,
			"Description": ss.Description,
			"Created":     ss.Created.String()[:10],
			"Updated":     ss.Updated.String()[:10],
		})
	}

	t.AddRow(map[string]interface{}{})

	t.Print()
}
