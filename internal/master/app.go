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

	"github.com/lastbackend/lastbackend/internal/master/cache"
	"github.com/lastbackend/lastbackend/internal/master/ipam"
	"github.com/lastbackend/lastbackend/internal/master/runtime"
	"github.com/lastbackend/lastbackend/internal/master/server"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/spf13/viper"
)

const defaultCIDR = "172.0.0.0/24"

type App struct {
	v          *viper.Viper
	HttpServer *server.HttpServer
	Runtime    *runtime.Runtime
	IPAM       ipam.IPAM
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

	if !app.v.IsSet("storage") {
		return errors.New("storage not configured")
	}

	stg, err := storage.Get(app.v)
	if err != nil {
		return fmt.Errorf("cannot initialize storage: %v", err)
	}

	app.HttpServer = server.NewServer(stg, app.v)
	app.Runtime = runtime.NewRuntime(context.Background(), stg, app.IPAM, cache.NewCache())
	app.Runtime.Loop()

	cidr := defaultCIDR
	if app.v.IsSet("service") && app.v.IsSet("service.cidr") {
		cidr = app.v.GetString("service.cidr")
	}

	ipm, err := ipam.New(stg, cidr)
	if err != nil {
		return fmt.Errorf("cannot initialize ipam: %v", err)
	}

	app.IPAM = ipm

	return nil
}
