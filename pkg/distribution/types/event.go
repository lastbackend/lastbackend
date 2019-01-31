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

package types

const (
	EventActionCreate = "create"
	EventActionUpdate = "update"
	EventActionDelete = "delete"
	EventActionError  = "error"
)

type event struct {
	Action   string
	Name     string
	SelfLink string
}

type Event struct {
	event
	Kind string
	Data interface{}
}

type NamespaceEvent struct {
	event
	Data *Namespace
}

type ClusterEvent struct {
	event
	Data *Cluster
}

type ServiceEvent struct {
	event
	Data *Service
}

type VolumeEvent struct {
	event
	Data *Volume
}

type NetworkEvent struct {
	event
	Data *Network
}

type SubnetEvent struct {
	event
	Data *Subnet
}

type SecretEvent struct {
	event
	Data *Secret
}

type ConfigEvent struct {
	event
	Data *Config
}

type RouteEvent struct {
	event
	Data *Route
}

type IngressEvent struct {
	event
	Data *Ingress
}

type IngressStatusEvent struct {
	event
	Data *IngressStatus
}

type DiscoveryEvent struct {
	event
	Data *Discovery
}

type DiscoveryStatusEvent struct {
	event
	Data *DiscoveryStatus
}

type EndpointEvent struct {
	event
	Data *Endpoint
}

type DeploymentEvent struct {
	event
	Data *Deployment
}

type PodEvent struct {
	event
	Data *Pod
}

type PodManifestEvent struct {
	event
	Node string
	Data *PodManifest
}

type VolumeManifestEvent struct {
	event
	Node string
	Data *VolumeManifest
}

type EndpointManifestEvent struct {
	event
	Data *EndpointManifest
}

type SubnetManifestEvent struct {
	event
	Data *SubnetManifest
}

type SecretManifestEvent struct {
	event
	Data *SecretManifest
}

type RouteManifestEvent struct {
	event
	Ingress string
	Data    *RouteManifest
}

type NodeEvent struct {
	event
	Data *Node
}

type NodeStatusEvent struct {
	event
	Data *NodeStatus
}

type JobEvent struct {
	event
	Data *Task
}

type JobRunnerEvent struct {
	event
	Data *Job
}

func (e *event) IsActionCreate() bool {
	return e.Action == EventActionCreate
}

func (e *event) IsActionUpdate() bool {
	return e.Action == EventActionUpdate
}

func (e *event) IsActionRemove() bool {
	return e.Action == EventActionDelete
}

func (e *event) IsActionError() bool {
	return e.Action == EventActionError
}
