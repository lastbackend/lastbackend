//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type EventView struct{}

func (nv *EventView) New(obj *types.Event) *Event {

	n := Event{
		Name: fmt.Sprintf("%s:%s", obj.Kind, obj.Action),
	}

	switch obj.Kind {
	case types.KindNamespace:
		n.Payload = new(NamespaceView).New(obj.Data.(*types.Namespace))
	case types.KindService:
		n.Payload = new(ServiceView).New(obj.Data.(*types.Service))
	case types.KindDeployment:
		n.Payload = new(DeploymentView).New(obj.Data.(*types.Deployment), nil)
	case types.KindPod:
		n.Payload = new(PodView).New(obj.Data.(*types.Pod))
	case types.KindRoute:
		n.Payload = new(RouteView).New(obj.Data.(*types.Route))
	case types.KindSecret:
		n.Payload = new(SecretView).New(obj.Data.(*types.Secret))
	case types.KindConfig:
		n.Payload = new(ConfigView).New(obj.Data.(*types.Config))
	case types.KindVolume:
		n.Payload = new(VolumeView).New(obj.Data.(*types.Volume))
	case types.KindNode:
		n.Payload = new(NodeView).New(obj.Data.(*types.Node))
	case types.KindDiscovery:
		n.Payload = new(DiscoveryView).New(obj.Data.(*types.Discovery))
	case types.KindIngress:
		n.Payload = new(IngressView).New(obj.Data.(*types.Ingress))
	default:

	}

	return &n
}

func (obj *Event) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
