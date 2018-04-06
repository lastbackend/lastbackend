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

package ingress

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/ingress/runtime"
	"github.com/lastbackend/lastbackend/pkg/ingress/state"

	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/client"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/ingress/envs"
	"github.com/spf13/viper"
	"github.com/lastbackend/lastbackend/pkg/ingress/events/exporter"
	"github.com/lastbackend/lastbackend/pkg/ingress/events"
)

// Daemon - run node daemon
func Daemon() {

	var (
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log.New(viper.GetInt("verbose"))
	log.Info("Start Ingress")

	state := state.New()

	envs.Get().SetState(state)

	r := runtime.NewRuntime(context.Background())

	host := viper.GetString("api.uri")
	port := viper.GetInt("api.port")
	tls := viper.GetBool("api.tls")

	schema := "http"
	if tls {
		schema = "https"
	}

	endpoint := fmt.Sprintf("%s://%s:%d", schema, host, port)
	types.SecretAccessToken = viper.GetString("token")

	rest, err := client.NewHTTP(endpoint, &client.Config{
		BearerToken: types.SecretAccessToken,
		Timeout:     5,
	})

	if err != nil {
		log.Errorf("node:initialize client err: %s", err.Error())
		os.Exit(0)
	}

	c := rest.V1().Cluster().Ingress(viper.GetString("ingress.name"))

	envs.Get().SetClient(c)
	e := exporter.NewExporter()
	e.SetDispatcher(events.Dispatcher)
	envs.Get().SetExporter(e)

	if err := r.Connect(context.Background()); err != nil {
		log.Fatalf("node:initialize: connect err %s", err.Error())
	}

	e.Loop()
	r.Loop()

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

	return
}