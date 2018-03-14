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

package v1

const logLevel = 5

type IView interface {
	Cluster() *ClusterView
	Node() *NodeView

	Namespace() *NamespaceView
	Route() *RouteView
	Service() *ServiceView
	Deployment() *DeploymentView
	Pod() *PodView
	Container() *ContainerView
}

type View struct{}

func (View) Cluster() *ClusterView {
	return new(ClusterView)
}
func (View) Node() *NodeView {
	return new(NodeView)
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
func (View) Deployment() *DeploymentView {
	return new(DeploymentView)
}
func (View) Pod() *PodView {
	return new(PodView)
}
func (View) Container() *ContainerView {
	return new(ContainerView)
}
