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

	"github.com/lastbackend/lastbackend/pkg/controller/config"
	"github.com/lastbackend/lastbackend/pkg/controller/context"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"os/signal"
	"syscall"

	"os"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime"
)

func Daemon(_cfg *_cfg.Config ) {

	var (
		ctx = context.Get()
		cfg = config.Set(_cfg)
		log = logger.New(*cfg.Debug, 9)
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log.Info("Start State Contoller")

	ctx.SetConfig(cfg)
	ctx.SetLogger(log)

	stg, err := storage.Get(cfg.GetEtcdDB())
	if err != nil {
		panic(err)
	}
	ctx.SetStorage(stg)

	// Initialize Runtime
	r := runtime.NewRuntime(ctx)
  r.Loop()
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
}
