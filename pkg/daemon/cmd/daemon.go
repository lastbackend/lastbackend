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

package cmd

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/pkg/daemon/api"
	"github.com/lastbackend/lastbackend/pkg/daemon/config"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"os"
	"os/signal"
	"syscall"
)

func Daemon(cmd *cli.Cmd) {
	var ctx = context.Get()
	var cfg = config.Get()

	cmd.Spec = "[-c][-d]"

	var debug = cmd.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "Enable debug mode"})
	var configPath = cmd.String(cli.StringOpt{Name: "c config", Value: "./config.yaml", Desc: "Path to config file", HideValue: true})

	cmd.Before = func() {
		if *configPath != "" {
			if err := cfg.Configure(*configPath); err != nil {
				panic(err)
			}
		}

		cfg.Debug = *debug
		ctx.Init(cfg)
	}

	cmd.Action = func() {

		var (
			sigs = make(chan os.Signal)
			done = make(chan bool, 1)
		)

		go func() {
			if err := api.Listen(cfg.HttpServer.Host, cfg.HttpServer.Port); err != nil {
				ctx.Log.Warnf("Http server start error: %s", err.Error())
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

		ctx.Log.Info("Handle SIGINT and SIGTERM.")
	}
}
