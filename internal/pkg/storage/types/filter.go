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

package types

type Filter interface {
	Namespace() NamespaceFilter
	Service() ServiceFilter
	Config() ConfigFilter
	Deployment() DeploymentFilter
	Pod() PodFilter
	Endpoint() EndpointFilter
	Route() RouteFilter
	Secret() SecretFilter
	Volume() VolumeFilter
	Task() TaskFilter
	Job() JobFilter
}

type NamespaceFilter interface {
}

type ServiceFilter interface {
	ByNamespace(namespace string) string
}

type DeploymentFilter interface {
	ByNamespace(namespace string) string
	ByService(namespace, service string) string
}

type PodFilter interface {
	ByNamespace(namespace string) string
	ByService(namespace, service string) string
	ByDeployment(namespace, service, deployment string) string
	ByJob(namespace, job string) string
	ByTask(namespace, job, task string) string
}

type EndpointFilter interface {
	ByNamespace(namespace string) string
}

type RouteFilter interface {
	ByNamespace(namespace string) string
}

type SecretFilter interface {
	ByNamespace(namespace string) string
}

type ConfigFilter interface {
	ByNamespace(namespace string) string
}

type VolumeFilter interface {
	ByNamespace(namespace string) string
}

type TaskFilter interface {
	ByNamespace(namespace string) string
	ByJob(namespace, job string) string
}

type JobFilter interface {
	ByNamespace(namespace string) string
}
