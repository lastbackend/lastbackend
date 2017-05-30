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
	_cfg "github.com/lastbackend/lastbackend/pkg/common/config"

	"fmt"
	"github.com/lastbackend/lastbackend/pkg/agent/config"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/agent/events/listener"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime/cri/cri"
	"github.com/lastbackend/lastbackend/pkg/cache"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"os"
	"os/signal"
	"syscall"
)

func Daemon(_cfg *_cfg.Config) {

	var (
		ctx  = context.Get()
		cfg  = config.Set(_cfg)
		log  = logger.New("Agent", *cfg.LogLevel)
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log.Info("Start Agent")

	rntm := runtime.Get()

	id, err := rntm.Register()
	if err != nil {
		log.Errorf("Agent can not be registered: %s", err.Error())
		return
	}

	ctx.SetID(id)

	crii, err := cri.New(cfg.Runtime)
	if err != nil {
		ctx.GetLogger().Errorf("Cannot initialize runtime: %s", err.Error())
	}

	ctx.SetConfig(cfg)
	ctx.SetLogger(log)
	ctx.SetCache(cache.New(log))

	var host string = "0.0.0.0"
	if cfg.APIServer.Host != nil && *cfg.APIServer.Host != "" {
		host = *cfg.APIServer.Host
	}

	client, err := http.New(fmt.Sprintf("%s:%d", host, *cfg.APIServer.Port), &http.ReqOpts{})
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
		ctx.GetLogger().Errorf("Cannot initialize event listener: %s", err.Error())
	}

	rntm.Loop()

	go func() {
		if err := Listen(*cfg.AgentServer.Host, *cfg.AgentServer.Port); err != nil {
			log.Warnf("Http agent server start error: %s", err.Error())
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

	log.Info("Handle SIGINT and SIGTERM.")

	return
}
