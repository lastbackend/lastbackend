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
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"

	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"

	"context"
	"io"
)

const logLevel = 3

func (r *Runtime) Auth(ctx context.Context, secret *types.SecretAuthData) (string, error) {

	config := types.AuthConfig{
		Username: secret.Username,
		Password: secret.Password,
	}

	js, err := json.Marshal(config)
	if err != nil {
		return types.EmptyString, err
	}

	return base64.URLEncoding.EncodeToString(js), nil
}

func (r *Runtime) Pull(ctx context.Context, spec *types.ImageManifest) (*types.Image, error) {

	log.V(logLevel).Debugf("Docker: Name pull: %s", spec.Name)

	image := new(types.Image)
	image.Meta.Name = spec.Name
	image.Status.State = types.StateReady

	options := docker.ImagePullOptions{
		PrivilegeFunc: func() (string, error) {
			panic(0)
			return "", errors.New("Access denied")
		},
		RegistryAuth: spec.Auth,
	}

	rc, err := r.client.ImagePull(ctx, spec.Name, options)
	if err != nil {
		return nil, err
	}
	io.Copy(os.Stdout, rc)

	image.Status.State = types.StateReady
	return image, nil
}

func (r *Runtime) Push(ctx context.Context, spec *types.ImageManifest) error {

	log.V(logLevel).Debugf("Docker: Name push: %s", spec.Name)

	options := docker.ImagePushOptions{
		RegistryAuth: spec.Auth,
	}

	rc, err := r.client.ImagePush(ctx, spec.Name, options)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, rc)
	return err
}

func (r *Runtime) Build(ctx context.Context, stream io.Reader, spec *types.SpecBuildImage) (*types.Image, error) {
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

	img := new(types.Image)
	io.Copy(os.Stdout, res.Body)

	return img, err
}

func (r *Runtime) Remove(ctx context.Context, ID string) error {
	log.V(logLevel).Debugf("Docker: Name remove: %s", ID)
	var options docker.ImageRemoveOptions

	options = docker.ImageRemoveOptions{
		Force:         false,
		PruneChildren: true,
	}

	_, err := r.client.ImageRemove(ctx, ID, options)
	return err
}

func (r *Runtime) List(ctx context.Context) ([]*types.Image, error) {

	var images = make([]*types.Image, 0)

	il, err := r.client.ImageList(ctx, docker.ImageListOptions{All: true})
	if err != nil {
		return images, err
	}

	for _, i := range il {
		img := new(types.Image)

		img.Meta.Name = i.ID
		img.Meta.ID = i.ID
		img.Status.Size = i.Size
		img.Status.VirtualSize = i.VirtualSize

	}

	return images, nil
}

func (r *Runtime) Inspect(ctx context.Context, id string) (*types.Image, error) {
	info, _, err := r.client.ImageInspectWithRaw(ctx, id)

	image := new(types.Image)
	image.Meta.ID = info.ID
	image.Status.Size = info.Size
	image.Status.VirtualSize = info.VirtualSize

	return image, err
}
