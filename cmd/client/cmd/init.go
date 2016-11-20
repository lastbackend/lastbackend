package cmd

import (
	"github.com/jawher/mow.cli"
	user "github.com/lastbackend/lastbackend/cmd/client/cmd/user"
	"github.com/lastbackend/lastbackend/cmd/client/context"
)

func Init(app *cli.Cli, ctx *context.Context) {
	app.Command("signup", "Create new account", func(c *cli.Cmd) {
		c.Action = func() {
			user.SignUp(ctx)
		}
	})

	app.Command("login", "Auth to account", func(c *cli.Cmd) {
		c.Action = func() {
			user.SignIn(ctx)
		}
	})

	app.Command("whoami", "Display the current user's login name", func(c *cli.Cmd) {
		c.Action = func() {
			user.Whoami(ctx)
		}
	})
}
