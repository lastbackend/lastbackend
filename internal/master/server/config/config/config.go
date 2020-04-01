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

package config

import (
	"context"
	"net/http"

	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logPrefix = "api:handler:config"
	logLevel  = 3
)

func Fetch(ctx context.Context, namespace, name string) (*models.Config, *errors.Err) {

	cm := service.NewConfigModel(ctx, envs.Get().GetStorage())
	cfg, err := cm.Get(namespace, name)

	if err != nil {
		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("config").InternalServerError(err)
	}

	if cfg == nil {
		err := errors.New("namespace not found")
		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("config").NotFound()
	}

	return cfg, nil
}

func Apply(ctx context.Context, ns *models.Namespace, mf *request.ConfigManifest) (*models.Config, *errors.Err) {

	if mf.Meta.Name == nil {
		return nil, errors.New("config").BadParameter("meta.name")
	}

	cfg, err := Fetch(ctx, ns.Meta.Name, *mf.Meta.Name)
	if err != nil {
		if err.Code != http.StatusText(http.StatusNotFound) {
			return nil, errors.New("config").InternalServerError()
		}
	}

	if cfg == nil {
		return Create(ctx, ns, mf)
	}

	return Update(ctx, ns, cfg, mf)
}

func Create(ctx context.Context, ns *models.Namespace, mf *request.ConfigManifest) (*models.Config, *errors.Err) {

	cm := service.NewConfigModel(ctx, envs.Get().GetStorage())
	if mf.Meta.Name != nil {

		cf, err := cm.Get(ns.Meta.Name, *mf.Meta.Name)
		if err != nil {
			log.Errorf("%s:create:> get config by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
			return nil, errors.New("config").InternalServerError()
		}

		if cf != nil {
			log.Warnf("%s:create:> config name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
			return nil, errors.New("config").NotUnique("name")
		}
	}

	cfg := new(models.Config)
	cfg.Meta.SetDefault()
	cfg.Meta.Namespace = ns.Meta.Name
	cfg.Meta.SelfLink = *models.NewConfigSelfLink(ns.Meta.Name, *mf.Meta.Name)

	mf.SetConfigMeta(cfg)
	mf.SetConfigSpec(cfg)

	if _, err := cm.Create(ns, cfg); err != nil {
		log.Errorf("%s:create:> create config err: %s", logPrefix, ns.Meta.Name, err.Error())
		return nil, errors.New("config").InternalServerError()
	}

	return cfg, nil
}

func Update(ctx context.Context, ns *models.Namespace, cfg *models.Config, mf *request.ConfigManifest) (*models.Config, *errors.Err) {

	cm := service.NewConfigModel(ctx, envs.Get().GetStorage())

	mf.SetConfigMeta(cfg)
	mf.SetConfigSpec(cfg)

	if _, err := cm.Update(cfg); err != nil {
		log.Errorf("%s:update:> update config err: %s", logPrefix, err.Error())
		return nil, errors.New("config").InternalServerError()
	}

	return cfg, nil
}
