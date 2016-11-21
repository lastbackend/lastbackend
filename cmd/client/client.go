package main

import (
	"github.com/boltdb/bolt"
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/cmd/client/cmd"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/log"
	"github.com/lastbackend/lastbackend/utils"
	"os"
)

func main() {

	var (
		err error
		cfg = config.Get()
		ctx = context.Get()
	)

	app := cli.App("last.backend", "apps cloud hosting with integrated deployment tools")

	app.Version("v version", "0.3.0")

	app.Spec = "[-d][-H]"

	var debug = app.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "Enable debug mode"})
	var host = app.String(cli.StringOpt{Name: "H host", Value: "http://localhost:3000", Desc: "Host for rest api", HideValue: true})
	var help = app.Bool(cli.BoolOpt{Name: "h help", Value: false, Desc: "Show the help info and exit", HideValue: true})

	app.Before = func() {

		cfg.ApiHost = *host

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

		ctx.HTTP = http.New(cfg.ApiHost)

		dir := utils.GetHomeDir() + "/.lb"

		utils.MkDir(dir, 0755)
		ctx.Storage, err = bolt.Open(dir+"/lb.db", 0755, nil)
		if err != nil {
			ctx.Log.Fatal(err)
		}
	}

	cmd.Init(app)

	err = app.Run(os.Args)
	if err != nil {
		ctx.Log.Panic("Error: run application", err.Error())
		return
	}
}
