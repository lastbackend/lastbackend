//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package cli

import (
	"fmt"
	"github.com/jawher/mow.cli"
	p "github.com/lastbackend/lastbackend/pkg/cli/cmd/namespace"
	s "github.com/lastbackend/lastbackend/pkg/cli/cmd/service"
	"github.com/lastbackend/lastbackend/pkg/cli/config"
	"github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/cli/storage"
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"os"
)

const DEFAULT_HOST = "https://api.lastbackend.com"

func Run() {

	var (
		err error
		cfg = config.Get()
		ctx = context.Get()
	)

	app := cli.App("lb", "apps cloud hosting with integrated deployment tools")

	app.Version("v version", "0.3.0")

	app.Spec = "[-d][-H][--tls]"

	var debug = app.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "enable debug mode"})
	var host = app.String(cli.StringOpt{Name: "H host", Value: DEFAULT_HOST, Desc: "host for rest api", HideValue: true})
	var tls = app.Bool(cli.BoolOpt{Name: "tls", Value: false, Desc: "enable tls", HideValue: true})
	var help = app.Bool(cli.BoolOpt{Name: "h help", Value: false, Desc: "show the help info and exit", HideValue: true})

	app.Before = func() {

		cfg.ApiHost = *host

		if *help {
			app.PrintLongHelp()
		}

		if *debug {
			cfg.Debug = *debug
		}

		if cfg.ApiHost == DEFAULT_HOST {
			*tls = true
		}

		hcli, err := http.New(cfg.ApiHost, &http.ReqOpts{TLS: *tls})
		if err != nil {
			return
		}
		ctx.SetHttpClient(hcli)

		strg, err := storage.Get()
		if err != nil {
			panic(fmt.Sprintf("Error: init local storage %s", err.Error()))
			return
		}
		ctx.SetStorage(strg)
	}

	configure(app)

	err = app.Run(os.Args)
	if err != nil {
		panic(fmt.Sprintf("Error: run application %s", err.Error()))
		return
	}
}

func configure(app *cli.Cli) {

	app.Command("namespaces", "Display the namespace list", func(c *cli.Cmd) {
		c.Action = p.ListNamespaceCmd
	})

	app.Command("services", "Display the service list", func(c *cli.Cmd) {
		c.Action = s.ListServiceCmd
	})

	app.Command("deploy", "Deploy management", func(c *cli.Cmd) {

		c.Spec = "[URL][-t][-i][-n][--replicas]" // [-e][-p][-v]

		var url = c.String(cli.StringArg{Name: "URL", Value: "", Desc: "git repo url", HideValue: true})
		var name = c.String(cli.StringOpt{Name: "n name", Desc: "service name", HideValue: true})
		var image = c.String(cli.StringOpt{Name: "i image", Desc: "docker image name", HideValue: true})
		var template = c.String(cli.StringOpt{Name: "t template", Desc: "tempalte name", HideValue: true})

		var replicas = c.Int(cli.IntOpt{Name: "replicas", Desc: "service replicas", HideValue: true})
		//var env = c.Strings(cli.StringsOpt{Name: "e", Desc: "environment", HideValue:true})
		//var ports = c.Strings(cli.StringsOpt{Name: "p", Desc: "ports", HideValue:true})
		//var volumes = c.Strings(cli.StringsOpt{Name: "v", Desc: "volumes", HideValue:true})

		c.Action = func() {
			if len(*url) == 0 && len(*image) == 0 && len(*template) == 0 {
				c.PrintHelp()
				return
			}

			s.DeployCmd(*name, *image, *template, *url, *replicas)
		}

	})

	app.Command("namespace", "", func(c *cli.Cmd) {

		c.Spec = "[NAME]"

		var name = c.String(cli.StringArg{Name: "NAME", Value: "", Desc: "name of your namespace", HideValue: true})

		c.Command("create", "Create new namespace", func(sc *cli.Cmd) {

			sc.Spec = "[--desc]"

			var desc = sc.String(cli.StringOpt{Name: "desc", Value: "", Desc: "set description info", HideValue: true})

			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				p.CreateCmd(*name, *desc)
			}
		})

		c.Command("inspect", "Get namespace info by name", func(sc *cli.Cmd) {
			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				p.GetCmd(*name)
			}
		})

		c.Command("remove", "Remove namespace by name", func(sc *cli.Cmd) {
			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				p.RemoveCmd(*name)
			}
		})

		c.Command("switch", "switch to namespace", func(sc *cli.Cmd) {
			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				p.SwitchCmd(*name)
			}
		})

		c.Command("update", "if you wish to change name or description of the namespace", func(sc *cli.Cmd) {

			sc.Spec = "[--desc][--name]"

			var desc = sc.String(cli.StringOpt{Name: "desc", Value: "", Desc: "set description info", HideValue: true})
			var newNamespace = sc.String(cli.StringOpt{Name: "name", Value: "", Desc: "set new namespace name", HideValue: true})

			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				p.UpdateCmd(*name, *newNamespace, *desc)
			}
		})

		c.Command("current", "information about current namespace", func(sc *cli.Cmd) {
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

		var name = c.String(cli.StringArg{Name: "SERVICE_NAME", Value: "", Desc: "name of service", HideValue: true})

		c.Command("inspect", "inspect the service", func(sc *cli.Cmd) {
			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				s.InspectCmd(*name)
			}
		})

		c.Command("update", "if you wish to change configuration of the service", func(sc *cli.Cmd) {

			sc.Spec = "[--desc][--nname][--replicas]"

			var desc = sc.String(cli.StringOpt{Name: "desc", Value: "", Desc: "set description", HideValue: true})
			var nname = sc.String(cli.StringOpt{Name: "nname", Value: "", Desc: "set new name", HideValue: true})
			var replicas = sc.Int(cli.IntOpt{Name: "replicas", Value: 0, Desc: "set replicas number", HideValue: true})

			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				s.UpdateCmd(*name, *nname, *desc, *replicas)
			}
		})

		c.Command("logs", "show service logs", func(sc *cli.Cmd) {
			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				s.LogsServiceCmd(*name)
			}
		})

		c.Command("remove", "remove an existing service", func(sc *cli.Cmd) {
			sc.Action = func() {
				if len(*name) == 0 {
					c.PrintHelp()
					return
				}

				s.RemoveCmd(*name)
			}
		})
	})
}
