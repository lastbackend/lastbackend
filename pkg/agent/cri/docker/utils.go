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
	c.Image.ID = dc.ImageID
	c.Image.Name = dc.Image
	return c

}
