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

package runtime

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"strings"
)

func ConfigManage(ctx context.Context, name string, cfg *types.ConfigManifest) error {

	log.V(logLevel).Debugf("Manage config: %s", name)

	if cfg.State == types.StateDestroyed {
		ConfigRemove(ctx, name)
		return nil
	}

	ok := envs.Get().GetState().Configs().GetConfig(name)
	if ok != nil {
		return ConfigUpdate(ctx, name, cfg)
	}

	return ConfigCreate(ctx, name, cfg)
}

func ConfigCreate(ctx context.Context, name string, cfg *types.ConfigManifest) error {

	log.V(logLevel).Debugf("create config: %s", name)

	ok := envs.Get().GetState().Configs().GetConfig(name)
	if ok != nil {
		return nil
	}

	envs.Get().GetState().Configs().AddConfig(name, cfg)
	return nil
}

func ConfigUpdate(ctx context.Context, name string, cfg *types.ConfigManifest) error {

	log.V(logLevel).Debugf("update config: %s", name)

	envs.Get().GetState().Configs().SetConfig(name, cfg)
	return nil

}

func ConfigRemove(ctx context.Context, name string) {

	log.V(logLevel).Debugf("remove config: %s", name)

	envs.Get().GetState().Configs().DelConfig(name)
}

func parseConfigSelflink(selflink string) (string, string) {
	var namespace, name string

	parts := strings.Split(selflink, ":")

	if len(parts) == 1 {
		namespace = types.DEFAULT_NAMESPACE
		name = parts[0]
	}

	if len(parts) > 1 {
		namespace = parts[0]
		name = parts[1]
	}

	return namespace, name

}
