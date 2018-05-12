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

package controller

import (
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/log"
	"os/signal"
	"syscall"

	"context"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
	"os"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/ipam"
)

// Daemon - controller entrypoint
func Daemon() bool {

	var (
		env  = envs.Get()
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log.Info("Start Stage Controller")

	stg, err := storage.Get(viper.GetString("etcd"))
	if err != nil {
		log.Fatalf("Cannot initialize storage: %s", err.Error())
	}
	env.SetStorage(stg)

	ipm, err := ipam.New(viper.GetString("service.cidr"))
	if err != nil {
		log.Fatalf("Cannot initialize ipam service: %s", err.Error())
	}
	env.SetIPAM(ipm)

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
