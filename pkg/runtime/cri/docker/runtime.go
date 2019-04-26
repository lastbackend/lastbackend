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
	"github.com/docker/docker/api"
	"github.com/docker/docker/client"
	"github.com/lastbackend/lastbackend/pkg/log"
	"strconv"
)

const (
	logLevel = 5
	logPrefix = "runtime:docker"
)

type Runtime struct {
	client *client.Client
}

type config map[string]interface{}

func New(cfg config) (*Runtime, error) {

	var (
		err           error
		r             = new(Runtime)
		clientOptions = make([]client.Opt, 0)
	)

	if cfg == nil {
		cfg = make(config, 0)
	}

	log.V(logLevel).Debug("Use docker runtime interface")

	host := client.DefaultDockerHost
	if v, ok := cfg["host"]; ok {
		host = v.(string)
	}
	clientOptions = append(clientOptions, client.WithHost(host))

	version := api.DefaultVersion
	if v, ok := cfg["version"]; ok {
		switch i := v.(type) {
		case string:
			version = string(i)
		case float64:
			version = strconv.FormatFloat(float64(i), 'f', 2, 64)
		}
	}
	clientOptions = append(clientOptions, client.WithVersion(version))

	if v, ok := cfg["tls"]; ok {

		var tls = v.(map[string]string)
		var caFile, certFile, keyFile string

		if v, ok := tls["ca_file"]; ok {
			caFile = v
		}

		if v, ok := tls["cert_file"]; ok {
			certFile = v
		}

		if v, ok := tls["key_file"]; ok {
			keyFile = v
		}

		clientOptions = append(clientOptions, client.WithTLSClientConfig(caFile, certFile, keyFile))
	}

	r.client, err = client.NewClientWithOpts(clientOptions...)
	if err != nil {
		return nil, err
	}

	return r, nil
}
