package cmd

import (
	"github.com/deployithq/deployit/cmd/daemon"
	"github.com/jawher/mow.cli"
	"os"
	"github.com/deployithq/deployit/cmd/daemon/context"
)

func Start() {

	var ctx = context.Get()
	ctx.Version = "0.1.0"

	app := cli.App("deployit", "deploy it tool service")

	app.Version("v version", "deployit "+ctx.Version)

	var help = app.Bool(cli.BoolOpt{Name: "h help", Value: false, Desc: "Show the help info and exit", HideValue: true})

	app.Before = func() {
		if *help {
			app.PrintLongHelp()
		}
	}

	app.Command("daemon", "Run deployit in daemon mode", daemon.Run)

	app.Run(os.Args)
}
