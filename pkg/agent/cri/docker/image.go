package docker

import (
	"context"
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"io"
)

func (r *Runtime) ImagePull(spec types.ImageSpec) (io.ReadCloser, error) {
	var image, auth string
	image = spec.Image()
	auth = spec.Auth()
	options := docker.ImagePullOptions{
		RegistryAuth: auth,
	}
	return r.client.ImagePull(context.Background(), image, options)
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

func (r *Runtime) ImageList() {}
