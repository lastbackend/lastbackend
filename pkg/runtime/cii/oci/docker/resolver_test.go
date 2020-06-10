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

package docker

import (
	"context"
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolver(t *testing.T) {

	ctx := context.Background()

	cfg := Config{}
	cfg.Debug = true

	hostOptions := HostOptions{}
	hostOptions.DefaultScheme = "http"
	//hostOptions.HostDir = HostDirFromRoot(hostDir)
	hostOptions.DefaultTLS = &tls.Config{
		InsecureSkipVerify: true,
	}

	cfg.Hosts = ConfigureHosts(ctx, hostOptions)

	r, err := New(cfg)
	if assert.NotNil(t, err, "create resolver error") {
		return
	}
	if !assert.NotNil(t, r, "resolver can not be nil") {
		return
	}

	reader, err := r.Pull(ctx, "localhost:5000/redis")
	if assert.NotNil(t, err, "pull image failed") {
		return
	}
}
