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
	"github.com/lastbackend/lastbackend/pkg/agent/daemon"
	"github.com/lastbackend/lastbackend/pkg/common/config"
	"os"
)

func main() {
	var cfg config.Config

	app := cli.App("lb", "apps cloud hosting with integrated deployment tools")

	app.Spec = "[OPTIONS]"

	cfg.LogLevel = app.Int(cli.IntOpt{
		Name: "debug", Desc: "Debug level mode",
		EnvVar: "DEBUG", Value: 0, HideValue: true,
	})

	app.Version("v version", "0.3.0")

	cfg.APIServer.Host = app.String(cli.StringOpt{
		Name: "http-server-host", Desc: "Http server host",
		EnvVar: "HTTP-SERVER-HOST", Value: "", HideValue: true,
	})
	cfg.APIServer.Port = app.Int(cli.IntOpt{
		Name: "http-server-port", Desc: "Http server port",
		EnvVar: "HTTP-SERVER-PORT", Value: 2967, HideValue: true,
	})

	cfg.AgentServer.Host = app.String(cli.StringOpt{
		Name: "host", Value: "", Desc: "Agent API server listen address",
		EnvVar: "HOST", HideValue: true,
	})
	cfg.AgentServer.Port = app.Int(cli.IntOpt{
		Name: "port", Value: 2968, Desc: "Agent API server listen port",
		EnvVar: "PORT", HideValue: true,
	})
	cfg.Host.Hostname = app.String(cli.StringOpt{
		Name: "hostname", Value: "", Desc: "Agent hostname",
		EnvVar: "HOSTNAME", HideValue: true,
	})
	cfg.Host.IP = app.String(cli.StringOpt{
		Name: "overwrite-ip", Value: "", Desc: "Agent host ip",
		EnvVar: "OVERWRITE_IP", HideValue: true,
	})
	cfg.Runtime.Docker.Host = app.String(cli.StringOpt{
		Name: "docker-host", Value: "", Desc: "Provide path to Docker daemon",
		EnvVar: "DOCKER_HOST", HideValue: true,
	})
	cfg.Runtime.Docker.Certs = app.String(cli.StringOpt{
		Name: "docker-certs", Value: "", Desc: "Provide path to Docker certificates",
		EnvVar: "DOCKER_CERT_PATH", HideValue: true,
	})
	cfg.Runtime.Docker.Version = app.String(cli.StringOpt{
		Name: "docker-api-version", Value: "", Desc: "Docker daemon API version",
		EnvVar: "DOCKER_API_VERSION", HideValue: true,
	})
	cfg.Runtime.Docker.TLS = app.Bool(cli.BoolOpt{
		Name: "docker-tls", Value: false, Desc: "Use secure connection to docker daemon",
		EnvVar: "DOCKER_TLS_VERIFY", HideValue: true,
	})
	cfg.Runtime.CRI = app.String(cli.StringOpt{
		Name: "cri", Value: "docker", Desc: "Default container runtime interface",
		EnvVar: "CRI", HideValue: true,
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
		daemon.Daemon(&cfg)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Errorf("Error: run application: %s", err.Error())
		return
	}
}
