package daemon

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/pkg/daemon/cmd"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"os"
)

func Run() {

	var (
		er  error
		ctx = context.Get()
	)

	app := cli.App("lb", "apps cloud hosting with integrated deployment tools")

	app.Version("v version", "0.3.0")

	var help = app.Bool(cli.BoolOpt{Name: "h help", Value: false, Desc: "Show the help info and exit", HideValue: true})

	app.Before = func() {
		if *help {
			app.PrintLongHelp()
		}
	}
	
	app.Command("daemon", "Run last.backend daemon", cmd.Daemon)

	er = app.Run(os.Args)
	if er != nil {
		ctx.Log.Panic("Error: run application", er.Error())
		return
	}
}
