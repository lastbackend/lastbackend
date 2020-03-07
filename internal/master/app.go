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
	"github.com/lastbackend/lastbackend/internal/master/server"

	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/spf13/viper"
)

const defaultCIDR = "172.0.0.0/24"

type App struct {
	HttpServer *server.HttpServer
	v          *viper.Viper
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

	return nil
}

//// Daemon - controller entrypoint
//func Daemon(v *viper.Viper) bool {
//
//	var (
//		env  = envs.Get()
//		sigs = make(chan os.Signal)
//		done = make(chan bool, 1)
//	)
//
//	log := l.New(v.GetInt("verbose"))
//
//	log.Info("Start lastbackend master daemon")
//
//	stg, err := storage.Get(v)
//	if err != nil {
//		log.Fatalf("Cannot initialize storage: %s", err.Error())
//	}
//	env.SetStorage(stg)
//
//	cidr := defaultCIDR
//	if v.IsSet("service") && v.IsSet("service.cidr") {
//		cidr = v.GetString("service.cidr")
//	}
//
//	ipm, err := ipam.New(cidr)
//	if err != nil {
//		log.Fatalf("Cannot initialize ipam service: %s", err.Error())
//	}
//	env.SetIPAM(ipm)
//
//	envs.Get().SetCache(cache.NewCache())
//	envs.Get().SetClusterInfo(v.GetString("name"), v.GetString("description"))
//	envs.Get().SetDomain(v.GetString("domain.internal"), v.GetString("domain.external"))
//	envs.Get().SetAccessToken(v.GetString("token"))
//
//	// Initialize Container
//	r := runtime.NewRuntime(context.Background())
//	r.Loop()
//
//	if v.IsSet("vault") {
//		vault := &types.Vault{
//			Endpoint: v.GetString("vault.endpoint"),
//			Token:    v.GetString("vault.token"),
//		}
//		envs.Get().SetVault(vault)
//	}
//
//	go func() {
//
//		opts := new(http.HttpOpts)
//		opts.BearerToken = v.GetString("token")
//		if v.IsSet("server.tls") {
//			opts.Insecure = v.GetBool("server.tls.insecure")
//			opts.CertFile = v.GetString("server.tls.cert")
//			opts.KeyFile = v.GetString("server.tls.key")
//			opts.CaFile = v.GetString("server.tls.ca")
//		} else {
//			opts.Insecure = true
//		}
//
//		host := v.GetString("server.host")
//		port := v.GetInt("server.port")
//
//		if err := http.Listen(host, port, opts); err != nil {
//			log.Fatalf("Http server start error: %v", err)
//		}
//
//	}()
//
//	// Handle SIGINT and SIGTERM.
//	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
//
//	go func() {
//		for {
//			select {
//			case <-sigs:
//				done <- true
//				return
//			}
//		}
//	}()
//
//	<-done
//
//	log.Info("Handle SIGINT and SIGTERM.")
//
//	return true
//}
