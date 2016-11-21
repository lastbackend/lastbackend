package cmd

import (
	"github.com/jawher/mow.cli"
	u "github.com/lastbackend/lastbackend/cmd/client/cmd/user"
)

func Init(app *cli.Cli) {
	app.Command("signup", "Create new account", func(c *cli.Cmd) {
		c.Action = u.SignUp
	})

	app.Command("login", "Auth to account", func(c *cli.Cmd) {
		c.Action = u.SignIn
	})

	app.Command("whoami", "Display the current user's login name", func(c *cli.Cmd) {
		c.Action = u.Whoami
	})
}
