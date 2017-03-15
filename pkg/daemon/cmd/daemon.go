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
	"github.com/lastbackend/lastbackend/pkg/proxy/server"
	"os"
	"os/signal"
	"syscall"
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

		if cfg.ProxyServer.Port == 0 {
			cfg.ProxyServer.Port = 9999
		}

		if cfg.TemplateRegistry.Host == "" {
			cfg.TemplateRegistry.Host = "http://localhost:3003"
		}

		ctx.TemplateRegistry = http_client.New(cfg.TemplateRegistry.Host)
	}

	cmd.Action = func() {

		var (
			sigs   = make(chan os.Signal)
			done   = make(chan bool, 1)
			routes = http.NewRouter()
		)

		go http.RunHttpServer(routes, cfg.HttpServer.Port)

		proxy := server.New(ctx.K8S)
		go proxy.Start(cfg.ProxyServer.Port)

		// Handle SIGINT and SIGTERM.
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			for {
				select {
				case <-proxy.Ready:
					ctx.Log.Info("Listen proxy on", cfg.ProxyServer.Port, "port")
				case <-proxy.Done:
					done <- true
					return
				case <-sigs:
					proxy.Shutdown()
					<-proxy.Done
					done <- true
					return
				}
			}
		}()

		<-done

		ctx.Log.Info("Handle SIGINT and SIGTERM.")
	}
}
