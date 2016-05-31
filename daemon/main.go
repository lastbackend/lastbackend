package daemon

import (
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/drivers/log"
	"gopkg.in/urfave/cli.v2"
)

var Host string
var Port int
var Debug bool

func Init(c *cli.Context) error {

	log := &log.Log{
		Logger: log.New(),
	}

	if Debug {
		log.SetDebugLevel()
		log.Debug("Debug mode enabled")
	}

	log.Info("Init daemon")

	env := &env.Env{
		Log: log,
		Host: Host,
	}

	log.Info("Context inited")

	Route{}.Init(env)

	return nil
}
