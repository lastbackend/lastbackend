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

package image

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"strings"
)

type util struct {
	IUtil
}

func (util) Name(_ context.Context, hub, name string) string {
	cfg := c.Get().GetConfig()
	namespace := name
	if cfg.Registry.Username != "" {
		namespace = fmt.Sprintf("%s/%s", cfg.Registry.Username, namespace)
	}
	if cfg.Registry.Server != "" {
		server := cfg.Registry.Server
		switch true {
		case strings.HasPrefix(server, "http://") == true:
			server = server[7:]
		case strings.HasPrefix(server, "https://") == true:
			server = server[8:]
		}
		namespace = fmt.Sprintf("%s/%s", server, namespace)
	}
	return namespace
}

func (util) RegistryAuth(_ context.Context, _ string) *types.RegistryAuth {
	cfg := c.Get().GetConfig()
	return &types.RegistryAuth{
		Username: cfg.Registry.Username,
		Password: cfg.Registry.Password,
		Server:   cfg.Registry.Server,
	}
}
