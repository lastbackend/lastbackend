//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"errors"

	"github.com/docker/docker/api"
	"github.com/docker/docker/client"
)

const (
	logLevel  = 5
	logPrefix = "runtime:docker"
)

type Runtime struct {
	client *client.Client
}

type Config struct {
	Host    string
	Version string
	TLS     *TLSConfig
}

type TLSConfig struct {
	CAPath   string
	CertPath string
	KeyPath  string
}

func New(cfg Config) (*Runtime, error) {

	var (
		err           error
		r             = new(Runtime)
		clientOptions = make([]client.Opt, 0)
	)

	host := client.DefaultDockerHost
	if len(cfg.Host) != 0 {
		host = cfg.Host
	}
	clientOptions = append(clientOptions, client.WithHost(host))

	version := api.DefaultVersion
	if len(cfg.Version) != 0 {
		version = cfg.Version
	}

	clientOptions = append(clientOptions, client.WithVersion(version))

	if cfg.TLS != nil {

		if len(cfg.TLS.CAPath) == 0 {
			return nil, errors.New("path to ca file not set")
		}
		if len(cfg.TLS.CertPath) == 0 {
			return nil, errors.New("path to cert file not set")
		}
		if len(cfg.TLS.KeyPath) == 0 {
			return nil, errors.New("path to key file not set")
		}

		caPath := cfg.TLS.CAPath
		certPath := cfg.TLS.CertPath
		keyPath := cfg.TLS.KeyPath

		clientOptions = append(clientOptions, client.WithTLSClientConfig(caPath, certPath, keyPath))
	}

	r.client, err = client.NewClientWithOpts(clientOptions...)
	if err != nil {
		return nil, err
	}

	return r, nil
}
