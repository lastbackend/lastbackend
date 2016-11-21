package cmd

import (
	"github.com/jawher/mow.cli"
	user "github.com/lastbackend/lastbackend/cmd/client/cmd/user"
	"github.com/lastbackend/lastbackend/cmd/client/context"


	//"lb_cli/cli"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects"
	//"github.com/lastbackend/lastbackend/cmd/client/cmd/user/structs"
	"fmt"
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
	//--------------------------------------------------------------------------------------------------------
	// метод для получения токена авторизации еще не написан. в дальнейшем он обозначен как ctx.GetUserToken()
	//--------------------------------------------------------------------------------------------------------
	app.Command("project", "project managment", func(c *cli.Cmd) {
		var (
			p_name *string
			desc *string
		)
		c.Spec = "NAME [[-d] DESC]"
		p_name = c.String(cli.StringArg{Name: "NAME", Value: "", Desc: "name of your project"})
		c.Bool(cli.BoolOpt{Name: "d description", Value: false})
		desc = c.String(cli.StringArg{Name: "DESC", Desc: "desc text"})
		c.Command("create", "create new project", func (c *cli.Cmd) {
			c.Action = func() {
				projects.Create(*p_name, *desc, ctx)
			}
		})

		c.Command("remove", "remove an existing project", func (c *cli.Cmd) {
			c.Action = func() {
				projects.Remove(*p_name, ctx)
			}
		})

		c.Action = func() {
			fmt.Println(*p_name)
			if *p_name == "" {
				projects.Get(*p_name, ctx)
			} else {
				projects.List(ctx)
			}
		}

	})
}
