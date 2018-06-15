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

package storage

import (
	"strings"
)

type Pod struct{}

func (Pod) Query(namespace, service, deployment, pod string) string {
	return strings.Join([]string{namespace, service, deployment, pod}, ":")
}

type Deployment struct{}

func (Deployment) Query(namespace, service, deployment string) string {
	return strings.Join([]string{namespace, service, deployment}, ":")
}

type Service struct{}

func (Service) Query(namespace string) string {
	return strings.Join([]string{namespace}, ":")
}
