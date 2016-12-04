package registry

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/registry/pkg/registry/cmd"
	"github.com/lastbackend/registry/pkg/registry/context"
	"os"
)

func Run() {

	var (
		er  error
		ctx = context.Get()
	)

	app := cli.App("registry", "")

	app.Version("v version", "0.1.0")

	var help = app.Bool(cli.BoolOpt{Name: "h help", Value: false, Desc: "Show the help info and exit", HideValue: true})

	app.Before = func() {
		if *help {
			app.PrintLongHelp()
		}
	}

	app.Command("daemon", "Run registry daemon", cmd.Daemon)

	er = app.Run(os.Args)
	if er != nil {
		ctx.Log.Panic("Error: run application", er.Error())
		return
	}
}
