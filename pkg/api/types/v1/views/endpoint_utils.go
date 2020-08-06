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
	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

type EndpointView struct{}

func (ev *EndpointView) New(obj *models.Endpoint) *Endpoint {
	e := Endpoint{}
	e.Meta = ev.ToEndpointMeta(obj.Meta)
	e.Status = ev.ToEndpointStatus(obj.Status)
	e.Spec = ev.ToEndpointSpec(obj.Spec)
	return &e
}

func (ev *EndpointView) ToEndpointMeta(meta models.EndpointMeta) EndpointMeta {
	return EndpointMeta{}
}

func (ev *EndpointView) ToEndpointStatus(meta models.EndpointStatus) EndpointStatus {
	return EndpointStatus{}
}

func (ev *EndpointView) ToEndpointSpec(meta models.EndpointSpec) EndpointSpec {
	return EndpointSpec{}
}

func (ev *EndpointView) NewList(obj map[string]*models.Endpoint) *EndpointList {
	if obj == nil {
		return nil
	}
	el := make(EndpointList, 0)
	for _, v := range obj {
		nn := ev.New(v)
		el[nn.Meta.SelfLink] = nn
	}

	return &el
}

func (obj *Endpoint) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *EndpointList) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
