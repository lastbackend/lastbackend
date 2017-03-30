package main

import (
	"os"
	"github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/pkg/agent/cmd"
)

func main() {
	var er error

	app := cli.App("lb", "apps cloud hosting with integrated deployment tools")

	app.Version("v version", "0.3.0")

	var help = app.Bool(cli.BoolOpt{Name: "h help", Value: false, Desc: "Show the help info and exit", HideValue: true})

	app.Before = func() {
		if *help {
			app.PrintLongHelp()
		}
	}

	app.Command("daemon", "Run last.backend daemon", cmd.Agent)

	er = app.Run(os.Args)
	if er != nil {
		logrus.Panic("Error: run application", er.Error())
		return
	}
}
