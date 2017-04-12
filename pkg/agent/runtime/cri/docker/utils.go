//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

import (
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"time"
)

func GetContainer(dc docker.Container, info docker.ContainerJSON) *types.Container {

	var c *types.Container

	c = &types.Container{
		ID:      dc.ID,
		State:   dc.State,
		Status:  dc.Status,
		Created: time.Unix(dc.Created, 0),
	}

	t, _ := time.Parse(time.RFC3339Nano, info.State.StartedAt)
	c.Started = t
	c.Image = dc.Image
	return c

}
