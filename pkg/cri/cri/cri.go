//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/pkg/config"
	"github.com/lastbackend/lastbackend/pkg/cri"
	"github.com/lastbackend/lastbackend/pkg/cri/docker"
)

func New(cfg config.Runtime) (cri.CRI, error) {
	var cri cri.CRI
	var err error

	switch cfg.CRI {
	case "docker":
		cri, err = docker.New(cfg.Docker)
	}

	if err != nil {
		return cri, err
	}

	return cri, err
}
