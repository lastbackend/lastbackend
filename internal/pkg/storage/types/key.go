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

type Key interface {
	Namespace(name string) string
	Service(namespace, name string) string
	Deployment(namespace, service, name string) string
	Pod(namespace, service, deployment, name string) string
	Endpoint(namespace, service string) string
	Config(namespace, name string) string
	Secret(namespace, name string) string
	Volume(namespace, name string) string
	Ingress(name string) string
	Discovery(name string) string
	Process(kind, hostname string, pid int, lead bool) string
	Manifest(name string) string
	Node(name string) string
	Route(namespace, name string) string
	Subnet(name string) string
}
