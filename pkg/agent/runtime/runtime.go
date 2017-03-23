package runtime

import (
	"context"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/lastbackend/lastbackend/libs/model"
)

type Runtime struct {
	Context    context.Context
	Client     *client.Client
	Containers *model.ContainerList
	Images     *model.ImageList
}

func New() *Runtime {
	var runtime = new(Runtime)

	return runtime
}

func (r *Runtime) Init() {

	r.Client, _ = client.NewEnvClient()

	// Get Container list
	containers, _ := r.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	cj, _ := json.Marshal(containers)
	logrus.Debugf("%s", cj)

	// Get Images list
	images, _ := r.Client.ImageList(context.Background(), types.ImageListOptions{})
	ci, _ := json.Marshal(images)
	logrus.Debugf("%s", ci)

	events, _ := r.Client.Events(context.Background(), types.EventsOptions{})
	go func() {
		for {
			select {
			case e := <-events:
				if e.Type == "container" && e.Status != "destroy" {
					ej, _ := json.Marshal(e)
					logrus.Debugf("%s", ej)
					continue
				}
			}
		}
	}()

}
