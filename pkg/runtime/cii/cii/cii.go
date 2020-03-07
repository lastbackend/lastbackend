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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/runtime/cii"
	"github.com/lastbackend/lastbackend/pkg/runtime/cii/docker"
	"github.com/lastbackend/lastbackend/tools/log"
	"github.com/spf13/viper"
)

const (
	logLevel     = 5
	dockerDriver = "docker"
)

func New(v *viper.Viper) (cii.CII, error) {
	switch v.GetString("container.iri.type") {
	case dockerDriver:
		log.V(logLevel).Debugf("Use docker runtime interface for cii")
		cfg := docker.Config{}
		cfg.Host = v.GetString("container.iri.docker.host")
		cfg.Version = v.GetString("container.iri.docker.version")

		if v.IsSet("container.iri.docker.tls.verify") && v.GetBool("container.iri.docker.tls.verify") {
			cfg.TLS = new(docker.TLSConfig)
			cfg.TLS.CAPath = v.GetString("container.iri.docker.tls.ca_file")
			cfg.TLS.CertPath = v.GetString("container.iri.docker.tls.cert_file")
			cfg.TLS.KeyPath = v.GetString("container.iri.docker.tls.key_file")
		}

		return docker.New(cfg)
	default:
		return nil, fmt.Errorf("image runtime <%s> interface not supported", v.GetString("container.iri.type"))
	}
}
