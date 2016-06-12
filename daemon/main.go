package daemon

import (
	"flag"
	"fmt"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/drivers/docker"
	"github.com/deployithq/deployit/drivers/localDB"
	"github.com/deployithq/deployit/drivers/log"
	"github.com/deployithq/deployit/utils"
	"os"
	"strconv"
)

type DaemonCommand struct {
	Debug bool
}

func (c *DaemonCommand) Run(args []string) int {

	log := &log.Log{
		Logger: log.New(),
	}

	paths := []string{
		fmt.Sprintf("%s/apps", env.Default_root_path),
		fmt.Sprintf("%s/tmp", env.Default_root_path),
	}

	if err := utils.CreateDirs(paths); err != nil {
		log.Fatal(err)
		return 1
	}

	// Creating flags set
	cmdFlags := flag.NewFlagSet("daemon", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		fmt.Print(c.Help())
	}

	cmdFlags.BoolVar(&c.Debug, "debug", false, "Enables debug mode")
	if c.Debug == false {
		if os.Getenv("DEPLOYIT_DEBUG") != "" {
			c.Debug = true
		}
	}

	if c.Debug {
		log.SetDebugLevel()
		log.Debug("Debug mode enabled")
	}

	log.Info("Init local db")
	ldb, _ := localDB.Init(env.Default_root_path)

	log.Info("Init daemon")

	env := &env.Env{
		LDB:        ldb,
		Log:        log,
		Containers: &docker.Containers{},
	}

	cmdFlags.IntVar(&env.Port, "port", 3000, "Daemon port")
	if c.Debug == false {
		if os.Getenv("DEPLOYIT_DAEMON_PORT") != "" {
			env.Port, _ = strconv.Atoi(os.Getenv("DEPLOYIT_DAEMON_PORT"))
		}
	}

	log.Info("Context inited")

	Route{}.Init(env)

	return 0
}

func (c *DaemonCommand) Help() string {
	return ""
}

func (c *DaemonCommand) Synopsis() string {
	return ""
}
