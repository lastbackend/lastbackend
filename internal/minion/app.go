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

package minion

import (
	"context"
	"errors"
	"os"

	"github.com/lastbackend/lastbackend/internal/minion/controller"
	"github.com/lastbackend/lastbackend/internal/minion/exporter"
	"github.com/lastbackend/lastbackend/internal/minion/runtime"
	"github.com/lastbackend/lastbackend/internal/minion/server"
	"github.com/lastbackend/lastbackend/internal/minion/state"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/pkg/client"
	"github.com/lastbackend/lastbackend/pkg/network"
	"github.com/lastbackend/lastbackend/pkg/runtime/cii"
	"github.com/lastbackend/lastbackend/pkg/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/runtime/csi"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/spf13/viper"
)

type App struct {
	v          *viper.Viper
	HttpServer *server.HttpServer
	State      *state.State
	Runtime    *runtime.Runtime
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
	log.Infof("Run minion service")
	return nil
}

func (app *App) Stop() {
	app.Runtime.Stop()
}

func (app *App) init() error {

	log := logger.WithContext(context.Background())
	log.Infof("Init minion service")

	_cri, err := cri.New(app.v)
	if err != nil {
		log.Errorf("Cannot initialize cri: %v", err)
		return err
	}

	_cii, err := cii.New(app.v)
	if err != nil {
		log.Errorf("Cannot initialize iri: %v", err)
		return err
	}

	_csi := make(map[string]csi.CSI, 0)

	csis := app.v.GetStringMap("container.csi")
	if csis != nil {
		for kind := range csis {
			si, err := csi.New(kind, app.v)
			if err != nil {
				log.Errorf("Cannot initialize csi: %s > %v", kind, err)
				return err
			}
			_csi[kind] = si
		}
	}

	_net, err := network.New(app.v)
	if err != nil {
		log.Errorf("can not initialize network: %s", err.Error())
		os.Exit(1)
	}

	_state := state.New()
	// TODO: Need implement logic
	//_state.Node().Info = runtime.NodeInfo()
	//_state.Node().Status = runtime.NodeStatus()

	_exporter, err := exporter.NewExporter(_state.Node().Info.Hostname, types.EmptyString)
	if err != nil {
		log.Errorf("can not initialize collector: %s", err.Error())
		return err
	}

	r, err := runtime.New(_cri, _cii, _csi, _net, _state, _exporter, app.v.GetString("workdir"))
	if err != nil {
		log.Errorf("can not initialize runtime: %s", err.Error())
		return err
	}

	if err := r.Restore(); err != nil {
		log.Errorf("restore err: %v", err)
		return err
	}
	r.Subscribe()
	r.Loop()

	if app.v.IsSet("manifest.dir") {
		dir := app.v.GetString("manifest.dir")
		if dir != types.EmptyString {
			r.Provision(dir)
		}
	}

	app.Runtime = r

	go _exporter.Listen()

	if app.v.IsSet("api.uri") && len(app.v.GetString("api.uri")) != 0 {

		cfg := client.NewConfig()
		cfg.BearerToken = app.v.GetString("token")

		if app.v.IsSet("api.tls") {
			cfg.TLS = client.NewTLSConfig()
			cfg.TLS.Verify = app.v.GetBool("api.tls.verify")
			cfg.TLS.CertFile = app.v.GetString("api.tls.cert")
			cfg.TLS.KeyFile = app.v.GetString("api.tls.key")
			cfg.TLS.CAFile = app.v.GetString("api.tls.ca")
		}

		endpoint := app.v.GetString("api.uri")

		rest, err := client.New(client.ClientHTTP, endpoint, cfg)
		if err != nil {
			log.Errorf("Init client err: %s", err)
		}

		ctl := controller.New(r, rest, _net, _state)

		if err := ctl.Connect(app.v); err != nil {
			log.Errorf("node:initialize: connect err %s", err.Error())

		}
		go ctl.Subscribe()
		go ctl.Sync()
	}

	go func() {
		if err := app.HttpServer.Run(); err != nil {
			log.Fatalf("http rest server start err: %v", err)
			return
		}
	}()

	return nil
}
