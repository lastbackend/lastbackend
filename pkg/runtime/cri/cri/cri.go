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

package cri

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/runtime/cri/docker"
	"github.com/spf13/viper"
)

const (
	logLevel     = 5
	dockerDriver = "docker"
	runcDriver   = "runc"
)

func New(v *viper.Viper) (cri.CRI, error) {
	switch v.GetString("container.cri.type") {
	case dockerDriver:
		log.V(logLevel).Debugf("Use docker runtime interface for cri")

		cfg := docker.Config{}
		cfg.Host = v.GetString("container.cri.docker.host")
		cfg.Version = v.GetString("container.cri.docker.version")

		if v.IsSet("container.cri.docker.tls.verify") && v.GetBool("container.cri.docker.tls.verify") {
			cfg.TLS = new(docker.TLSConfig)
			cfg.TLS.CAPath = v.GetString("container.cri.docker.tls.ca_file")
			cfg.TLS.CertPath = v.GetString("container.cri.docker.tls.cert_file")
			cfg.TLS.KeyPath = v.GetString("container.cri.docker.tls.key_file")
		}

		return docker.New(cfg)
	default:
		return nil, fmt.Errorf("container runtime <%s> interface not supported", v.GetString("container.cri.type"))
	}
}
