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

package request

const (
	DEFAULT_DESCRIPTION_LIMIT = 512
)

const (
	StateProvision = "provision"
	StateReady     = "ready"
	StateDestroy   = "destroy"
	StateCreated   = "created"
	StateDestroyed = "destroyed"
	StateError     = "error"
)

type IRequest interface {
	Cluster() *ClusterRequest
	Deployment() *DeploymentRequest
	Namespace() *NamespaceRequest
	Node() *NodeRequest
	Endpoint() *EndpointRequest
	Route() *RouteRequest
	Service() *ServiceRequest
	Secret() *SecretRequest
	Config() *ConfigRequest
	Volume() *VolumeRequest
	Ingress() *IngressRequest
	Discovery() *DiscoveryRequest
	Job() *JobRequest
	Task() *TaskRequest
}

type Request struct{}

func (Request) Cluster() *ClusterRequest {
	return new(ClusterRequest)
}
func (Request) Deployment() *DeploymentRequest {
	return new(DeploymentRequest)
}
func (Request) Namespace() *NamespaceRequest {
	return new(NamespaceRequest)
}
func (Request) Node() *NodeRequest {
	return new(NodeRequest)
}
func (Request) Endpoint() *EndpointRequest {
	return new(EndpointRequest)
}
func (Request) Route() *RouteRequest {
	return new(RouteRequest)
}
func (Request) Service() *ServiceRequest {
	return new(ServiceRequest)
}
func (Request) Secret() *SecretRequest {
	return new(SecretRequest)
}
func (Request) Config() *ConfigRequest {
	return new(ConfigRequest)
}
func (Request) Volume() *VolumeRequest {
	return new(VolumeRequest)
}
func (Request) Ingress() *IngressRequest {
	return new(IngressRequest)
}
func (Request) Discovery() *DiscoveryRequest {
	return new(DiscoveryRequest)
}

func (Request) Job() *JobRequest {
	return new(JobRequest)
}

func (Request) Task() *TaskRequest {
	return new(TaskRequest)
}
