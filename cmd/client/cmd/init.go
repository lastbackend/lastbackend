package cmd

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/cmd/client/context"
)

func Init(app *cli.Cli, ctx *context.Context) {
	app.Command("signup", "create new account", func(c *cli.Cmd) {
		c.Action = func() {
			SignUp()
		}
	})

	app.Command("login", "Auth to account", func(c *cli.Cmd) {
		c.Action = func() {
			Auth(ctx)
		}
	})
}
