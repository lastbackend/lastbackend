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
	"github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/pkg/agent/config"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/agent/cri"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime"
	"github.com/lastbackend/lastbackend/pkg/agent/storage"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func Agent(cmd *cli.Cmd) {

	var ctx = context.Get()
	var cfg = config.Get()

	cmd.Spec = "[-d]"

	cfg.Debug = cmd.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "Enable debug mode"})

	cfg.Runtime.Docker.Host = cmd.String(cli.StringOpt{
		Name: "docker-host", Value: "", Desc: "Provide path to Docker daemon",
		EnvVar: "DOCKER_HOST", HideValue: true,
	})

	cfg.Runtime.Docker.Certs = cmd.String(cli.StringOpt{
		Name: "docker-certs", Value: "", Desc: "Provide path to Docker certificates",
		EnvVar: "DOCKER_CERT_PATH", HideValue: true,
	})

	cfg.Runtime.Docker.Version = cmd.String(cli.StringOpt{
		Name: "docker-api-version", Value: "", Desc: "Docker daemon API version",
		EnvVar: "DOCKER_API_VERSION", HideValue: true,
	})

	cfg.Runtime.Docker.TLS = cmd.Bool(cli.BoolOpt{
		Name: "docker-tls", Value: false, Desc: "Use secure connection to docker daemon",
		EnvVar: "DOCKER_TLS_VERIFY", HideValue: true,
	})

	cfg.Runtime.CRI = cmd.String(cli.StringOpt{
		Name: "cri", Value: "docker", Desc: "Default container runtime interface",
		EnvVar: "LB_CRI", HideValue: true,
	})

	cmd.Before = func() {

	}

	cmd.Action = func() {

		var (
			err  error
			crii cri.CRI
			sigs = make(chan os.Signal)
			done = make(chan bool, 1)
		)

		ctx.SetConfig(cfg)
		ctx.SetLogger(logger.New(*cfg.Debug, 0))
		ctx.SetStorage(storage.New())

		rntm := &runtime.Runtime{}
		crii, err = rntm.SetCri(cfg.Runtime)
		if err != nil {
			ctx.GetLogger().Errorf("Cannot initialize runtime: %s", err.Error())
		}

		ctx.SetCri(crii)

		if err = rntm.StartPodManager(); err != nil {
			ctx.GetLogger().Errorf("Cannot initialize pod manager: %s", err.Error())
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

		logrus.Info("Handle SIGINT and SIGTERM.")
	}

}
