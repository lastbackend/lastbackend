package daemon

import (
	"fmt"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/drivers/localDB"
	"github.com/deployithq/deployit/drivers/log"
	"github.com/deployithq/deployit/utils"
	"gopkg.in/urfave/cli.v2"
)

var Host string
var Port int
var Debug bool

func Init(c *cli.Context) error {

	log := &log.Log{
		Logger: log.New(),
	}

	paths := []string{
		fmt.Sprintf("%s/apps", env.Default_root_path),
		fmt.Sprintf("%s/tmp", env.Default_root_path),
	}

	if err := utils.CreateDirs(paths); err != nil {
		log.Fatal(err)
		return err
	}

	if Debug {
		log.SetDebugLevel()
		log.Debug("Debug mode enabled")
	}

	log.Info("Init local db")
	ldb, _ := localDB.Init(fmt.Sprintf("%s/storage/", env.Default_root_path))

	log.Info("Init daemon")

	env := &env.Env{
		LDB:  ldb,
		Log:  log,
		Host: Host,
	}

	log.Info("Context inited")

	Route{}.Init(env)

	return nil
}
