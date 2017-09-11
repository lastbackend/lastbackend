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

package app

import (
	"context"
	ctx "github.com/lastbackend/lastbackend/pkg/api/context"
	ins "github.com/lastbackend/lastbackend/pkg/api/app/interfaces"
	"github.com/lastbackend/lastbackend/pkg/api/app/routes/request"
	u "github.com/lastbackend/lastbackend/pkg/api/app/utils"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logLevel = 3

var utils ins.IUtil = new(u.Util)

type app struct {
	Context context.Context
}

func New(ctx context.Context) *app {
	return &app{ctx}
}

func SetUtils(u ins.IUtil) {
	utils = u
}

func (a *app) List() (types.AppList, error) {
	var (
		storage = ctx.Get().GetStorage()
		list    = types.AppList{}
	)

	log.V(logLevel).Debug("App: list app")

	apps, err := storage.App().List(a.Context)
	if err != nil {
		log.V(logLevel).Error("App: list app err: %s", err.Error())
		return list, err
	}

	log.V(logLevel).Debugf("App: list app result: %d", len(apps))

	for _, item := range apps {
		var app = item
		list = append(list, app)
	}

	return list, nil
}

func (a *app) Get(name string) (*types.App, error) {
	var (
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("App: get app %s", name)

	app, err := storage.App().GetByName(a.Context, name)
	if err != nil {
		if err.Error() == store.ErrKeyNotFound {
			log.V(logLevel).Warnf("App: app by name `%s` not found", name)
			return nil, nil
		}
		log.V(logLevel).Errorf("App: get app by name `%s` err: %s", name, err.Error())
		return nil, err
	}

	return app, nil
}

func (a *app) CheckExistServices(name string) (bool, error) {
	var (
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("App: check exist service for app %s", name)

	count, err := storage.Service().CountByApp(a.Context, name)
	if err != nil {
		log.V(logLevel).Errorf("App: get service list for app `%s` err: %s", name, err.Error())
		return false, err
	}

	return count > 0, nil
}

func (a *app) Create(rq *request.RequestAppCreateS) (*types.App, error) {

	var (
		err     error
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("App: create app %#v", rq)

	var app = types.App{}
	app.Meta.SetDefault()
	app.Meta.Name = utils.NameCreate(a.Context, rq.Name)
	app.Meta.Description = rq.Description

	if err = storage.App().Insert(a.Context, &app); err != nil {
		log.V(logLevel).Errorf("App: insert app err: %s", err.Error())
		return nil, err
	}

	return &app, nil
}

func (a *app) Update(n *types.App) (*types.App, error) {
	var (
		err     error
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("App: update app %#v", n)

	if err = storage.App().Update(a.Context, n); err != nil {
		log.V(logLevel).Errorf("App: update app err: %s", err.Error())
		return n, err
	}

	return n, nil
}

func (a *app) Remove(name string) error {
	var (
		err     error
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("App: remove app %s", name)

	err = storage.App().Remove(a.Context, name)
	if err != nil {
		log.V(logLevel).Errorf("App: remove app err: %s", err.Error())
		return err
	}

	return nil
}

func (a *app) WatchService(service chan *types.Service) error {
	var (
		storage = ctx.Get().GetStorage().Service()
	)

	log.V(logLevel).Debugf("App: watch services in app")

	if err := storage.PodsWatch(a.Context, service); err != nil {
		log.V(logLevel).Errorf("App: watch services in app err: %s", err.Error())
		return err
	}

	return nil
}
