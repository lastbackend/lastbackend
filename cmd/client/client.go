package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/cmd/client/cmd"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/libs/log"
	"os"
	"github.com/boltdb/bolt"
	"github.com/lastbackend/lastbackend/utils"
)

func main() {

	var (
		er  error
		cfg = config.Get()
		ctx = context.Get()
	)

	ctx.Info.Version = "0.3.0"
	ctx.Info.ApiVersion = "0.3.0"

	app := cli.App("last.backend", "apps cloud hosting with integrated deployment tools")

	app.Version("v version", fmt.Sprintf(""+
		"Client:\r\n"+
		" Version:\t%s\r\n"+
		" API version:\t%s", ctx.Info.Version, ctx.Info.ApiVersion))

	app.Spec = "[-d]"

	var debug = app.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "Enable debug mode"})
	var help = app.Bool(cli.BoolOpt{Name: "h help", Value: false, Desc: "Show the help info and exit", HideValue: true})

	app.Before = func() {
		if *help {
			app.PrintLongHelp()
		}

		ctx.Log = new(log.Log)
		ctx.Log.Init()

		if *debug {
			cfg.Debug = *debug
			ctx.Log.SetDebugLevel()
			ctx.Log.Info("Logger debug mode enabled")
		}

		db, err := bolt.Open(utils.GetHomeDir + "/.lb/.session", 0755, nil)
		if err != nil {
			ctx.Log.Fatal(err)
		}
		defer db.Close()

	}

	cmd.Init(app)

	er = app.Run(os.Args)
	if er != nil {
		ctx.Log.Panic("Error: run application", er.Error())
		return
	}
}
