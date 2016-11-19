package cmd

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/cmd/client/context"
)

func Init(app *cli.Cli, ctx *context.Context) {
	app.Command("signup", "Create new account", func(c *cli.Cmd) {
		c.Action = func() {
			SignUp(ctx)
		}
	})

	app.Command("login", "Auth to account", func(c *cli.Cmd) {
		c.Action = func() {
			Auth(ctx)
		}
	})

	app.Command("whoami", "Display the current user's login name", func(c *cli.Cmd) {
		c.Action = func() {
			Whoami(ctx)
		}
	})
}
