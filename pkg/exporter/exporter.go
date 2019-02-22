//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/pkg/api/client"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/exporter/controller"
	"github.com/lastbackend/lastbackend/pkg/exporter/envs"
	"github.com/lastbackend/lastbackend/pkg/exporter/http"
	"github.com/lastbackend/lastbackend/pkg/exporter/runtime"
	"github.com/lastbackend/lastbackend/pkg/exporter/state"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/network"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func Daemon() bool {

	var (
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log.New(viper.GetInt("verbose"))
	log.Info("Start Exporter server")

	net, err := network.New()
	if err != nil {
		log.Errorf("can not initialize network: %s", err.Error())
		os.Exit(1)
	}
	envs.Get().SetNet(net)

	st := state.New()

	st.Exporter().Info = runtime.ExporterInfo()
	st.Exporter().Status = runtime.ExporterStatus()

	envs.Get().SetState(st)

	r, err := runtime.NewRuntime()
	if err != nil {
		log.Errorf("can not start runtime: %s", err.Error())
		os.Exit(1)
	}

	go func() {
		if err := r.Logger(context.Background()); err != nil {
			log.Errorf("can not start logger listener: %s", err.Error())
			os.Exit(1)
		}
	}()

	types.SecretAccessToken = viper.GetString("token")

	if viper.IsSet("api") || viper.IsSet("api_uri") {

		cfg := client.NewConfig()
		cfg.BearerToken = viper.GetString("token")

		if viper.IsSet("api.tls") && !viper.GetBool("api.tls.insecure") {
			cfg.TLS = client.NewTLSConfig()
			cfg.TLS.CertFile = viper.GetString("api.tls.cert")
			cfg.TLS.KeyFile = viper.GetString("api.tls.key")
			cfg.TLS.CAFile = viper.GetString("api.tls.ca")
		}

		endpoint := viper.GetString("api.uri")
		if viper.IsSet("api_uri") {
			endpoint = viper.GetString("api_uri")
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

		go ctl.Sync(context.Background())
	}

	go func() {
		opts := new(http.HttpOpts)
		opts.Insecure = viper.GetBool("exporter.http.tls.insecure")
		opts.CertFile = viper.GetString("exporter.http.tls.cert")
		opts.KeyFile = viper.GetString("exporter.http.tls.key")
		opts.CaFile = viper.GetString("exporter.http.tls.ca")

		types.SecretAccessToken = viper.GetString("token")

		if err := http.Listen(viper.GetString("exporter.http.host"), viper.GetInt("exporter.http.port"), opts); err != nil {
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
