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

package ingress

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/lastbackend/lastbackend/pkg/api/client"
	"github.com/lastbackend/lastbackend/pkg/ingress/controller"
	"github.com/lastbackend/lastbackend/pkg/ingress/envs"
	"github.com/lastbackend/lastbackend/pkg/ingress/runtime"
	"github.com/lastbackend/lastbackend/pkg/ingress/state"
	l "github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/network"
	"github.com/spf13/viper"
)

func Daemon(v *viper.Viper) bool {

	var (
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log := l.New(v.GetInt("verbose"))
	log.Info("Start Ingress server")

	if !v.IsSet("haproxy") {
		log.Fatalf("Haproxy not configured")
	}

	if v.IsSet("network") {
		net, err := network.New(v)
		if err != nil {
			log.Errorf("can not initialize network: %s", err.Error())
			os.Exit(1)
		}
		envs.Get().SetNet(net)
	}

	st := state.New()

	envs.Get().SetState(st)
	envs.Get().SetTemplate(template.Must(template.New("").Parse(runtime.HaproxyTemplate)),
		v.GetString("haproxy.config"),
		v.GetString("haproxy.pid"))

	envs.Get().SetHaproxy(v.GetString("haproxy.exec"))

	conf := runtime.NewHAProxyConfig(
		uint16(v.GetInt("haproxy.port")),
		v.GetString("haproxy.username"),
		v.GetString("haproxy.password"),
	)
	iface := v.GetString("network.interface")
	r := runtime.New(iface, conf)

	st.Ingress().Info = r.IngressInfo()
	st.Ingress().Status = r.IngressStatus()

	go func() {
		r.Restore()
		r.Loop()
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
		rest, err := client.New(client.ClientHTTP, endpoint, cfg)
		if err != nil {
			log.Errorf("Init client err: %s", err)
		}

		c := rest.V1().Cluster().Ingress(st.Ingress().Info.Hostname)
		envs.Get().SetClient(c)

		ctl := controller.New(r)
		if err := ctl.Connect(context.Background()); err != nil {
			log.Errorf("ingress:initialize: connect err %s", err.Error())
		}

		go ctl.Subscribe()
		go ctl.Sync()
	}

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
