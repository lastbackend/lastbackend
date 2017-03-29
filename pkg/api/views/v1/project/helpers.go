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

package project

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/api/types"
)

func New(obj *types.Project) *Project {
	p := new(Project)

	p.User = obj.User
	p.Name = obj.Name
	p.Description = obj.Description
	p.Updated = obj.Updated
	p.Created = obj.Created

	return p
}

func (obj *Project) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func NewList(obj *types.ProjectList) *ProjectList {
	p := new(ProjectList)
	if obj == nil {
		return nil
	}
	for _, v := range *obj {
		*p = append(*p, *New(&v))
	}
	return p
}

func (obj *ProjectList) ToJson() ([]byte, error) {
	if obj == nil || len(*obj) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(obj)
}
