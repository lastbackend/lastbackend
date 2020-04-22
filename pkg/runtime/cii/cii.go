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
	"github.com/lastbackend/lastbackend/pkg/runtime/cii/containerd"
	"github.com/spf13/viper"
)

const (
	dockerDriver     = "docker"
	containerDDriver = "containerd"
)

// IMI - Image System Interface
type CII interface {
	Auth(ctx context.Context, secret *models.SecretAuthData) (string, error)
	Pull(ctx context.Context, spec *models.ImageManifest, out io.Writer) (*models.Image, error)
	Remove(ctx context.Context, image string) error
	Push(ctx context.Context, spec *models.ImageManifest, out io.Writer) (*models.Image, error)
	Build(ctx context.Context, stream io.Reader, spec *models.SpecBuildImage, out io.Writer) (*models.Image, error)
	List(ctx context.Context) ([]*models.Image, error)
	Inspect(ctx context.Context, id string) (*models.Image, error)
	Subscribe(ctx context.Context) (chan *models.Image, error)
	Close() error
}

func New(v *viper.Viper) (CII, error) {
	switch v.GetString("container.iri.type") {
	//case dockerDriver:
	//	cfg := docker.Config{}
	//	cfg.Host = v.GetString("container.iri.docker.host")
	//	cfg.Version = v.GetString("container.iri.docker.version")
	//
	//	if v.IsSet("container.iri.docker.tls.verify") && v.GetBool("container.iri.docker.tls.verify") {
	//		cfg.TLS = new(docker.TLSConfig)
	//		cfg.TLS.CAPath = v.GetString("container.iri.docker.tls.ca_file")
	//		cfg.TLS.CertPath = v.GetString("container.iri.docker.tls.cert_file")
	//		cfg.TLS.KeyPath = v.GetString("container.iri.docker.tls.key_file")
	//	}
	//	return docker.New(cfg)
	case containerDDriver:
		cfg := containerd.Config{}
		//cfg.Host = v.GetString("container.iri.containerd.host")
		//cfg.Version = v.GetString("container.iri.containerd.version")

		if v.IsSet("container.iri.containerd.tls.verify") && v.GetBool("container.iri.containerd.tls.verify") {
			cfg.TLS = new(containerd.TLSConfig)
			cfg.TLS.CAPath = v.GetString("container.iri.containerd.tls.ca_file")
			cfg.TLS.CertPath = v.GetString("container.iri.containerd.tls.cert_file")
			cfg.TLS.KeyPath = v.GetString("container.iri.containerd.tls.key_file")
		}
		return containerd.New(cfg)
	default:
		return nil, fmt.Errorf("image runtime <%s> interface not supported", v.GetString("container.iri.type"))
	}
}
