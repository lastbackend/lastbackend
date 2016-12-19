package client

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/libs/db"
	"github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/log"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/deploy"
	p "github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/service"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/template"
	u "github.com/lastbackend/lastbackend/pkg/client/cmd/user"
	"github.com/lastbackend/lastbackend/pkg/client/config"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"os"
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

		ctx.Storage, err = db.Init()
		if err != nil {
			ctx.Log.Panic("Error: init local storage", err.Error())
			return
		}

		session := struct {
			Token string `json:"token,omitempty"`
		}{}

		ctx.Storage.Get("session", &session)
		ctx.Token = session.Token
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

	app.Command("deploy", "Deploy management", func(c *cli.Cmd) {

		//c.Command("it", "deploy local sources", func(sc *cli.Cmd) {
		//	sc.Action = func() {
		//		// TODO: Deploy it command for deploy local sources in LB host
		//	}
		//})

		c.Spec = "[TARGET]"

		var target = c.String(cli.StringArg{
			Name:      "TARGET",
			Value:     "",
			Desc:      "template name or git repo url",
			HideValue: true,
		})

		c.Action = func() {
			if len(*target) == 0 {
				c.PrintHelp()
				return
			}

			deploy.DeployTargetCmd(*target)
		}

	})

	app.Command("project", "", func(c *cli.Cmd) {

		c.Spec = "[NAME]"

		var name = c.String(cli.StringArg{
			Name:      "NAME",
			Value:     "",
			Desc:      "name of your project",
			HideValue: true,
		})

		c.Command("create", "Create new project", func(sc *cli.Cmd) {

			sc.Spec = "[--desc]"

			var desc = sc.String(cli.StringOpt{
				Name:      "desc",
				Value:     "",
				Desc:      "Set description info",
				HideValue: true,
			})

			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				p.CreateCmd(*name, *desc)
			}
		})

		c.Command("inspect", "Get project info by name", func(sc *cli.Cmd) {
			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				p.GetCmd(*name)
			}
		})

		c.Command("remove", "Remove project by name", func(sc *cli.Cmd) {
			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				p.RemoveCmd(*name)
			}
		})

		c.Command("switch", "switch to project", func(sc *cli.Cmd) {
			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				p.SwitchCmd(*name)
			}
		})

		c.Command("update", "if you wish to change name or description of the project", func(sc *cli.Cmd) {

			sc.Spec = "[--desc][--name]"

			var desc = sc.String(cli.StringOpt{
				Name:      "desc",
				Value:     "",
				Desc:      "Set description info",
				HideValue: true,
			})

			var newProjectName = sc.String(cli.StringOpt{
				Name:      "name",
				Value:     "",
				Desc:      "Set new project name",
				HideValue: true,
			})

			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				p.UpdateCmd(*name, *newProjectName, *desc)
			}
		})

		c.Command("current", "information about current project", func(sc *cli.Cmd) {
			sc.Action = func() {
				if len(*name) != 0 {
					c.PrintHelp()
					return
				}
				p.CurrentCmd()
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
		c.Command("create", "create new service", func(sc *cli.Cmd) {
			sc.Action = func() {
				service.CreateCmd(*service_name)
			}
		})
		c.Command("inspect", "inspect the service", func(sc *cli.Cmd) {
			sc.Action = func() {
				service.InspectCmd(*service_name)
			}
		})
		c.Command("update", "if you wish to change configuration of the service", func(sc *cli.Cmd) {

			sc.Action = func() {
				service.UpdateCmd(*service_name)
			}
		})
		c.Command("remove", "remove an existing service", func(sc *cli.Cmd) {
			sc.Action = func() {
				service.RemoveCmd(*service_name)
			}
		})
	})

	app.Command("templates", "view templates", func(c *cli.Cmd) {
		c.Action = func() {
			template.ListCmd()
		}
	})
}
