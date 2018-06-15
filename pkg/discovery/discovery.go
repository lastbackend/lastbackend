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
	"github.com/lastbackend/lastbackend/pkg/discovery/cache"
	"github.com/lastbackend/lastbackend/pkg/discovery/envs"
	"github.com/lastbackend/lastbackend/pkg/discovery/runtime"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
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

	log.Info("Start service discovery")

	stg, err := storage.Get(viper.GetString("storage"))
	if err != nil {
		log.Fatalf("Can not initialize storage: %v", err)
	}
	envs.Get().SetStorage(stg)
	envs.Get().SetCache(cache.New())

	r := runtime.NewRuntime(context.Background())
	r.Restore()
	r.Loop()

	sd, err := Listen(viper.GetInt("discovery.port"))
	if err != nil {
		log.Fatalf("Start discovery server error: %v", err)
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

	sd.Shutdown()

	log.Info("Handle SIGINT and SIGTERM.")
	return true
}
