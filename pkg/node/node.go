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

package node

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/node/runtime"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/node"
	"github.com/lastbackend/lastbackend/pkg/node/state"

	"fmt"

	"github.com/lastbackend/lastbackend/pkg/api/client"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/events"
	"github.com/lastbackend/lastbackend/pkg/node/events/exporter"
	"github.com/lastbackend/lastbackend/pkg/node/http"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cni/cni"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cpi/cpi"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cri/cri"
	"github.com/spf13/viper"
)

// Daemon - run node daemon
func Daemon() {

	var (
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log.New(viper.GetInt("verbose"))
	log.Info("Start Node")

	cri, err := cri.New()
	if err != nil {
		log.Errorf("Cannot initialize cri: %v", err)
	}

	cni, err := cni.New()
	if err != nil {
		log.Errorf("Cannot initialize cni: %v", err)
	}

	cpi, err := cpi.New()
	if err != nil {
		log.Errorf("Cannot initialize cni: %v", err)
	}

	state := state.New()

	envs.Get().SetState(state)
	envs.Get().SetCRI(cri)
	envs.Get().SetCNI(cni)
	envs.Get().SetCPI(cpi)

	r := runtime.NewRuntime(context.Background())
	r.Restore()

	state.Node().Info = node.GetInfo()
	state.Node().Status = node.GetStatus()

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

	c := rest.V1().Cluster().Node(state.Node().Info.Hostname)

	envs.Get().SetClient(c)
	e := exporter.NewExporter()
	e.SetDispatcher(events.Dispatcher)
	envs.Get().SetExporter(e)

	if err := r.Connect(context.Background()); err != nil {
		log.Fatalf("node:initialize: connect err %s", err.Error())
	}

	r.Subscribe()

	e.Loop()
	r.Loop()

	go func() {

		if err := http.Listen(viper.GetString("node.host"), viper.GetInt("node.port")); err != nil {
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

	return
}
