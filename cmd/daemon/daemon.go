package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/cmd/daemon/cmd"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	"os"
)

func main() {

	var (
		er  error
		ctx = context.Get()
	)

	ctx.Info.Version = "0.3.0"
	ctx.Info.ApiVersion = "0.3.0"

	app := cli.App("last.backend", "apps cloud hosting with integrated deployment tools")

	app.Version("v version", fmt.Sprintf(""+
		" Version:\t%s\r\n"+
		" API version:\t%s", ctx.Info.Version, ctx.Info.ApiVersion))

	var help = app.Bool(cli.BoolOpt{Name: "h help", Value: false, Desc: "Show the help info and exit", HideValue: true})

	app.Before = func() {
		if *help {
			app.PrintLongHelp()
		}
	}

	cmd.Init(app)

	er = app.Run(os.Args)
	if er != nil {
		ctx.Log.Panic("Error: run application", er.Error())
		return
	}
}
