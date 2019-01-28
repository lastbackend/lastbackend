//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

const logLevel = 5

type IView interface {
	Cluster() *ClusterView
	Node() *NodeView
	Ingress() *IngressView
	Discovery() *DiscoveryView

	Namespace() *NamespaceView
	Route() *RouteView
	Service() *ServiceView
	Secret() *SecretView
	Config() *ConfigView
	Deployment() *DeploymentView
	Endpoint() *EndpointView
	Pod() *Pod
	Container() *ContainerView
	Volume() *VolumeView

	Event() *EventView
}

type View struct{}

func (View) Cluster() *ClusterView {
	return new(ClusterView)
}
func (View) Node() *NodeView {
	return new(NodeView)
}
func (View) Ingress() *IngressView {
	return new(IngressView)
}
func (View) Discovery() *DiscoveryView {
	return new(DiscoveryView)
}

func (View) Namespace() *NamespaceView {
	return new(NamespaceView)
}
func (View) Route() *RouteView {
	return new(RouteView)
}
func (View) Service() *ServiceView {
	return new(ServiceView)
}
func (View) Secret() *SecretView {
	return new(SecretView)
}
func (View) Config() *ConfigView {
	return new(ConfigView)
}
func (View) Deployment() *DeploymentView {
	return new(DeploymentView)
}
func (View) Pod() *Pod {
	return new(Pod)
}
func (View) Endpoint() *EndpointView {
	return new(EndpointView)
}
func (View) Container() *ContainerView {
	return new(ContainerView)
}

func (View) Volume() *VolumeView {
	return new(VolumeView)
}

func (View) Event() *EventView {
	return new(EventView)
}
