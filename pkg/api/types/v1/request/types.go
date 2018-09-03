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

package request

const (
	DEFAULT_MEMORY_MIN        = 128
	DEFAULT_REPLICAS_MIN      = 1
	DEFAULT_DESCRIPTION_LIMIT = 512
)

const (
	StateProvision   = "provision"
	StateInitialized = "initialized"
	StateWarning     = "warning"
	StateReady       = "ready"
)

const (
	StatePull    = "pull"
	StateDestroy = "destroy"
	StateCancel  = "cancel"
)

const (
	StateCreated   = "created"
	StateStarting  = "starting"
	StateStarted   = "started"
	StateStopped   = "stopped"
	StateDestroyed = "destroyed"
)

const (
	StateExited  = "exited"
	StateRunning = "running"
	StateError   = "error"
)

const (
	StepInitialized = "initialized"
	StepPull        = "pull"
	StepReady       = "ready"
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
	Trigger() *TriggerRequest
	Volume() *VolumeRequest
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
func (Request) Trigger() *TriggerRequest {
	return new(TriggerRequest)
}
func (Request) Volume() *VolumeRequest {
	return new(VolumeRequest)
}
