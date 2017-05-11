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
	"github.com/lastbackend/lastbackend/pkg/common/config"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime/cri/docker"
	"github.com/pkg/errors"
)

func New(cfg config.Runtime) (cri.CRI, error) {
	switch cfg.CRI {
	case "docker":
		return docker.New(cfg.Docker)
	default:
		return nil, errors.New(`container runtime interface not support`)
	}
}
