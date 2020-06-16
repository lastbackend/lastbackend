// +build linux
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

package daemon

import (
	"context"
	"fmt"
	"path"

	"github.com/lastbackend/lastbackend/internal/agent"
	"github.com/lastbackend/lastbackend/internal/daemon/config"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/server"
	"github.com/lastbackend/lastbackend/internal/util/filesystem"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/pkg/errors"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	config config.Config

	agent  *agent.App
	server *server.App
}

func New(config config.Config) (*App, error) {

	loggerOpts := logger.Configuration{}
	loggerOpts.EnableConsole = true

	if config.DisableServer && config.DisableSchedule {
		return nil, errors.New("cannot use 'no-schedule' and 'agent' flags together")
	}

	if config.Debug {
		loggerOpts.ConsoleLevel = logger.Debug
	}

	if err := logger.NewLogger(loggerOpts); err != nil {
		return nil, errors.New("logger initialize failed")
	}

	workdir, err := filesystem.DetermineHomePath(config.RootDir, false)
	if err != nil {
		return nil, err
	}

	if err := filesystem.MkDir(workdir, 0755); err != nil {
		return nil, err
	}

	app := new(App)
	app.ctx, app.cancel = context.WithCancel(context.Background())
	config.RootDir = workdir
	app.config = config

	return app, nil
}

func (app App) Run() error {
	log := logger.WithContext(context.Background())
	log.Infof("Run daemon process")

	if err := filesystem.MkDir(path.Join(app.config.RootDir, "data"), 0755); err != nil {
		return err
	}

	stg, err := storage.Get(storage.BboltDriver, storage.BboltConfig{DbDir: path.Join(app.config.RootDir, "data"), DbName: "state"})
	if err != nil {
		return fmt.Errorf("cannot initialize storage: %v", err)
	}

	if !app.config.DisableSchedule {
		app.agent, err = agent.New(stg, app.config.GetAgentConfig())
		if err != nil {
			return err
		}
		if err := app.agent.Run(); err != nil {
			return err
		}
	}

	if !app.config.DisableServer {
		app.server, err = server.New(stg, app.config.GetServerConfig())
		if err != nil {
			return err
		}
		if err := app.server.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) Stop() {
	log := logger.WithContext(context.Background())
	log.Infof("Stop daemon process")

	if !app.config.DisableSchedule {
		app.agent.Stop()
	}
	if !app.config.DisableServer {
		app.server.Stop()
	}

	app.cancel()
}
