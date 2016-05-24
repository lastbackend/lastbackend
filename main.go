package main

import (
	"github.com/codegangsta/cli"
	"github.com/deployithq/deployit/drivers/db"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/deployithq/deployit/drivers/log"
	"os"
)

type Env struct {
	Log      interfaces.Log
	Database interfaces.DB
}

func main() {
	app := cli.NewApp()
	app.Name = "deployit"
	app.Usage = ""

	app.Action = Action

	app.Run(os.Args)

}

func Action(c *cli.Context) error {

	env := new(Env)

	env.Log = &log.Log{
		Logger: log.New(),
	}
	env.Log.SetDebugLevel()

	database := db.Open()
	defer database.Close()

	env.Database = &db.Bolt{
		DB: database,
	}

	switch c.Args()[0] {
	case "it":
		DeployIt(env)
	}

	return nil
}

func DeployIt(env *Env) {
	env.Log.Debug("Deploy it")

}
