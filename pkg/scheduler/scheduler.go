//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package scheduler

import (
	"github.com/lastbackend/lastbackend/pkg/scheduler/runtime"
	"os/signal"
	"syscall"

	"context"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/scheduler/envs"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
	"os"
)

const app = "scheduler"

func Daemon() bool {

	var (
		env  = envs.Get()
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log.Info("Start State Scheduler")

	stg, err := storage.Get(viper.GetString("etcd"))
	if err != nil {
		log.Fatalf("Cannot initialize storage: %v", err)
	}
	env.SetStorage(stg)

	// Initialize Runtime
	r := runtime.NewRuntime(context.Background())
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

	return true
}
