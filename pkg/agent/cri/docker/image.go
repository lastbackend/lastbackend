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
	"context"
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"io"
)

func (r *Runtime) ImagePull(spec *types.ImageSpec) (io.ReadCloser, error) {
	options := docker.ImagePullOptions{
		RegistryAuth: spec.Auth,
	}
	return r.client.ImagePull(context.Background(), spec.Name, options)
}

func (r *Runtime) ImagePush() {}

func (r *Runtime) ImageBuild() {}

func (r *Runtime) ImageRemove(ID string) error {

	var options docker.ImageRemoveOptions

	options = docker.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	}

	_, err := r.client.ImageRemove(context.Background(), ID, options)
	return err
}

func (r *Runtime) ImageList() {
	_, err := r.client.ImageList(context.Background(), docker.ImageListOptions{All: true})
	if err != nil {
		return
	}
}

func (r *Runtime) ImageInspect(name, tag string) error {
	_, _, err := r.client.ImageInspectWithRaw(context.Background(), name)
	if err != nil {
		return err
	}

	return nil
}
