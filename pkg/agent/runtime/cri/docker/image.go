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
	"github.com/lastbackend/lastbackend/pkg/common/context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"io"
)

func (r *Runtime) ImagePull(ctx context.IContext, spec *types.ImageSpec) (io.ReadCloser, error) {
	log := ctx.GetLogger()
	log.Debugf("Docker: Image pull: %s", spec.Name)
	options := docker.ImagePullOptions{
		RegistryAuth: spec.Auth,
	}
	return r.client.ImagePull(ctx.Background(), spec.Name, options)
}

func (r *Runtime) ImagePush(ctx context.IContext) {}

func (r *Runtime) ImageBuild(ctx context.IContext) {}

func (r *Runtime) ImageRemove(ctx context.IContext, ID string) error {

	log := ctx.GetLogger()
	log.Debugf("Docker: Image remove: %s", ID)
	var options docker.ImageRemoveOptions

	options = docker.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	}

	_, err := r.client.ImageRemove(ctx.Background(), ID, options)
	return err
}

func (r *Runtime) ImageList(ctx context.IContext) {
	_, err := r.client.ImageList(ctx.Background(), docker.ImageListOptions{All: true})
	if err != nil {
		return
	}
}

func (r *Runtime) ImageInspect(ctx context.IContext, name, tag string) error {
	_, _, err := r.client.ImageInspectWithRaw(ctx.Background(), name)
	if err != nil {
		return err
	}

	return nil
}
