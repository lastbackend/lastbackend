//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package exporter

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/api/client"
	"github.com/lastbackend/lastbackend/pkg/exporter/controller"
	"github.com/lastbackend/lastbackend/pkg/exporter/envs"
	"github.com/lastbackend/lastbackend/pkg/exporter/http"
	"github.com/lastbackend/lastbackend/pkg/exporter/logger"
	"github.com/lastbackend/lastbackend/pkg/exporter/runtime"
	"github.com/lastbackend/lastbackend/pkg/exporter/state"
	l "github.com/lastbackend/lastbackend/pkg/log"
	"github.com/spf13/viper"
)

const (
	defaultIface = "eth0"
)

func Daemon(v *viper.Viper) bool {

	var (
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log := l.New(v.GetInt("verbose"))
	log.Info("Start Exporter server")

	iface := defaultIface
	if v.IsSet("network") {
		iface = v.GetString("network.interface")
	}

	ro := &runtime.RuntimeOpts{
		Port:  uint16(v.GetInt("server.port")),
		Iface: iface,
	}

	if v.IsSet("logger") {
		ro.Logger = &logger.LoggerOpts{
			Host:    v.GetString("logger.host"),
			Port:    uint16(v.GetInt("logger.port")),
			Workdir: v.GetString("logger.workdir"),
		}
	}

	r, err := runtime.New(ro)
	if err != nil {
		log.Errorf("can not start runtime: %s", err.Error())
		os.Exit(1)
	}

	st := state.New()
	st.Exporter().Info = r.ExporterInfo()
	st.Exporter().Status = r.ExporterStatus()

	envs.Get().SetState(st)

	go func() {
		if err := r.Start(); err != nil {
			log.Errorf("can not start runtime listener: %s", err.Error())
			os.Exit(1)
		}
	}()

	if v.IsSet("api") {

		cfg := client.NewConfig()
		cfg.BearerToken = v.GetString("token")

		if v.IsSet("api.tls") && !v.GetBool("api.tls.insecure") {
			cfg.TLS = client.NewTLSConfig()
			cfg.TLS.CertFile = v.GetString("api.tls.cert")
			cfg.TLS.KeyFile = v.GetString("api.tls.key")
			cfg.TLS.CAFile = v.GetString("api.tls.ca")
		}

		endpoint := v.GetString("api.uri")
		if len(endpoint) == 0 {
			log.Fatalf("Api endpoint is not set")
		}

		rest, err := client.New(client.ClientHTTP, endpoint, cfg)
		if err != nil {
			log.Errorf("Init client err: %s", err)
		}

		c := rest.V1().Cluster().Exporter(st.Exporter().Info.Hostname)
		envs.Get().SetClient(c)
		ctl := controller.New(r)

		if err := ctl.Connect(context.Background()); err != nil {
			log.Errorf("ingress:initialize: connect err %s", err.Error())
		}

		go ctl.Sync()
	}

	go func() {

		opts := new(http.HttpOpts)
		if v.IsSet("server.tls") {
			opts.Insecure = v.GetBool("server.tls.insecure")
			opts.CertFile = v.GetString("server.tls.cert")
			opts.KeyFile = v.GetString("server.tls.key")
			opts.CaFile = v.GetString("server.tls.ca")
		} else {
			opts.Insecure = true
		}

		envs.Get().SetAccessToken(v.GetString("token"))

		host := v.GetString("server.host")
		port := v.GetInt("server.port")

		if err := http.Listen(host, port, opts); err != nil {
			log.Fatalf("Http server start error: %v", err)
		}

	}()

	// Handle SIGINT and SIGTERM.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-sigs:
				done <- true
				return
			}
		}
	}()

	<-done

	log.Info("Handle SIGINT and SIGTERM.")

	return true
}
