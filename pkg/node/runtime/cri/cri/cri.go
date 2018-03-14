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

package cri

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cri/docker"
	"github.com/spf13/viper"
)

func New() (cri.CRI, error) {
	switch viper.GetString("node.cri.type") {
	case "docker":
		return docker.New()
	default:
		return nil, fmt.Errorf("container runtime <%s> interface not supported", viper.GetString("node.cri.type"))
	}
}
