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
