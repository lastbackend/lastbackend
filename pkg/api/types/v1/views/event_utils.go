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
	"fmt"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

type EventView struct{}

func (nv *EventView) New(obj *models.Event) *Event {

	n := Event{
		Name: fmt.Sprintf("%s:%s", obj.Kind, obj.Action),
	}

	switch obj.Kind {
	case models.KindNamespace:
		n.Payload = new(NamespaceView).New(obj.Data.(*models.Namespace))
	case models.KindService:
		n.Payload = new(ServiceView).New(obj.Data.(*models.Service))
	case models.KindDeployment:
		n.Payload = new(DeploymentView).New(obj.Data.(*models.Deployment))
	case models.KindPod:
		n.Payload = new(PodView).New(obj.Data.(*models.Pod))
	case models.KindRoute:
		n.Payload = new(RouteView).New(obj.Data.(*models.Route))
	case models.KindSecret:
		n.Payload = new(SecretView).New(obj.Data.(*models.Secret))
	case models.KindConfig:
		n.Payload = new(ConfigView).New(obj.Data.(*models.Config))
	case models.KindVolume:
		n.Payload = new(VolumeView).New(obj.Data.(*models.Volume))
	case models.KindJob:
		n.Payload = new(JobView).New(obj.Data.(*models.Job))
	case models.KindTask:
		n.Payload = new(TaskView).New(obj.Data.(*models.Task))
	case models.KindNode:
		n.Payload = new(NodeView).New(obj.Data.(*models.Node))
	case models.KindDiscovery:
		n.Payload = new(DiscoveryView).New(obj.Data.(*models.Discovery))
	case models.KindIngress:
		n.Payload = new(IngressView).New(obj.Data.(*models.Ingress))
	default:

	}

	return &n
}

func (obj *Event) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
