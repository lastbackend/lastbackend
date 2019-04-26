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

package api

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/api/cache"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http"
	"github.com/lastbackend/lastbackend/pkg/api/runtime"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	l "github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/monitor"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
)

func Daemon(v *viper.Viper) {

	var (
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log := l.New(v.GetInt("verbose"))

	log.Info("Start API server")

	if !v.IsSet("storage") {
		log.Fatalf("Storage not configured")
	}

	stg, err := storage.Get(v)
	if err != nil {
		log.Fatalf("Cannot initialize storage: %s", err.Error())
	}
	envs.Get().SetStorage(stg)

	envs.Get().SetCache(cache.NewCache())

	mnt := monitor.New()
	envs.Get().SetMonitor(mnt)

	go func() {
		if err := mnt.Watch(context.Background(), stg, nil); err != nil {
			log.Fatalf("Cannot initialize monitor: %s", err.Error())
		}
	}()

	runtime.New().Run()

	if v.IsSet("vault") {
		vault := &types.Vault{
			Endpoint: v.GetString("vault.endpoint"),
			Token:    v.GetString("vault.token"),
		}
		envs.Get().SetVault(vault)
	}

	go func() {

		opts := new(http.HttpOpts)
		opts.BearerToken = v.GetString("token")
		if v.IsSet("server.tls") {
			opts.Insecure = v.GetBool("server.tls.insecure")
			opts.CertFile = v.GetString("server.tls.cert")
			opts.KeyFile = v.GetString("server.tls.key")
			opts.CaFile = v.GetString("server.tls.ca")
		} else {
			opts.Insecure = true
		}

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
}