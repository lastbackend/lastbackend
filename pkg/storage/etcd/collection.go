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

package etcd

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/storage/types"
)

const (
	namespaceCollection  = "namespace"
	secretCollection     = "secret"
	endpointCollection   = "endpoint"
	serviceCollection    = "service"
	deploymentCollection = "deployment"
	podCollection        = "pod"
	volumeCollection     = "volume"

	manifestCollection = "manifest"

	clusterCollection = "cluster"
	nodeCollection    = "node"
	networkCollection = "network"
	subnetCollection  = "subnet"

	discoveryCollection = "discovery"
	ingressCollection = "ingress"
	routeCollection   = "route"

	systemCollection  = "system"
	triggerCollection = "trigger"
	testCollection    = "test"
)

type Collection struct{}

type ManifestCollection struct{}

func (Collection) Namespace() string {
	return namespaceCollection
}

func (Collection) Secret() string {
	return secretCollection
}

func (Collection) Endpoint() string {
	return endpointCollection
}

func (Collection) Service() string {
	return serviceCollection
}

func (Collection) Deployment() string {
	return deploymentCollection
}

func (Collection) Pod() string {
	return podCollection
}

func (Collection) Volume() string {
	return volumeCollection
}

func (Collection) Discovery() string {
	return discoveryCollection
}

func (Collection) Ingress() string {
	return ingressCollection
}

func (Collection) Route() string {
	return routeCollection
}

func (Collection) System() string {
	return systemCollection
}

func (Collection) Trigger() string {
	return triggerCollection
}

func (Collection) Cluster() string {
	return clusterCollection
}

func (Collection) Node() string {
	return nodeCollection
}

func (Collection) Network() string {
	return networkCollection
}

func (Collection) Subnet() string {
	return subnetCollection
}

func (Collection) Manifest() types.ManifestCollection {
	return new(ManifestCollection)
}

func (Collection) Test() string {
	return testCollection
}

func (ManifestCollection) Node() string {
	return fmt.Sprintf("%s/%s", manifestCollection, nodeCollection)
}

func (ManifestCollection) Cluster() string {
	return fmt.Sprintf("%s/%s", manifestCollection, clusterCollection)
}

func (ManifestCollection) Pod(node string) string {
	return fmt.Sprintf("%s/%s/%s/%s", manifestCollection, nodeCollection, node, podCollection)
}

func (ManifestCollection) Volume(node string) string {
	return fmt.Sprintf("%s/%s/%s/%s", manifestCollection, nodeCollection, node, volumeCollection)
}

func (ManifestCollection) Ingress() string {
	return fmt.Sprintf("%s/%s", manifestCollection, ingressCollection)
}

func (ManifestCollection) Subnet() string {
	return fmt.Sprintf("%s/%s/%s", manifestCollection, clusterCollection, subnetCollection)
}

func (ManifestCollection) Endpoint() string {
	return fmt.Sprintf("%s/%s/%s", manifestCollection, clusterCollection, endpointCollection)
}

func (ManifestCollection) Secret() string {
	return fmt.Sprintf("%s/%s/%s", manifestCollection, clusterCollection, secretCollection)
}
