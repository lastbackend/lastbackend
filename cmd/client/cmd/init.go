package cmd

import (
	"github.com/jawher/mow.cli"
)

func Init(app *cli.Cli) {
	app.Command("login", "User authentication", Login)
}
