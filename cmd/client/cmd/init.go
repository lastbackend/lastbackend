package cmd

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/project"
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

	app.Command("projects", "Display the project list", func(c *cli.Cmd) {
		c.Action = project.List
	})

	app.Command("project", "project managment", func(c *cli.Cmd) {
		var (
			name *string
			desc *string
		)

		app.Spec = "[NAME]"

		name = c.String(cli.StringArg{Name: "NAME", Value: "", Desc: "name of your project"})
		desc = c.String(cli.StringOpt{Name: "description", Value: "", HideValue: true})

		fmt.Println(*name)

		c.Command("create", "create new project", func(c *cli.Cmd) {
			c.Action = func() {
				project.Create(*name, *desc)
			}
		})

		c.Command("remove", "remove an existing project", func(c *cli.Cmd) {
			c.Action = func() {
				project.Remove(*name)
			}
		})

		c.Action = func() {
			fmt.Println(*name)
			project.Get(*name)
		}

	})
}
