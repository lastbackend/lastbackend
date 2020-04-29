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

package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/server/config"
	"github.com/lastbackend/lastbackend/internal/server/server"
	"github.com/lastbackend/lastbackend/internal/server/state"
	"github.com/lastbackend/lastbackend/internal/util/filesystem"
	"github.com/lastbackend/lastbackend/tools/logger"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	config config.Config

	HttpServer *server.HttpServer
	State      *state.State
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

func (app App) Run() error {
	log := logger.WithContext(context.Background())

	log.Infof("Run master service")

	go func() {
		if err := app.HttpServer.Run(); err != nil {
			log.Fatalf("http rest server start err: %v", err)
			return
		}
	}()

	return nil
}

func (app *App) Stop() {
	app.cancel()
}

func (app *App) init() error {

	workdir, err := filesystem.DetermineHomePath(app.config.WorkDir, app.config.Rootless)
	if err != nil {
		return err
	}

	if err := filesystem.MkDir(workdir, 0755); err != nil {
		return err
	}

	stg, err := storage.Get(storage.BboltDriver, storage.BboltConfig{DbDir: workdir})
	if err != nil {
		return fmt.Errorf("cannot initialize storage: %v", err)
	}

	app.State = state.NewState(context.Background(), stg)
	app.HttpServer = server.NewServer(app.State, stg, app.config)

	return nil
}
