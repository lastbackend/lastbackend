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

package discovery

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/network"
	"os"
	"os/signal"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/api/client"
	"github.com/lastbackend/lastbackend/pkg/discovery/cache"
	"github.com/lastbackend/lastbackend/pkg/discovery/controller"
	"github.com/lastbackend/lastbackend/pkg/discovery/envs"
	"github.com/lastbackend/lastbackend/pkg/discovery/runtime"
	"github.com/lastbackend/lastbackend/pkg/discovery/state"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
)

func Daemon() bool {

	var (
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log.New(viper.GetInt("verbose"))
	log.Info("Start service discovery")

	net, err := network.New()
	if err != nil {
		log.Errorf("can not initialize network: %s", err.Error())
		os.Exit(1)
	}
	envs.Get().SetNet(net)

	st := state.New()
	envs.Get().SetState(st)
	st.Discovery().Info = runtime.DiscoveryInfo()
	st.Discovery().Status = runtime.DiscoveryStatus()

	stg, err := storage.Get(viper.GetString("etcd"))
	if err != nil {
		log.Fatalf("Cannot initialize storage: %v", err)
	}
	envs.Get().SetStorage(stg)
	envs.Get().SetCache(cache.New())

	r := runtime.NewRuntime(context.Background())
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

		c := rest.V1().Cluster().Discovery(st.Discovery().Info.Hostname)
		envs.Get().SetClient(c)
		ctl := controller.New(r)

		if err := ctl.Connect(context.Background()); err != nil {
			log.Errorf("ingress:initialize: connect err %s", err.Error())
		}

		go ctl.Sync(context.Background())
	}

	r.Restore(context.Background())
	r.Loop(context.Background())

	sd, err := Listen(viper.GetInt("dns.port"))
	if err != nil {
		log.Fatalf("Start discovery server error: %v", err)
	}

	st.Discovery().Status.Ready = true

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

	sd.Shutdown()

	log.Info("Handle SIGINT and SIGTERM.")
	return true
}
