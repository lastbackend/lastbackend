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

package runtime

import (
	"context"
	"github.com/lastbackend/lastbackend/tools/logger"
	"strings"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

func (r Runtime) ConfigManage(ctx context.Context, name string, cfg *models.ConfigManifest) error {
	log := logger.WithContext(context.Background())
	log.Debugf("Manage config: %s", name)

	if cfg.State == models.StateDestroyed {
		r.ConfigRemove(ctx, name)
		return nil
	}

	ok := r.state.Configs().GetConfig(name)
	if ok != nil {
		return r.ConfigUpdate(ctx, name, cfg)
	}

	return r.ConfigCreate(ctx, name, cfg)
}

func (r Runtime) ConfigCreate(ctx context.Context, name string, cfg *models.ConfigManifest) error {
	log := logger.WithContext(context.Background())
	log.Debugf("create config: %s", name)

	ok := r.state.Configs().GetConfig(name)
	if ok != nil {
		return nil
	}

	r.state.Configs().AddConfig(name, cfg)
	return nil
}

func (r Runtime) ConfigUpdate(ctx context.Context, name string, cfg *models.ConfigManifest) error {
	log := logger.WithContext(context.Background())
	log.Debugf("update config: %s", name)

	r.state.Configs().SetConfig(name, cfg)
	return nil

}

func (r Runtime) ConfigRemove(ctx context.Context, name string) {
	log := logger.WithContext(context.Background())
	log.Debugf("remove config: %s", name)

	r.state.Configs().DelConfig(name)
}

func (r Runtime) parseConfigSelflink(selflink string) (string, string) {
	var namespace, name string

	parts := strings.Split(selflink, ":")

	if len(parts) == 1 {
		namespace = models.DEFAULT_NAMESPACE
		name = parts[0]
	}

	if len(parts) > 1 {
		namespace = parts[0]
		name = parts[1]
	}

	return namespace, name

}
