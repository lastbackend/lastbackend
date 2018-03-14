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

import (
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"

	"context"
	"io"
)

func (r *Runtime) ImagePull(ctx context.Context, spec *types.SpecTemplateContainerImage) (io.ReadCloser, error) {
	log.Debugf("Docker: Image pull: %s", spec.Name)

	options := docker.ImagePullOptions{
		RegistryAuth: spec.Auth,
		PrivilegeFunc: func() (string, error) {
			panic(0)
			return "", errors.New("Access denied")
		},
	}
	return r.client.ImagePull(ctx, spec.Name, options)
}

func (r *Runtime) ImagePush(ctx context.Context, spec *types.SpecTemplateContainerImage) (io.ReadCloser, error) {
	log.Debugf("Docker: Image push: %s", spec.Name)
	options := docker.ImagePushOptions{
		RegistryAuth: spec.Auth,
	}
	return r.client.ImagePush(ctx, spec.Name, options)
}

func (r *Runtime) ImageBuild(ctx context.Context, stream io.Reader, spec *types.SpecBuildImage) (io.ReadCloser, error) {
	options := docker.ImageBuildOptions{
		Tags:           spec.Tags,
		Memory:         spec.Memory,
		Dockerfile:     spec.Dockerfile,
		ExtraHosts:     spec.ExtraHosts,
		Context:        spec.Context,
		NoCache:        spec.NoCache,
		SuppressOutput: spec.SuppressOutput,
	}
	if spec.AuthConfigs != nil {
		options.AuthConfigs = make(map[string]docker.AuthConfig, 0)
		for k, v := range spec.AuthConfigs {
			options.AuthConfigs[k] = docker.AuthConfig(v)
		}
	}
	res, err := r.client.ImageBuild(ctx, stream, options)
	return res.Body, err
}

func (r *Runtime) ImageRemove(ctx context.Context, ID string) error {
	log.Debugf("Docker: Image remove: %s", ID)
	var options docker.ImageRemoveOptions

	options = docker.ImageRemoveOptions{
		Force:         false,
		PruneChildren: true,
	}

	_, err := r.client.ImageRemove(ctx, ID, options)
	return err
}

func (r *Runtime) ImageList(ctx context.Context) ([]docker.ImageSummary, error) {
	return r.client.ImageList(ctx, docker.ImageListOptions{All: true})
}

func (r *Runtime) ImageInspect(ctx context.Context, id string) (*types.ImageInfo, []byte, error) {
	info, buf, err := r.client.ImageInspectWithRaw(ctx, id)

	image := new(types.ImageInfo)
	image.ID = info.ID
	image.Size = info.Size
	image.VirtualSize = info.VirtualSize

	return image, buf, err
}
