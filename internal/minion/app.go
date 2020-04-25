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
	"fmt"
	"path/filepath"

	"github.com/lastbackend/lastbackend/internal/minion/containerd"
	"github.com/lastbackend/lastbackend/internal/minion/controller"
	"github.com/lastbackend/lastbackend/internal/minion/exporter"
	"github.com/lastbackend/lastbackend/internal/minion/rootless"
	"github.com/lastbackend/lastbackend/internal/minion/runtime"
	"github.com/lastbackend/lastbackend/internal/minion/server"
	"github.com/lastbackend/lastbackend/internal/minion/state"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/filesystem"
	client "github.com/lastbackend/lastbackend/pkg/client/cluster"
	"github.com/lastbackend/lastbackend/pkg/network"
	"github.com/lastbackend/lastbackend/pkg/runtime/cii"
	"github.com/lastbackend/lastbackend/pkg/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/runtime/csi"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/spf13/viper"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	v          *viper.Viper
	HttpServer *server.HttpServer
	State      *state.State
	Runtime    *runtime.Runtime
	Containerd *containerd.Containerd
	Controller *controller.Controller
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
	app.ctx, app.cancel = context.WithCancel(context.Background())
	app.v = v

	if err := app.init(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run() error {
	log := logger.WithContext(context.Background())
	log.Infof("Run minion service")

	if err := app.Containerd.Run(); err != nil {
		log.Errorf("Run containerd server err: %v", err)
		return err
	}

	if err := app.Runtime.Run(); err != nil {
		log.Errorf("Run runtime err: %v", err)
		return err
	}

	if app.Controller != nil {
		if err := app.Controller.Connect(app.v); err != nil {
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
	app.Containerd.Stop()
	app.cancel()
}

func (app *App) init() error {

	var err error

	log := logger.WithContext(context.Background())
	log.Infof("Init minion service")

	workdir, err := filesystem.DetermineHomePath(app.v.GetString("workdir"), app.v.GetBool("rootless"))
	if err != nil {
		return err
	}

	if err := filesystem.MkDir(workdir, 0755); err != nil {
		return err
	}
fmt.Println(">>>>", app.v.GetBool("rootless"))
	if app.v.GetBool("rootless") {
		if err := rootless.Rootless(workdir); err != nil {
			return err
		}
	}

	copts := new(containerd.Config)
	copts.Registry = app.v.GetString("registry.config")
	copts.ConfigPath = filepath.Join(workdir, "etc/containerd/config.toml")
	copts.Root = filepath.Join(workdir, "containerd")
	copts.Opt = filepath.Join(workdir, "containerd")
	copts.State = "/run/lastbackend/containerd"
	copts.Address = filepath.Join(copts.State, "containerd.sock")
	copts.Template = filepath.Join(workdir, "etc/containerd/config.toml.tmpl")
	if !app.v.GetBool("debug") {
		copts.Log = filepath.Join(workdir, "containerd/containerd.log")
	}

	app.Containerd, err = containerd.New(copts)
	if err != nil {
		return fmt.Errorf("Cannot initialize containerd: %v", err)
	}

	app.v.Set("runtime.iri.type", cii.ContainerdDriver)
	app.v.Set("runtime.iri.containerd.address", copts.Address)

	_cii, err := cii.New(app.v)
	if err != nil {
		return fmt.Errorf("Cannot initialize iri: %v", err)
	}

	app.v.Set("runtime.cri.type", cri.ContainerdDriver)
	app.v.Set("runtime.cri.containerd.address", copts.Address)

	_cri, err := cri.New(app.v)
	if err != nil {
		return fmt.Errorf("Cannot initialize cri: %v", err)
	}

	_csi := make(map[string]csi.CSI, 0)

	csis := app.v.GetStringMap("container.csi")
	if csis != nil {
		for kind := range csis {
			si, err := csi.New(kind, app.v)
			if err != nil {
				log.Errorf("Cannot initialize [%s] csi: %v", kind, err)
				return err
			}
			_csi[kind] = si
		}
	}

	_net, err := network.New(app.v)
	if err != nil {
		return fmt.Errorf("Can not initialize network: %v", err)
	}

	_state := state.New()

	// TODO: Need implement logic
	//_state.Node().Info = runtime.NodeInfo()
	//_state.Node().Status = runtime.NodeStatus()

	_exp, err := exporter.NewExporter(_state.Node().Info.Hostname, models.EmptyString)
	if err != nil {
		return fmt.Errorf("Can not initialize collector: %v", err)
	}

	app.Runtime, err = runtime.New(_cri, _cii, _csi, _net, _state, _exp, app.v.GetString("workdir"), app.v.GetString("manifest_dir"))
	if err != nil {
		return fmt.Errorf("Can not initialize runtime: %v", err)
	}

	if app.v.IsSet("api.uri") && len(app.v.GetString("api.uri")) != 0 {

		cfg := client.NewConfig()
		cfg.BearerToken = app.v.GetString("token")

		if app.v.IsSet("api.tls") && app.v.GetBool("api.tls.verify") {
			cfg.TLS = client.NewTLSConfig()
			cfg.TLS.Verify = app.v.GetBool("api.tls.verify")
			cfg.TLS.CertFile = app.v.GetString("api.tls.cert")
			cfg.TLS.KeyFile = app.v.GetString("api.tls.key")
			cfg.TLS.CAFile = app.v.GetString("api.tls.ca")
		}

		endpoint := app.v.GetString("api.uri")

		rest, err := client.New(client.ClientHTTP, endpoint, cfg)
		if err != nil {
			return fmt.Errorf("Can not initialize http client: %v", err)
		}

		app.Controller = controller.New(app.Runtime, rest, _net, _state)
	}

	return nil
}
