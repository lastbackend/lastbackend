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

package docker

import "github.com/lastbackend/lastbackend/pkg/util/proxy"

type Message struct {
	Data             string            `json:"message"`
	ContainerId      string            `json:"container_id"`
	ContainerName    string            `json:"container_name"`
	Selflink         string            `json:"selflink"`
	ContainerCreated proxy.JsonTime    `json:"container_created"`
	Tag              string            `json:"tag"`
	Extra            map[string]string `json:"extra"`
	Host             string            `json:"host"`
	Timestamp        proxy.JsonTime    `json:"timestamp"`
}
