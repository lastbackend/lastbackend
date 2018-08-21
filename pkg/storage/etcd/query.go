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
	"strings"
)

const emptyString = ""

func BuildDeploymentQuery(namespace, service string) string {
	return strings.Join([]string{namespace, service, emptyString}, ":")
}

func BuildSecretQuery(namespace, service string) string {
	return strings.Join([]string{namespace, service, emptyString}, ":")
}

func BuildServiceQuery(namespace string) string {
	return strings.Join([]string{namespace, emptyString}, ":")
}

func BuildEndpointQuery(namespace string) string {
	return strings.Join([]string{namespace, emptyString}, ":")
}

func BuildPodQuery(namespace, service, deployment string) string {
	return strings.Join([]string{namespace, service, deployment, emptyString}, ":")
}

func BuildRouteQuery(namespace, route string) string {
	return strings.Join([]string{namespace, route, emptyString}, ":")
}

func BuildTriggerQuery(namespace, service, name string) string {
	return strings.Join([]string{namespace, service, name, emptyString}, ":")
}

func BuildVolumeQuery(namespace string) string {
	return strings.Join([]string{namespace, emptyString}, ":")
}
