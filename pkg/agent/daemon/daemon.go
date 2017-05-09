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

package daemon

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/pkg/agent/config"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/agent/events/listener"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime/cri/cri"
	"github.com/lastbackend/lastbackend/pkg/agent/storage"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"os"
	"os/signal"
	"syscall"
)

func Agent(cmd *cli.Cmd) {

	var ctx = context.Get()
	var cfg = config.Get()

	cmd.Spec = ""

	cfg.Debug = *cmd.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "Enable debug mode"})

	cfg.Runtime.Docker.Host = *cmd.String(cli.StringOpt{
		Name: "docker-host", Value: "", Desc: "Provide path to Docker daemon",
		EnvVar: "DOCKER_HOST", HideValue: true,
	})

	cfg.Runtime.Docker.Certs = *cmd.String(cli.StringOpt{
		Name: "docker-certs", Value: "", Desc: "Provide path to Docker certificates",
		EnvVar: "DOCKER_CERT_PATH", HideValue: true,
	})

	cfg.Runtime.Docker.Version = *cmd.String(cli.StringOpt{
		Name: "docker-api-version", Value: "", Desc: "Docker daemon API version",
		EnvVar: "DOCKER_API_VERSION", HideValue: true,
	})

	cfg.Runtime.Docker.TLS = *cmd.Bool(cli.BoolOpt{
		Name: "docker-tls", Value: false, Desc: "Use secure connection to docker daemon",
		EnvVar: "DOCKER_TLS_VERIFY", HideValue: true,
	})

	cfg.Runtime.CRI = *cmd.String(cli.StringOpt{
		Name: "cri", Value: "docker", Desc: "Default container runtime interface",
		EnvVar: "LB_CRI", HideValue: true,
	})

	cfg.APIServer.Host = *cmd.String(cli.StringOpt{
		Name: "host", Value: "0.0.0.0", Desc: "API server listen address",
		EnvVar: "LB_AGENT_HOST", HideValue: true,
	})

	cfg.APIServer.Port = *cmd.Int(cli.IntOpt{
		Name: "port", Value: 2968, Desc: "API server listen port",
		EnvVar: "LB_AGENT_PORT", HideValue: true,
	})

	cfg.Host.Hostname = *cmd.String(cli.StringOpt{
		Name: "hostname", Value: "", Desc: "Agent hostname",
		EnvVar: "LB_HOSTNAME", HideValue: true,
	})

	cmd.Before = func() {

	}

	cmd.Action = func() {

		var (
			err  error
			sigs = make(chan os.Signal)
			done = make(chan bool, 1)
		)

		rntm := runtime.Get()
		crii, err := cri.New(cfg.Runtime)
		if err != nil {
			ctx.GetLogger().Errorf("Cannot initialize runtime: %s", err.Error())
		}

		ctx.SetConfig(cfg)
		ctx.SetLogger(logger.New(cfg.Debug, 9))
		ctx.SetStorage(storage.New())

		client, err := http.New(fmt.Sprintf("%s:%d", cfg.APIServer.Host, cfg.APIServer.Port), &http.ReqOpts{})
		if err != nil {
			ctx.GetLogger().Errorf("Cannot initialize http client: %s", err.Error())
		}
		ctx.SetHttpClient(client)
		ctx.SetEventListener(listener.New(ctx.GetHttpClient(), rntm.GetSpecChan()))

		ctx.SetCri(crii)

		if err = rntm.StartPodManager(); err != nil {
			ctx.GetLogger().Errorf("Cannot initialize pod manager: %s", err.Error())
		}

		if err = rntm.StartEventListener(); err != nil {
			ctx.GetLogger().Errorf("Cannot initialize pod manager: %s", err.Error())
		}

		rntm.Loop()

		go func() {
			if err := Listen(cfg.APIServer.Host, cfg.APIServer.Port); err != nil {
				ctx.GetLogger().Warnf("Http server start error: %s", err.Error())
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

		logrus.Info("Handle SIGINT and SIGTERM.")
	}

}
