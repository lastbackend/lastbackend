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

package discovery

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/discovery/cache"
	"github.com/lastbackend/lastbackend/internal/discovery/controller"
	"github.com/lastbackend/lastbackend/internal/discovery/envs"
	"github.com/lastbackend/lastbackend/internal/discovery/runtime"
	"github.com/lastbackend/lastbackend/internal/discovery/state"
	"os"
	"os/signal"
	"syscall"

	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/client"
	l "github.com/lastbackend/lastbackend/tools/log"
	"github.com/spf13/viper"
)

func Daemon(v *viper.Viper) bool {

	var (
		env  = envs.Get()
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log := l.New(v.GetInt("verbose"))

	log.Info("Start service discovery")

	port := uint(53)
	if v.IsSet("dns.port") {
		port = uint(v.GetInt("dns.port"))
	}

	ro := &runtime.RuntimeOpts{
		Iface: v.GetString("runtime.interface"),
		Port:  uint16(port),
	}
	r := runtime.New(ro)

	st := state.New()
	env.SetState(st)
	st.Discovery().Info = r.DiscoveryInfo()
	st.Discovery().Status = r.DiscoveryStatus()

	stg, err := storage.Get(v)
	if err != nil {
		log.Fatalf("Cannot initialize storage: %s", err.Error())
	}
	env.SetStorage(stg)
	env.SetCache(cache.New(v.GetDuration("dns.ttl")))

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

		rest, err := client.New(client.ClientHTTP, endpoint, cfg)
		if err != nil {
			log.Errorf("Init client err: %s", err)
		}

		c := rest.V1().Cluster().Discovery(st.Discovery().Info.Hostname)
		env.SetClient(c)

		ctl := controller.New(r)

		if err := ctl.Connect(context.Background()); err != nil {
			log.Errorf("ingress:initialize: connect err %s", err.Error())
		}

		go ctl.Sync(context.Background())
	}

	sd, err := Listen(v.GetString("dns.host"), v.GetInt("dns.port"))
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
