package cmd

import (
	"github.com/jawher/mow.cli"
)

func Init(app *cli.Cli) {
	app.Command("signup", "create new account", func(c *cli.Cmd) {
		c.Action = func() {
			SignUp()
		}
	})

	app.Command("login", "Auth to account", func(c *cli.Cmd) {
		c.Action = func() {
			Auth()
		}
	})
}
