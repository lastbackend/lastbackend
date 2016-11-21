package cmd

import (
	"github.com/jawher/mow.cli"
	user "github.com/lastbackend/lastbackend/cmd/client/cmd/user"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
)

func Init(app *cli.Cli, ctx *context.Context, cfg *config.Config) {
	app.Command("signup", "Create new account", func(c *cli.Cmd) {
		c.Action = func() {
			user.SignUp(ctx, cfg)
		}
	})

	app.Command("login", "Auth to account", func(c *cli.Cmd) {
		c.Action = func() {
			user.SignIn(ctx, cfg)
		}
	})

	app.Command("whoami", "Display the current user's login name", func(c *cli.Cmd) {
		c.Action = func() {
			user.Whoami(ctx, cfg)
		}
	})
}
