//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package views

import (
	"encoding/json"

	"github.com/lastbackend/lastbackend/internal/pkg/types"
)

type APIView struct{}

func (nv *APIView) New(obj *types.API) *API {
	n := API{}
	n.Meta = nv.ToAPIMeta(obj.Meta)
	n.Status = nv.ToAPIStatus(obj.Status)
	return &n
}

func (nv *APIView) ToAPIMeta(meta types.APIMeta) APIMeta {
	m := APIMeta{}
	m.Name = meta.Name
	m.Description = meta.Description
	m.Created = meta.Created
	m.Updated = meta.Updated
	return m
}

func (nv *APIView) ToAPIStatus(status types.APIStatus) APIStatus {
	return APIStatus{
		Ready: status.Ready,
	}
}

func (obj *API) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (nv *APIView) NewList(obj *types.APIList) *APIList {
	if obj == nil {
		return nil
	}
	ingresses := make(APIList, 0)
	for _, v := range obj.Items {
		nn := nv.New(v)
		ingresses[nn.Meta.Name] = nn
	}

	return &ingresses
}

func (obj *APIList) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
