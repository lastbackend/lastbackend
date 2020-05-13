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

package cri

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/runtime/cri/oci"
)

const (
	OCIDriver = "oci"
)

// CRI - Container System Interface
type CRI interface {
	List(ctx context.Context, all bool) ([]*models.Container, error)
	Create(ctx context.Context, spec *models.ContainerManifest) (string, error)
	Start(ctx context.Context, ID string) error
	Restart(ctx context.Context, ID string, timeout *time.Duration) error
	Stop(ctx context.Context, ID string, timeout *time.Duration) error
	Pause(ctx context.Context, ID string) error
	Resume(ctx context.Context, ID string) error
	Remove(ctx context.Context, ID string, clean bool, force bool) error
	Inspect(ctx context.Context, ID string) (*models.Container, error)
	Logs(ctx context.Context, ID string, stdout, stderr, follow bool) (io.ReadCloser, error)
	Copy(ctx context.Context, ID, path string, content io.Reader) error
	Wait(ctx context.Context, ID string) error
	Subscribe(ctx context.Context, container chan *models.Container) error
	Close() error
}

type OciConfig oci.Config

func New(driver string, opts interface{}) (CRI, error) {

	if opts == nil {
		return nil, fmt.Errorf("options can not be is nil")
	}

	switch driver {
	case OCIDriver:
		o := opts.(OciConfig)
		return oci.New(oci.Config(o))
	default:
		return nil, fmt.Errorf("container runtime <%s> interface not supported", driver)
	}
}
