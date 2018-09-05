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
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"io"
	"net/http"
)

const (
	logLevel = 3
)

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

func (r *Runtime) Pull(ctx context.Context, spec *types.ImageManifest, out io.Writer) (*types.Image, error) {

	log.V(logLevel).Debugf("Docker: Name pull: %s", spec.Name)

	image := new(types.Image)
	image.Meta.Name = spec.Name
	image.Status.State = types.StatusPull

	options := docker.ImagePullOptions{
		PrivilegeFunc: func() (string, error) {
			return "", errors.New("access denied")
		},
		RegistryAuth: spec.Auth,
	}

	res, err := r.client.ImagePull(ctx, spec.Name, options)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	const bufferSize = 1024
	var buffer = make([]byte, bufferSize)

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		default:

			readBytes, err := res.Read(buffer)
			if err != nil && err != io.EOF {
				return nil, err
			}
			if readBytes == 0 {
				ii, _, err := r.client.ImageInspectWithRaw(ctx, spec.Name)
				if err != nil {
					return nil, err
				}
				image.Meta.Hash = ii.ID
				image.Meta.Tags = ii.RepoTags
				image.Status.State = types.StateReady

				return image, nil
			}

			_, err = func(p []byte) (n int, err error) {

				if out != nil {
					n, err = out.Write(p)
					if err != nil {
						return n, err
					}

					if f, ok := out.(http.Flusher); ok {
						f.Flush()
					}
				}

				return n, nil
			}(buffer[0:readBytes])

			if err != nil {
				return nil, err
			}

			for i := 0; i < readBytes; i++ {
				buffer[i] = 0
			}
		}
	}
}

func (r *Runtime) Push(ctx context.Context, spec *types.ImageManifest, out io.Writer) (*types.Image, error) {

	log.V(logLevel).Debugf("Docker: Name push: %s", spec.Name)

	options := docker.ImagePushOptions{
		RegistryAuth: spec.Auth,
	}

	res, err := r.client.ImagePush(ctx, spec.Name, options)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	const bufferSize = 5e+6 //  5e+6 = 5MB
	var (
		readBytesLast = 0
		bufferLast    = make([]byte, bufferSize)
	)

	result := new(struct {
		Progress    map[string]interface{} `json:"progressDetail"`
		ErrorDetail *struct {
			Message string `json:"message"`
			Error   string `json:"error"`
		} `json:"errorDetail,omitempty"`
		Aux struct {
			Tag    string `json:"Tag"`
			Digest string `json:"Digest"`
			Size   int    `json:"Limit"`
		} `json:"aux"`
	})

	err = func(stream io.ReadCloser, data interface{}) error {
		for {
			buffer := make([]byte, bufferSize)
			readBytes, err := res.Read(buffer)
			if err != nil && err != io.EOF {
				return err
			}
			if readBytes == 0 {
				if err := json.Unmarshal(bufferLast[:readBytesLast], &data); err != nil {
					result = nil
					break
				}

				if result.ErrorDetail != nil {
					return fmt.Errorf("%s", result.ErrorDetail.Message)
				}

				break
			}

			bufferLast = make([]byte, bufferSize)

			readBytesLast = readBytes
			copy(bufferLast, buffer)

			if out != nil {
				out.Write(buffer[:readBytes])
			}
		}

		return nil
	}(res, result)
	if err != nil {
		return nil, err
	}

	imageID := spec.Name

	if result != nil {
		imageID = fmt.Sprintf("%s@%s", spec.Name, result.Aux.Digest)
	}

	image, err := r.Inspect(ctx, imageID)
	if err != nil {
		return nil, err
	}

	return image, err
}

func (r *Runtime) Build(ctx context.Context, stream io.Reader, spec *types.SpecBuildImage, out io.Writer) (*types.Image, error) {
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
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	const bufferSize = 1024
	var buffer = make([]byte, bufferSize)

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		default:

			readBytes, err := res.Body.Read(buffer)
			if err != nil && err != io.EOF {
				return nil, err
			}
			if readBytes == 0 {
				// TODO: get image info
				return new(types.Image), err
			}

			_, err = func(p []byte) (n int, err error) {

				if out != nil {
					n, err = out.Write(p)
					if err != nil {
						return n, err
					}

					if f, ok := out.(http.Flusher); ok {
						f.Flush()
					}
				}

				return n, nil
			}(buffer[0:readBytes])

			if err != nil {
				return nil, err
			}

			for i := 0; i < readBytes; i++ {
				buffer[i] = 0
			}
		}
	}
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

		img.Meta.ID = i.ID
		img.Meta.Tags = i.RepoTags
		img.Meta.Name = i.ID

		img.Status.Size = i.Size
		img.Status.VirtualSize = i.VirtualSize

		images = append(images, img)
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
