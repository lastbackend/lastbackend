package projects

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"fmt"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/structures"
	"time"
	"net/http"
	"bytes"
)




func Init(app *cli.Cli, ctx *context.Context) {

	var (
		p_name *string
		err error
		req http.Request
	)
	app.Spec = "NAME [[-d] DESC]"
	p_name = app.String(cli.StringArg{Name: "NAME", Value: "", Desc: "name of your project"})
	app.Bool(cli.BoolOpt{Name: "d description", Value: false})
	desc := app.String(cli.StringArg{Name: "DESC", Value: "no description yet",Desc: "description text"})

	app.Command("project", "project managment", )

	app.Command("create", "creates named project", func(*cli.Cmd) {

		fmt.Println("DESCRIPTION: ", *desc)
		var project = new(structures.Project)
		local_time := time.Now().String()
		project.Name, project.Description, project.Created, project.Updated = *p_name, *desc, local_time, local_time
		//project.User, project.Id = api.GetUser(), GetID()
		jsonStr := "\"name\": \"" + *p_name + "\""
		req, err = http.NewRequest("POST", "/project", bytes.NewReader(jsonStr))
		req.Header.Set("Authorization", "Bearer " /* + api.GetToken() */)



	})
	/*
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
	*/
}
