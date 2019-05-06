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

package main

import (
	"os"

	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/plugins/docker"
)

func main() {
	h := sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
	var level = 7
	levelVal := os.Getenv("LOG_LEVEL")
	switch levelVal {
	case "debug":
		level = 3
	case "info":
		level = 1
	}

	log.New(level)
	docker.Handlers(&h, docker.NewDriver())

	if err := h.ServeUnix("lb", 0); err != nil {
		panic(err)
	}
}
