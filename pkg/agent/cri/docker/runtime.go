package docker

import (
	"github.com/docker/docker/api"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/lastbackend/lastbackend/pkg/agent/config"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"net/http"
	"path/filepath"
)

type Runtime struct {
	client *client.Client
}

func New(cfg *config.Docker) (*Runtime, error) {

	var cli *http.Client
	var err error

	log := context.Get().GetLogger()

	log.Debug("Use docker CRI")
	r := &Runtime{}

	if *cfg.Certs != "" {

		log.Debugf("Create Docker secure client: %s", *cfg.Certs)

		options := tlsconfig.Options{
			CAFile:             filepath.Join(*cfg.Certs, "ca.pem"),
			CertFile:           filepath.Join(*cfg.Certs, "cert.pem"),
			KeyFile:            filepath.Join(*cfg.Certs, "key.pem"),
			InsecureSkipVerify: *cfg.TLS,
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

	host := *cfg.Host
	if host == "" {
		host = client.DefaultDockerHost
	}

	version := *cfg.Version
	if version == "" {
		version = api.DefaultVersion
	}

	r.client, err = client.NewClient(host, version, cli, nil)
	if err != nil {
		return r, err
	}

	return r, nil
}
