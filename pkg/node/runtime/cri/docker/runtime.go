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
	"github.com/docker/docker/api"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/lastbackend/lastbackend/pkg/log"

	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
)

type Runtime struct {
	client *client.Client
}

func New() (*Runtime, error) {

	var (
		err error
		cli *http.Client
		r   = new(Runtime)
	)

	log.V(logLevel).Debug("Use docker CRI")

	if viper.GetString("runtime.docker.certs") != "" {

		log.V(logLevel).Debugf("Create Docker secure client: %s", viper.GetString("runtime.docker.certs"))

		options := tlsconfig.Options{
			CAFile:             filepath.Join(viper.GetString("runtime.docker.certs"), "ca.pem"),
			CertFile:           filepath.Join(viper.GetString("runtime.docker.certs"), "cert.pem"),
			KeyFile:            filepath.Join(viper.GetString("runtime.docker.certs"), "key.pem"),
			InsecureSkipVerify: viper.GetBool("runtime.docker.ssl"),
		}

		tlsc, err := tlsconfig.Client(options)
		if err != nil {
			return nil, err
		}

		cli = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsc,
			},
		}
	}

	host := client.DefaultDockerHost
	if viper.GetString("runtime.docker.host") != "" {
		host = viper.GetString("runtime.docker.host")
	}

	version := api.DefaultVersion
	if viper.GetString("runtime.docker.version") != "" {
		version = viper.GetString("runtime.docker.version")
	}

	r.client, err = client.NewClient(host, version, cli, nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func NewWithHost(host string) (*Runtime, error) {

	var (
		err error
		cli *http.Client
		r   = new(Runtime)
	)

	log.V(logLevel).Debugf("Use docker CRI with host %s", host)

	if viper.GetString("runtime.docker.certs") != "" {

		log.V(logLevel).Debugf("Create Docker secure client: %s", viper.GetString("runtime.docker.certs"))

		options := tlsconfig.Options{
			CAFile:             filepath.Join(viper.GetString("runtime.docker.certs"), "ca.pem"),
			CertFile:           filepath.Join(viper.GetString("runtime.docker.certs"), "cert.pem"),
			KeyFile:            filepath.Join(viper.GetString("runtime.docker.certs"), "key.pem"),
			InsecureSkipVerify: viper.GetBool("runtime.docker.ssl"),
		}

		tlsc, err := tlsconfig.Client(options)
		if err != nil {
			return nil, err
		}

		cli = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsc,
			},
		}
	}

	version := api.DefaultVersion
	if viper.GetString("runtime.docker.version") != "" {
		version = viper.GetString("runtime.docker.version")
	}

	r.client, err = client.NewClient(host, version, cli, nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}
