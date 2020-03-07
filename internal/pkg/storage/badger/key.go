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

package badger

import "fmt"

type Key struct{}

func (Key) Namespace(name string) string {
	return fmt.Sprintf("%s", name)
}

func (Key) Service(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

func (Key) Deployment(namespace, service, name string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, service, name)
}

func (Key) Pod(namespace, service, deployment, name string) string {
	return fmt.Sprintf("%s:%s:%s:%s", namespace, service, deployment, name)
}

func (Key) Endpoint(namespace, service string) string {
	return fmt.Sprintf("%s:%s", namespace, service)
}

func (Key) Secret(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

func (Key) Config(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

func (Key) Volume(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

func (Key) Ingress(name string) string {
	return fmt.Sprintf("%s", name)
}

func (Key) Exporter(name string) string {
	return fmt.Sprintf("%s", name)
}

func (Key) Discovery(name string) string {
	return fmt.Sprintf("%s", name)
}

func (Key) Process(kind, hostname string, pid int, lead bool) string {
	if lead {
		return fmt.Sprintf("%s/lead", kind)
	}
	return fmt.Sprintf("%s:%s:%d", kind, hostname, pid)
}

func (Key) Manifest(name string) string {
	return fmt.Sprintf("%s", name)
}

func (Key) Node(name string) string {
	return fmt.Sprintf("%s", name)
}

func (Key) Route(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

func (Key) Subnet(name string) string {
	return fmt.Sprintf("%s", name)
}
