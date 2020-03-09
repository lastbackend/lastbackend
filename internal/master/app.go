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

package master

import (
	"context"
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/internal/master/state"

	"github.com/lastbackend/lastbackend/internal/master/server"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/spf13/viper"
)

const defaultCIDR = "172.0.0.0/24"

type App struct {
	v          *viper.Viper
	HttpServer *server.HttpServer
	State      *state.State
}

func New(v *viper.Viper) (*App, error) {

	loggerOpts := logger.Configuration{}
	loggerOpts.EnableConsole = true
	if v.GetBool("debug") {
		loggerOpts.ConsoleLevel = logger.Debug
	}

	if err := logger.NewLogger(loggerOpts); err != nil {
		return nil, errors.New("logger initialize failed")
	}

	app := new(App)
	app.v = v

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
}

func (app *App) init() error {

	stg, err := storage.Get(app.v)
	if err != nil {
		return fmt.Errorf("cannot initialize storage: %v", err)
	}
	app.State = state.NewState(context.Background(), stg)
	app.HttpServer = server.NewServer(app.State, stg, app.v)

	return nil
}
