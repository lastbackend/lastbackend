package cmd

import (
	"fmt"
	"github.com/lastbackend/lastbackend/cmd/client"
	cctx "github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/cmd/daemon"
	dctx "github.com/lastbackend/lastbackend/cmd/daemon/context"
	"github.com/jawher/mow.cli"
	"os"
)

func Start() {

	var client_ctx = cctx.Get()
	var daemon_ctx = dctx.Get()

	daemon_ctx.Info.Version = "0.1.0"
	daemon_ctx.Info.ApiVersion = "1.0"

	client_ctx.Info.Version = "0.1.0"
	client_ctx.Info.ApiVersion = "1.0"

	app := cli.App("deployit", "deploy it tool service")

	app.Version("v version", fmt.Sprintf(""+
		"Client:\r\n"+
		" Version:\t%s\r\n"+
		" API version:\t%s"+
		"\r\n\r\n"+
		"Server:\r\n"+
		" Version:\t%s\r\n"+
		" API version:\t%s", client_ctx.Info.Version, client_ctx.Info.ApiVersion, daemon_ctx.Info.Version, daemon_ctx.Info.ApiVersion))

	var help = app.Bool(cli.BoolOpt{Name: "h help", Value: false, Desc: "Show the help info and exit", HideValue: true})

	app.Before = func() {
		if *help {
			app.PrintLongHelp()
		}
	}

	app.Command("daemon", "Run deployit in daemon mode", daemon.Run)
	app.Command("init", "Init project", client.Init)

	app.Run(os.Args)
}
