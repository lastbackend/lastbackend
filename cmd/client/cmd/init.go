package cmd

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"fmt"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/structures"
	"time"
)

func Init(app *cli.Cli, ctx *context.Context) {

	var p_name *string
	app.Spec = "NAME [[-d] DESC]"
	p_name = app.String(cli.StringArg{Name: "NAME", Value: "", Desc: "name of your project"})
	app.Bool(cli.BoolOpt{Name: "d description", Value: false})
	desc := app.String(cli.StringArg{Name: "DESC", Desc: "desc text"})

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

	app.Command("create", "creates named project", func(*cli.Cmd) {

		if *desc == "" {
			*desc = "no description yet"
		}
		//local_time := time.Now()
		//local_time.
		fmt.Println("DESCRIPTION: ", *desc)
		var project = new(structures.Project)
		local_time := time.Now().String()
		project.ProjName, project.Description, project.Created, project.Updated = *p_name, *desc, local_time, local_time
		_, err:= rethink.Table("project_data").Insert(project).Run(session)
		if err != nil {
			fmt.Println(err)
		}

	})
}
