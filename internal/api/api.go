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
	"os"
	"os/signal"
	"syscall"

	"github.com/lastbackend/lastbackend/internal/api/cache"
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/api/http"
	"github.com/lastbackend/lastbackend/internal/api/runtime"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	l "github.com/lastbackend/lastbackend/tools/log"
	"github.com/spf13/viper"
)

func Daemon(v *viper.Viper) bool {

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
	envs.Get().SetClusterInfo(v.GetString("name"), v.GetString("description"))
	envs.Get().SetDomain(v.GetString("domain.internal"), v.GetString("domain.external"))
	envs.Get().SetAccessToken(v.GetString("token"))

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

	return true
}
