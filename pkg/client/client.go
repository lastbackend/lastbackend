package client

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/log"
	p "github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	u "github.com/lastbackend/lastbackend/pkg/client/cmd/user"
	"github.com/lastbackend/lastbackend/pkg/client/config"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"os"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/service"
)

func Run() {

	var (
		err error
		cfg = config.Get()
		ctx = context.Get()
	)

	app := cli.App("lb", "apps cloud hosting with integrated deployment tools")

	app.Version("v version", "0.3.0")

	app.Spec = "[-d][-H]"

	var debug = app.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "Enable debug mode"})
	var host = app.String(cli.StringOpt{Name: "H host", Value: "http://localhost:3000", Desc: "Host for rest api", HideValue: true})
	var help = app.Bool(cli.BoolOpt{Name: "h help", Value: false, Desc: "Show the help info and exit", HideValue: true})

	app.Before = func() {

		cfg.ApiHost = *host

		if *help {
			app.PrintLongHelp()
		}

		ctx.Log = new(log.Log)
		ctx.Log.Init()

		if *debug {
			cfg.Debug = *debug
			ctx.Log.SetDebugLevel()
			ctx.Log.Info("Logger debug mode enabled")
		}

		ctx.HTTP = http.New(cfg.ApiHost)

		ctx.Storage = new(context.LocalStorage)

		err = ctx.Storage.Init()
		if err != nil {
			ctx.Log.Error(err)
		}

		ctx.Storage.Get("session", nil)
	}

	configure(app)

	err = app.Run(os.Args)
	if err != nil {
		ctx.Log.Panic("Error: run application", err.Error())
		return
	}
}

func configure(app *cli.Cli) {

	app.Command("signup", "Create new account", func(c *cli.Cmd) {
		c.Action = u.SignUpCmd
	})
	app.Command("login", "Auth to account", func(c *cli.Cmd) {
		c.Action = u.SignInCmd
	})
	app.Command("whoami", "Display the current user's login name", func(c *cli.Cmd) {
		c.Action = u.WhoamiCmd
	})
	app.Command("logout", "logout from account", func(c *cli.Cmd) {
		c.Action = u.LogoutCmd
	})

	app.Command("projects", "Display the project list", func(c *cli.Cmd) {
		c.Action = p.ListCmd
	})

	app.Command("project", "", func(c *cli.Cmd) {

		c.Spec = "[NAME]"

		var name = c.String(cli.StringArg{
			Name:      "NAME",
			Value:     "",
			Desc:      "name of your project",
			HideValue: true,
		})

		c.Command("create", "Create new project", func(c *cli.Cmd) {

			c.Spec = "[--desc]"

			var desc = c.String(cli.StringOpt{
				Name:      "desc",
				Value:     "",
				Desc:      "Set description info",
				HideValue: true,
			})

			c.Action = func() {
				p.CreateCmd(*name, *desc)
			}
		})

		c.Command("inspect", "Get project info by name", func(c *cli.Cmd) {
			c.Action = func() {
				p.GetCmd(*name)
			}
		})

		c.Command("remove", "Remove project by name", func(c *cli.Cmd) {
			c.Action = func() {
				p.RemoveCmd(*name)
			}
		})

	})

	app.Command("service", "Service management", func(c *cli.Cmd) {
		c.Spec = "[SERVICE_NAME]"
		var service_name = c.String(cli.StringArg{
			Name:      "SERVICE_NAME",
			Value:     "",
			Desc:      "name of service",
			HideValue: true,
		})
		c.Command("create", "create new service", func(c *cli.Cmd) {
			c.Action = func() {
				service.Create(*service_name)
			}
		})
		c.Command("inspect", "examine the inspect", func(c *cli.Cmd) {
			c.Action = func() {
				service.Inspect(*service_name)
			}
		})
		c.Command("remove", "remove an existing service", func(c *cli.Cmd) {
			c.Action = func() {
				service.Remove(*service_name)
			}
		})

	})

}
