package cri

import (
	"github.com/lastbackend/lastbackend/pkg/agent/config"
	"github.com/lastbackend/lastbackend/pkg/agent/cri/docker"
	"github.com/lastbackend/lastbackend/pkg/agent/cri"
)

func New (cfg *config.Runtime) (cri.CRI, error) {
	var cri cri.CRI
	var err error

	switch *cfg.CRI {
	case "docker":
		cri, err = docker.New(cfg.Docker)
	}

	if err != nil {
		return cri, err
	}

	return cri, err
}
