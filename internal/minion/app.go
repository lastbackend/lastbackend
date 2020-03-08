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
	"log"

	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/spf13/viper"
)

type App struct {
	v *viper.Viper
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
}

func (app *App) init() error {

	if !app.v.IsSet("container") {
		log.Fatalf("Container not configured")
	}
	if !app.v.IsSet("container.cri") {
		log.Fatalf("CRI not configured")
	}
	if !app.v.IsSet("container.iri") {
		log.Fatalf("IRI not configured")
	}
	if !app.v.IsSet("container.csi") {
		log.Fatalf("CSI not configured")
	}

	return nil
}

//
//// Daemon - run node daemon
//func Daemon(v *viper.Viper) {
//
//	var (
//		sigs = make(chan os.Signal)
//		done = make(chan bool, 1)
//	)
//
//	log := l.New(v.GetInt("verbose"))
//	log.Info("Start Node")
//
//	if !v.IsSet("container") {
//		log.Fatalf("Container not configured")
//	}
//	if !v.IsSet("container.cri") {
//		log.Fatalf("CRI not configured")
//	}
//	if !v.IsSet("container.iri") {
//		log.Fatalf("IRI not configured")
//	}
//	if !v.IsSet("container.csi") {
//		log.Fatalf("CSI not configured")
//	}
//
//	if err := envs.Get().SetConfig(v); err != nil {
//		log.Fatalf("Parse config err: %v", err)
//	}
//
//	_cri, err := cri.New(v)
//	if err != nil {
//		log.Errorf("Cannot initialize cri: %v", err)
//	}
//
//	_cii, err := cii.New(v)
//	if err != nil {
//		log.Errorf("Cannot initialize iri: %v", err)
//	}
//
//	csis := v.GetStringMap("container.csi")
//	if csis != nil {
//		for kind := range csis {
//			si, err := csi.New(kind, v)
//			if err != nil {
//				log.Errorf("Cannot initialize csi: %s > %v", kind, err)
//				return
//			}
//			envs.Get().SetCSI(kind, si)
//		}
//	}
//
//	st := state.New()
//	envs.Get().SetState(st)
//	envs.Get().SetCRI(_cri)
//	envs.Get().SetCII(_cii)
//
//	if v.IsSet("network") {
//		net, err := network.New(v)
//		if err != nil {
//			log.Errorf("can not initialize network: %s", err.Error())
//			os.Exit(1)
//		}
//		envs.Get().SetNet(net)
//	}
//
//	r, err := runtime.New()
//	if err != nil {
//		log.Errorf("can not initialize runtime: %s", err.Error())
//		os.Exit(1)
//	}
//
//	st.Node().Info = runtime.NodeInfo()
//	st.Node().Status = runtime.NodeStatus()
//
//	if err := r.Restore(); err != nil {
//		log.Errorf("restore err: %v", err)
//		os.Exit(1)
//	}
//	r.Subscribe()
//	r.Loop()
//
//	c, err := exporter.NewExporter(st.Node().Info.Hostname, types.EmptyString)
//	if err != nil {
//		log.Errorf("can not initialize collector: %s", err.Error())
//		os.Exit(1)
//	}
//	envs.Get().SetExporter(c)
//	go c.Listen()
//
//	if v.IsSet("manifest.dir") {
//		dir := v.GetString("manifest.dir")
//		if dir != types.EmptyString {
//			r.Provision(dir)
//		}
//	}
//
//	if v.IsSet("api.uri") && len(v.GetString("api.uri")) != 0 {
//
//		cfg := client.NewConfig()
//		cfg.BearerToken = v.GetString("token")
//
//		if v.IsSet("api.tls") {
//			cfg.TLS = client.NewTLSConfig()
//			cfg.TLS.Verify = v.GetBool("api.tls.verify")
//			cfg.TLS.CertFile = v.GetString("api.tls.cert")
//			cfg.TLS.KeyFile = v.GetString("api.tls.key")
//			cfg.TLS.CAFile = v.GetString("api.tls.ca")
//		}
//
//		endpoint := v.GetString("api.uri")
//
//		rest, err := client.New(client.ClientHTTP, endpoint, cfg)
//		if err != nil {
//			log.Errorf("Init client err: %s", err)
//		}
//
//		n := rest.V1().Cluster().Node(st.Node().Info.Hostname)
//		s := rest.V1()
//		envs.Get().SetClient(n, s)
//
//		ctl := controller.New(r)
//
//		if err := ctl.Connect(v); err != nil {
//			log.Errorf("node:initialize: connect err %s", err.Error())
//
//		}
//		go ctl.Subscribe()
//		go ctl.Sync()
//	}
//
//	go func() {
//		opts := new(http.HttpOpts)
//		opts.BearerToken = v.GetString("token")
//
//		if v.IsSet("server.tls") {
//			opts.TLS = v.GetBool("server.tls.verify")
//			opts.CertFile = v.GetString("server.tls.cert")
//			opts.KeyFile = v.GetString("server.tls.key")
//			opts.CaFile = v.GetString("server.tls.ca")
//		}
//
//		if err := http.Listen(v.GetString("server.host"), v.GetInt("server.port"), opts); err != nil {
//			log.Fatalf("Http server start error: %v", err)
//		}
//	}()
//
//	// Handle SIGINT and SIGTERM.
//	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
//
//	go func() {
//		for {
//			select {
//			case <-sigs:
//				r.Stop()
//				done <- true
//				return
//			}
//		}
//	}()
//
//	<-done
//	log.Info("Handle SIGINT and SIGTERM.")
//
//	return
//}
