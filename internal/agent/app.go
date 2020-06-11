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

package agent

import (
	"context"

	"github.com/lastbackend/lastbackend/internal/agent/config"
	"github.com/lastbackend/lastbackend/internal/agent/controller"
	"github.com/lastbackend/lastbackend/internal/agent/runtime"
	"github.com/lastbackend/lastbackend/internal/agent/server"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/pkg/errors"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	config config.Config
	storage storage.IStorage

	HttpServer *server.HttpServer
	Runtime    *runtime.Runtime
	Controller *controller.Controller
}

func New(stg storage.IStorage, config config.Config) (*App, error) {

	loggerOpts := logger.Configuration{}
	loggerOpts.EnableConsole = true

	if config.Debug {
		loggerOpts.ConsoleLevel = logger.Debug
	}
	if err := logger.NewLogger(loggerOpts); err != nil {
		return nil, errors.Wrapf(err, "cat not logger initialize")
	}

	app := new(App)
	app.ctx, app.cancel = context.WithCancel(context.Background())
	app.config = config
	app.storage = stg

	if err := app.init(); err != nil {
		return nil, errors.Wrapf(err, "can not be init application")
	}

	return app, nil
}

func (app *App) Run() error {
	log := logger.WithContext(context.Background())
	log.Infof("Run minion service")

	if err := app.Runtime.Run(); err != nil {
		return errors.Wrapf(err, "can not be run runtime")
	}

	if app.Controller != nil {
		if err := app.Controller.Connect(app.config); err != nil {
			return errors.Wrapf(err, "can not be connect controller")
		}
		go app.Controller.Subscribe()
		go app.Controller.Sync()
	}

	go func() {
		if err := app.HttpServer.Run(); err != nil {
			log.Fatal(errors.Wrapf(err, "can not be run http rest server"))
			return
		}
	}()

	return nil
}

func (app *App) Stop() {
	app.Runtime.Stop()
	app.cancel()
}

func (app *App) init() error {

	var err error

	log := logger.WithContext(context.Background())
	log.Infof("Init agent service")

	app.Runtime, err = runtime.New(app.storage, app.config)
	if err != nil {
		return errors.Wrapf(err, "can not be runtime initialize")
	}

	app.Controller, err = controller.New(app.Runtime)
	if err != nil {
		return errors.Wrapf(err, "can not be controller initialize")
	}

	app.HttpServer = server.NewServer(app.Runtime.GetState(), app.storage, app.config)

	return nil
}
