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
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/daemon/config"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/daemon/api"
)

func Daemon(cmd *cli.Cmd) {
	var err error

	var ctx = context.Get()
	var cfg = config.Get()

	cmd.Spec = "[-c][-d]"

	var debug = cmd.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "Enable debug mode"})
	var configPath = cmd.String(cli.StringOpt{Name: "c config", Value: "./config.yaml", Desc: "Path to config file", HideValue: true})

	cmd.Before = func() {

		ctx.Log = logger.Init()

		// If you want to connect to local syslog (Ex. "/dev/log" or "/var/run/syslog" or "/var/run/log").
		// Just assign empty string to the first two parameters of NewSyslogHook. It should look like the following.
		ctx.Log.SetSyslog("", "", syslog.LOG_INFO, "")

		if *configPath != "" {
			if err := cfg.Configure(*configPath); err != nil {
				ctx.Log.Panic(err)
			}
		}

		if *debug {
			cfg.Debug = *debug
			ctx.Log.SetDebugLevel()
			ctx.Log.Info("Logger debug mode enabled")
		}

		// Initializing database
		ctx.Log.Info("Initializing daemon")

		ctx.Storage, err = storage.Get(cfg.GetEtcdDB())
		if err != nil {
			ctx.Log.Panic(err)
		}

		if cfg.HttpServer.Port == 0 {
			cfg.HttpServer.Port = 3000
		}

	}

	cmd.Action = func() {

		var (
			sigs   = make(chan os.Signal)
			done   = make(chan bool, 1)
		)

		go func () {
			if err := api.Listen(cfg.HttpServer.Port); err != nil {
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
