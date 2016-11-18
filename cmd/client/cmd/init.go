package cmd

import (
	"github.com/jawher/mow.cli"
)

func Init(app *cli.Cli) {
	app.Command("login", "login to lb", func(c *cli.Cmd) {
		c.Action = func() {
			Login()
		}
	})
}
