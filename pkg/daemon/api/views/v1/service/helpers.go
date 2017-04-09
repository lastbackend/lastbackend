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

package service

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

func New(obj *types.Service) *Service {
	s := new(Service)
	s.User = obj.User
	s.Project = obj.Project
	s.Name = obj.Name
	s.Description = obj.Description
	s.Updated = obj.Updated
	s.Created = obj.Created

	if obj.Config != nil {
		s.Config.Region = obj.Config.Region
		s.Config.Memory = obj.Config.Memory
		s.Config.Replicas = obj.Config.Replicas
	}

	if obj.Source != nil && obj.Source.Type == types.SourceGitType {
		s.Source = new(Source)
		s.Source.Hub = obj.Source.Hub
		s.Source.Owner = obj.Source.Owner
		s.Source.Repo = obj.Source.Repo
		s.Source.Branch = obj.Source.Branch
	}

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

func (obj *ServiceList) ToJson() ([]byte, error) {
	if obj == nil || len(*obj) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(obj)
}
