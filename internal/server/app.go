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

	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/server/config"
	"github.com/lastbackend/lastbackend/internal/server/server"
	"github.com/lastbackend/lastbackend/internal/server/state"
	"github.com/lastbackend/lastbackend/tools/logger"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	config  config.Config
	storage storage.IStorage

	Server *server.Server
	State  *state.State
}

func New(stg storage.IStorage, config config.Config) (*App, error) {

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
	app.storage = stg

	if err := app.init(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app App) Run(ctx context.Context) error {
	log := logger.WithContext(ctx)

	log.Infof("Run server")

	go func() {
		if err := app.Server.Run(ctx); err != nil {
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
	var ctx = context.Background()
	app.State = state.NewState(ctx, app.storage)
	app.Server =server.NewServer(app.config)
	return nil
}
