package cmd

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/libs/adapter/k8s"
	"github.com/lastbackend/lastbackend/libs/adapter/storage"
	http_client "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/log"
	"github.com/lastbackend/lastbackend/pkg/daemon/config"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/http"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"github.com/lastbackend/lastbackend/pkg/proxy/server"
)

func Daemon(cmd *cli.Cmd) {
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

		// Initializing database
		ctx.Log.Info("Initializing daemon")
		ctx.K8S, err = k8s.Get(config.GetK8S())
		if err != nil {
			ctx.Log.Panic(err)
		}

		ctx.Storage, err = storage.Get()
		if err != nil {
			ctx.Log.Panic(err)
		}

		if cfg.HttpServer.Port == 0 {
			cfg.HttpServer.Port = 3000
		}

		if cfg.TemplateRegistry.Host == "" {
			cfg.TemplateRegistry.Host = "http://localhost:3003"
		}

		ctx.TemplateRegistry = http_client.New(cfg.TemplateRegistry.Host)
	}

	cmd.Action = func() {

		go http.RunHttpServer(http.NewRouter(), cfg.HttpServer.Port)
		go server.StartProxyServer()

		// Handle SIGINT and SIGTERM.
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		ctx.Log.Debug(<-ch)

		ctx.Log.Info("Handle SIGINT and SIGTERM.")
	}
}
