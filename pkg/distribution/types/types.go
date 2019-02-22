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

package types

const logLevel = 3

const (
	DEFAULT_NAMESPACE = "default"
	SYSTEM_NAMESPACE  = "system"

	DEFAULT_RESOURCE_LIMITS_RAM = "128mib"
	DEFAULT_RESOURCE_LIMITS_CPU = "0.1"

	DEFAULT_MEMORY_MIN        = 128
	DEFAULT_REPLICAS_MIN      = 1
	DEFAULT_DESCRIPTION_LIMIT = 512

	KindSecret     = "secret"
	KindRoute      = "route"
	KindNamespace  = "namespace"
	KindService    = "service"
	KindDeployment = "deployment"
	KindJob        = "job"
	KindTask       = "task"
	KindPod        = "pod"
	KindEndpoint   = "endpoint"
	KindConfig     = "config"
	KindVolume     = "volume"
)

type Vault struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Token    string `yaml:"token"`
}
