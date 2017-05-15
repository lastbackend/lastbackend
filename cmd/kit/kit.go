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

package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	agent "github.com/lastbackend/lastbackend/pkg/agent/daemon"
	api "github.com/lastbackend/lastbackend/pkg/api/daemon"
	builder "github.com/lastbackend/lastbackend/pkg/builder/daemon"
	controller "github.com/lastbackend/lastbackend/pkg/controller/daemon"
	scheduler "github.com/lastbackend/lastbackend/pkg/scheduler/daemon"
	"os"

	"github.com/lastbackend/lastbackend/pkg/common/config"
	"os/signal"
	"syscall"
)

func main() {

	var (
		cfg config.Config
	)

	app := cli.App("", "Infrastructure management toolkit")

	app.Version("v version", "0.9.0")

	app.Spec = "[APP...] [OPTIONS]"

	cfg.Debug = app.Bool(cli.BoolOpt{
		Name:   "d debug", Desc: "Enable debug mode",
		EnvVar: "DEBUG", Value: false, HideValue: true,
	})

	var apps = app.Strings(cli.StringsArg{
		Name:      "APP", Desc: "schoose particular application to run [api, controller, scheduler, builder, agent]",
		HideValue: true,
	})

	cfg.Token = app.String(cli.StringOpt{
		Name:   "token", Desc: "Secret token for signature",
		EnvVar: "SECRET-TOKEN", Value: "b8tX!ae4", HideValue: true,
	})

	cfg.APIServer.Host = app.String(cli.StringOpt{
		Name:   "http-host", Desc: "Http server host",
		EnvVar: "HTTP-SERVER-HOST", Value: "", HideValue: true,
	})
	cfg.APIServer.Port = app.Int(cli.IntOpt{
		Name:   "http-port", Desc: "Http server port",
		EnvVar: "HTTP-SERVER-PORT", Value: 2967, HideValue: true,
	})

	cfg.Registry.Server = app.String(cli.StringOpt{
		Name:   "registry-server", Desc: "Http server port",
		EnvVar: "REGISTRY-SERVER", Value: "hub.registry.net", HideValue: true,
	})
	cfg.Registry.Username = app.String(cli.StringOpt{
		Name:   "registry-username", Desc: "Http server port",
		EnvVar: "REGISTRY-USERNAME", Value: "demo", HideValue: true,
	})
	cfg.Registry.Password = app.String(cli.StringOpt{
		Name:   "registry-password", Desc: "Http server port",
		EnvVar: "REGISTRY-PASSWORD", Value: "IU1yxkTD", HideValue: true,
	})

	cfg.Etcd.Endpoints = app.Strings(cli.StringsOpt{
		Name:   "etcd-endpoints", Desc: "Set etcd endpoints list",
		EnvVar: "ETCD-ENDPOINTS", Value: []string{"localhost:2379"}, HideValue: true,
	})
	cfg.Etcd.TLS.Key = app.String(cli.StringOpt{
		Name:   "etcd-tls-key", Desc: "Etcd tls key",
		EnvVar: "ETCD-TLS-KEY", Value: "", HideValue: true,
	})
	cfg.Etcd.TLS.Cert = app.String(cli.StringOpt{
		Name:   "etcd-tls-cert", Desc: "Etcd tls cert",
		EnvVar: "ETCD-TLS-CERT", Value: "", HideValue: true,
	})
	cfg.Etcd.TLS.CA = app.String(cli.StringOpt{
		Name:   "etcd-tls-ca", Desc: "Etcd tls ca",
		EnvVar: "ETCD-TLS-CA", Value: "", HideValue: true,
	})

	cfg.AgentServer.Host = app.String(cli.StringOpt{
		Name:   "host", Value: "", Desc: "Agent API server listen address",
		EnvVar: "HOST", HideValue: true,
	})
	cfg.AgentServer.Port = app.Int(cli.IntOpt{
		Name:   "port", Value: 2968, Desc: "Agent API server listen port",
		EnvVar: "PORT", HideValue: true,
	})
	cfg.Runtime.Docker.Host = app.String(cli.StringOpt{
		Name:   "docker-host", Value: "", Desc: "Provide path to Docker daemon",
		EnvVar: "DOCKER_HOST", HideValue: true,
	})
	cfg.Runtime.Docker.Certs = app.String(cli.StringOpt{
		Name:   "docker-certs", Value: "", Desc: "Provide path to Docker certificates",
		EnvVar: "DOCKER_CERT_PATH", HideValue: true,
	})
	cfg.Runtime.Docker.Version = app.String(cli.StringOpt{
		Name:   "docker-api-version", Value: "", Desc: "Docker daemon API version",
		EnvVar: "DOCKER_API_VERSION", HideValue: true,
	})
	cfg.Runtime.Docker.TLS = app.Bool(cli.BoolOpt{
		Name:   "docker-tls", Value: false, Desc: "Use secure connection to docker daemon",
		EnvVar: "DOCKER_TLS_VERIFY", HideValue: true,
	})
	cfg.Runtime.CRI = app.String(cli.StringOpt{
		Name:   "cri", Value: "docker", Desc: "Default container runtime interface",
		EnvVar: "LB_CRI", HideValue: true,
	})
	cfg.Host.Hostname = app.String(cli.StringOpt{
		Name:   "hostname", Value: "", Desc: "Agent hostname",
		EnvVar: "LB_HOSTNAME", HideValue: true,
	})

	var help = app.Bool(cli.BoolOpt{
		Name:      "h help",
		Value:     false,
		Desc:      "Show the help info and exit",
		HideValue: true,
	})

	app.Before = func() {
		if *help {
			app.PrintLongHelp()
		}
	}

	app.Action = func() {

		var (
			sigs = make(chan os.Signal)
			done = make(chan bool, 1)
		)

		if len(*apps) == 0 {
			go api.Daemon(&cfg)
			go controller.Daemon(&cfg)
			go scheduler.Daemon(&cfg)
			go builder.Daemon(&cfg)
			go agent.Daemon(&cfg)
		} else {

			for _, app := range *apps {

				if app == "api" {
					go api.Daemon(&cfg)
				}

				if app == "controller" {
					go controller.Daemon(&cfg)
				}

				if app == "scheduler" {
					go scheduler.Daemon(&cfg)
				}

				if app == "builder" {
					go builder.Daemon(&cfg)
				}

				if app == "agent" {
					go agent.Daemon(&cfg)
				}
			}
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

	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Errorf("Error: run application: %s", err.Error())
		return
	}
}
