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
	"errors"
	"fmt"

	"github.com/lastbackend/lastbackend/internal/agent/config"
	"github.com/lastbackend/lastbackend/internal/agent/controller"
	"github.com/lastbackend/lastbackend/internal/agent/runtime"
	"github.com/lastbackend/lastbackend/internal/agent/server"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/util/filesystem"
	"github.com/lastbackend/lastbackend/tools/logger"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	config config.Config

	HttpServer *server.HttpServer
	Runtime    *runtime.Runtime
	Controller *controller.Controller
}

func New(config config.Config) (*App, error) {

	loggerOpts := logger.Configuration{}
	loggerOpts.EnableConsole = true

	if config.Debug {
		loggerOpts.ConsoleLevel = logger.Debug
	}
	if err := logger.NewLogger(loggerOpts); err != nil {
		return nil, errors.New("logger initialize failed")
	}

	app := new(App)
	app.ctx, app.cancel = context.WithCancel(context.Background())
	app.config = config

	if err := app.init(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run() error {
	log := logger.WithContext(context.Background())
	log.Infof("Run minion service")

	if err := app.Runtime.Run(); err != nil {
		log.Errorf("Run runtime err: %v", err)
		return err
	}

	if app.Controller != nil {
		if err := app.Controller.Connect(app.config); err != nil {
			return fmt.Errorf("Connect controller err: %v", err)
		}

		go app.Controller.Subscribe()
		go app.Controller.Sync()
	}

	go func() {
		if err := app.HttpServer.Run(); err != nil {
			log.Fatalf("Run http rest server err: %v", err)
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
	log.Infof("Init minion service")

	workdir, err := filesystem.DetermineHomePath(app.config.WorkDir, app.config.Rootless)
	if err != nil {
		return err
	}

	if err := filesystem.MkDir(workdir, 0755); err != nil {
		return err
	}
	fmt.Println("workdir >>>", workdir)
	stg, err := storage.Get(storage.BboltDriver, storage.BboltConfig{DbDir: workdir, DbName: ".agent-db"})
	if err != nil {
		return fmt.Errorf("cannot initialize storage: %v", err)
	}

	app.Runtime, err = runtime.New(app.config)
	if err != nil {
		return fmt.Errorf("Can not initialize runtime: %v", err)
	}

	app.Controller, err = controller.New(app.Runtime)
	if err != nil {
		return fmt.Errorf("Can not initialize controller: %v", err)
	}

	app.HttpServer = server.NewServer(app.Runtime.GetState(), stg, app.config)

	return nil
}
