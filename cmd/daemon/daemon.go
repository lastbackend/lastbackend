package daemon

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/cmd/daemon/config"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	"github.com/lastbackend/lastbackend/libs/adapter/k8s"
	"github.com/lastbackend/lastbackend/libs/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

func Run(cmd *cli.Cmd) {
	var err error

	var ctx = context.Get()
	var cfg = config.Get()

	cmd.Spec = "[-c][-d]"

	var debug = cmd.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "Enable debug mode"})
	var configPath = cmd.String(cli.StringOpt{Name: "c config", Value: "./config.yaml", Desc: "Path to config file", HideValue: true})

	cmd.Before = func() {

		ctx.Log = new(log.Log)
		ctx.Log.Init()

		if *configPath != "" {

			// Parsing config file
			configBytes, err := ioutil.ReadFile(*configPath)
			if err != nil {
				ctx.Log.Panic(err)
			}

			err = yaml.Unmarshal(configBytes, cfg)
			if err != nil {
				ctx.Log.Panic(err)
			}
		}

		if *debug {
			cfg.Debug = *debug
			ctx.Log.SetDebugLevel()
			ctx.Log.Info("Logger debug mode enabled")
		}

		if cfg.HttpServer.Port == 0 {
			cfg.HttpServer.Port = 3000
		}
	}

	cmd.Action = func() {

		ctx.Log.Info("Initializing daemon")
		ctx.K8S, err = k8s.Get(config.GetK8S())
		if err != nil {
			ctx.Log.Panic(err)
		}

		go RunHttpServer(cfg.HttpServer.Port)

		// Handle SIGINT and SIGTERM.
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		ctx.Log.Debug(<-ch)

		ctx.Log.Info("Handle SIGINT and SIGTERM.")
	}
}
