//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package cii

import (
	"context"
	"fmt"
	"io"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/runtime/cii/oci"
)

const (
	DockerDriver = "docker"
)

// IMI - Image System Interface
type CII interface {
	Auth(ctx context.Context, secret *models.SecretAuthData) (string, error)
	Pull(ctx context.Context, spec *models.ImageManifest, out io.Writer) (*models.Image, error)
	Remove(ctx context.Context, image string) error
	Push(ctx context.Context, spec *models.ImageManifest, out io.Writer) (*models.Image, error)
	Build(ctx context.Context, stream io.Reader, spec *models.SpecBuildImage, out io.Writer) (*models.Image, error)
	List(ctx context.Context, filters ...string) ([]*models.Image, error)
	Inspect(ctx context.Context, id string) (*models.Image, error)
	Subscribe(ctx context.Context) (chan *models.Image, error)
	Close() error
}

type DockerConfig oci.ConfigDocker

func New(driver string, opts interface{}) (CII, error) {

	if opts == nil {
		return nil, fmt.Errorf("options can not be is nil")
	}

	switch driver {
	case DockerDriver:
		o := opts.(DockerConfig)
		return oci.NewDocker(oci.ConfigDocker(o))
	default:
		return nil, fmt.Errorf("container image runtime <%s> interface not supported", driver)
	}
}
