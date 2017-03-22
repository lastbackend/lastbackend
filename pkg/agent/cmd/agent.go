package cmd

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/pkg/agent/config"
	"os"
	"os/signal"
	"syscall"
	"github.com/lastbackend/lastbackend/pkg/agent/http"
	"github.com/Sirupsen/logrus"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime"
)

func Agent(cmd *cli.Cmd) {

	var cfg = config.Get()

	cmd.Spec = "[-d]"

	var debug = cmd.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "Enable debug mode"})

	cmd.Before = func() {

		if *debug {
			cfg.Debug = *debug
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("Logger debug mode enabled")
		}

		if cfg.HttpServer.Port == 0 {
			cfg.HttpServer.Port = 2967
		}
	}

	cmd.Action = func() {

		var (
			sigs   = make(chan os.Signal)
			done   = make(chan bool, 1)
			routes = http.NewRouter()
		)

		go http.RunHttpServer(routes, cfg.HttpServer.Port)

		runtime := runtime.New()
		runtime.Init()

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

		logrus.Info("Handle SIGINT and SIGTERM.")
	}

}