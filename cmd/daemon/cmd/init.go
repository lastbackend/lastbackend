package cmd

import (
	"github.com/jawher/mow.cli"
)

func Init(app *cli.Cli) {
	app.Command("daemon", "Run last.backend daemon", Run)
}
