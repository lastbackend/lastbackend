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

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/http"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cni/cni"
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

	state := state.New()

	envs.Get().SetState(state)
	envs.Get().SetCri(cri)
	envs.Get().SetCNI(cni)

	r := runtime.NewRuntime(context.Background())
	r.Restore()

	state.Node().Info = node.GetInfo()
	state.Node().Status = node.GetStatus()

	r.Subscribe()
	r.Loop()

	go func() {
		types.SecretAccessToken = viper.GetString("token")
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
