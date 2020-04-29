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
	"github.com/lastbackend/lastbackend/pkg/runtime/csi"
	"path/filepath"

	"github.com/lastbackend/lastbackend/internal/agent/config"
	"github.com/lastbackend/lastbackend/internal/agent/containerd"
	"github.com/lastbackend/lastbackend/internal/agent/controller"
	"github.com/lastbackend/lastbackend/internal/agent/exporter"
	"github.com/lastbackend/lastbackend/internal/agent/rootless"
	"github.com/lastbackend/lastbackend/internal/agent/runtime"
	"github.com/lastbackend/lastbackend/internal/agent/server"
	"github.com/lastbackend/lastbackend/internal/agent/state"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/filesystem"
	"github.com/lastbackend/lastbackend/pkg/network"
	"github.com/lastbackend/lastbackend/pkg/runtime/cii"
	"github.com/lastbackend/lastbackend/pkg/runtime/cri"
	"github.com/lastbackend/lastbackend/tools/logger"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	config config.Config

	HttpServer *server.HttpServer
	State      *state.State
	Runtime    *runtime.Runtime
	Containerd *containerd.Containerd
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

	if err := app.Containerd.Run(); err != nil {
		log.Errorf("Run containerd server err: %v", err)
		return err
	}

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
	app.Containerd.Stop()
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

	if app.config.Rootless {
		if err := rootless.Rootless(workdir); err != nil {
			return err
		}
	}

	copts := new(containerd.Config)
	copts.Registry = app.config.Registry.Config
	copts.ConfigPath = filepath.Join(workdir, "etc/containerd/config.toml")
	copts.Root = filepath.Join(workdir, "containerd")
	copts.Opt = filepath.Join(workdir, "containerd")
	copts.State = "/run/lastbackend/containerd"
	copts.Address = filepath.Join(copts.State, "containerd.sock")
	copts.Template = filepath.Join(workdir, "etc/containerd/config.toml.tmpl")
	if !app.config.Debug {
		copts.Log = filepath.Join(workdir, "containerd/containerd.log")
	}

	app.Containerd, err = containerd.New(copts)
	if err != nil {
		return fmt.Errorf("Cannot initialize containerd: %v", err)
	}

	_cii, err := cii.New(cii.ContainerdDriver, cii.ContainerdConfig{Address: copts.Address})
	if err != nil {
		return fmt.Errorf("Cannot initialize iri: %v", err)
	}

	_cri, err := cri.New(cri.ContainerdDriver, cri.ContainerdConfig{Address: copts.Address})
	if err != nil {
		return fmt.Errorf("Cannot initialize cri: %v", err)
	}

	_csi := make(map[string]csi.CSI, 0)

	// TODO: Implement csi initialization logic
	//csis := app.config.GetStringMap("container.csi")
	//if csis != nil {
	//	for kind := range csis {
	//		si, err := csi.New(kind, dir.Config{RootDir: filepath.Join(app.config.WorkDir, "csi")})
	//		if err != nil {
	//			log.Errorf("Cannot initialize [%s] csi: %v", kind, err)
	//			return err
	//		}
	//		_csi[kind] = si
	//	}
	//}

	_net, err := network.New(app.config)
	if err != nil {
		return fmt.Errorf("Can not initialize network: %v", err)
	}

	_state := state.New()

	// TODO: Implement cluster state logic
	//_state.Node().Info = runtime.NodeInfo()
	//_state.Node().Status = runtime.NodeStatus()

	_exp, err := exporter.NewExporter(_state.Node().Info.Hostname, models.EmptyString)
	if err != nil {
		return fmt.Errorf("Can not initialize collector: %v", err)
	}

	app.Runtime, err = runtime.New(_cri, _cii, _csi, _net, _state, _exp, app.config)
	if err != nil {
		return fmt.Errorf("Can not initialize runtime: %v", err)
	}

	// TODO: Implement controller initialization logic
	//if app.config.IsSet("api.uri") && len(app.config.GetString("api.uri")) != 0 {
	//
	//	cfg := client.NewConfig()
	//	cfg.BearerToken = app.config.Security.Token
	//
	//	if app.config.IsSet("api.tls") && app.config.GetBool("api.tls.verify") {
	//		cfg.TLS = client.NewTLSConfig()
	//		cfg.TLS.Verify = app.config.GetBool("api.tls.verify")
	//		cfg.TLS.CertFile = app.config.GetString("api.tls.cert")
	//		cfg.TLS.KeyFile = app.config.GetString("api.tls.key")
	//		cfg.TLS.CAFile = app.config.GetString("api.tls.ca")
	//	}
	//
	//	endpoint := app.config.GetString("api.uri")
	//
	//	rest, err := client.New(client.ClientHTTP, endpoint, cfg)
	//	if err != nil {
	//		return fmt.Errorf("Can not initialize http client: %v", err)
	//	}
	//
	//	app.Controller = controller.New(app.Runtime, rest, _net, _state)
	//}

	return nil
}
