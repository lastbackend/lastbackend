package main

import ()
import (
	"github.com/deployithq/deployit/handlers"
	"github.com/mitchellh/cli"
	"log"
	"os"
)

func main() {

	c := cli.NewCLI("deploy it", "0.1.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"it": func() (cli.Command, error) {
			return new(handlers.ItCommand), nil
		},
		"app start": func() (cli.Command, error) {
			return &handlers.AppCommand{
				Subcommand: "start",
			}, nil
		},
		"app stop": func() (cli.Command, error) {
			return &handlers.AppCommand{
				Subcommand: "stop",
			}, nil
		},
		"app restart": func() (cli.Command, error) {
			return &handlers.AppCommand{
				Subcommand: "restart",
			}, nil
		},
		"app remove": func() (cli.Command, error) {
			return &handlers.AppCommand{
				Subcommand: "remove",
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitStatus)

	//app := cli.NewApp()
	//
	//app.Name = "deployit"
	//app.Usage = "Deploy it command line tool for deploying great apps!"
	//app.Version = "0.1"
	//
	//app.Flags = []cli.Flag{
	//	&cli.BoolFlag{
	//		Name:        "debug",
	//		Usage:       "Shows you debug logs",
	//		Destination: &handlers.Debug,
	//	},
	//}
	//
	//app.Commands = []cli.Command{
	//	{
	//		Name:        "Deploy it daemon",
	//		Aliases:     []string{"daemon"},
	//		Usage:       "Building and deploying application to host",
	//		Description: "Deploy it daemon is a server-side component for building and deploying applications to host where it is ran.",
	//		Action:      daemon.Init,
	//		Flags: []cli.Flag{
	//			&cli.IntFlag{
	//				Name:        "port",
	//				Usage:       "Port, which daemon will listen",
	//				Value:       3000,
	//				Destination: &daemon.Port,
	//			},
	//			&cli.StringFlag{
	//				Name:        "docker-uri",
	//				Usage:       "Docker daemon adress",
	//				Destination: &docker.DOCKER_URI,
	//				EnvVar:      "DOCKER_URI",
	//			},
	//			&cli.StringFlag{
	//				Name:        "docker-cert",
	//				Usage:       "Docker client certificate",
	//				Destination: &docker.DOCKER_CERT,
	//				EnvVar:      "DOCKER_CERT",
	//			},
	//			&cli.StringFlag{
	//				Name:        "docker-ca",
	//				Usage:       "Docker certificate authority that signed the registry certificate",
	//				Destination: &docker.DOCKER_CA,
	//				EnvVar:      "DOCKER_CA",
	//			},
	//			&cli.StringFlag{
	//				Name:        "docker-key",
	//				Usage:       "Docker client key",
	//				Destination: &docker.DOCKER_KEY,
	//				EnvVar:      "DOCKER_KEY",
	//			}},
	//	},
	//	{
	//		Name:    "Deploy it",
	//		Aliases: []string{"it"},
	//		Usage:   "Deploying sources from current directory",
	//		Action:  handlers.DeployIt,
	//		Flags: []cli.Flag{
	//			&cli.StringFlag{
	//				Name:        "host",
	//				Usage:       "Adress of your host, where daemon is running",
	//				Value:       "api.deployit.co",
	//				Destination: &handlers.Host,
	//			},
	//			&cli.IntFlag{
	//				Name:        "port",
	//				Usage:       "Port of daemon host",
	//				Value:       3000,
	//				Destination: &handlers.Port,
	//			},
	//			&cli.BoolFlag{
	//				Name:        "ssl",
	//				Usage:       "HTTPS mode if your daemon uses ssl",
	//				Destination: &handlers.SSL,
	//			},
	//			&cli.BoolFlag{
	//				Name:        "log",
	//				Usage:       "Show build logs",
	//				Destination: &handlers.Log,
	//			},
	//			&cli.BoolFlag{
	//				Name:        "force",
	//				Usage:       "",
	//				Destination: &handlers.Force,
	//			},
	//			&cli.StringFlag{
	//				Name:        "tag",
	//				Usage:       `Version of your app, examples: "latest", "master", "0.3", "1.9.9", etc.`,
	//				Value:       "latest",
	//				Destination: &handlers.Tag,
	//			}},
	//	},
	//	{
	//		Name:    "App management",
	//		Aliases: []string{"app"},
	//		Usage:   "App management command which allows to stop/start/restart/remove application and see its logs",
	//		Subcommands: []cli.Command{
	//			{
	//				Name:    "App start",
	//				Aliases: []string{"start"},
	//				Usage:   "Start application binded to this sources",
	//				Action:  handlers.AppStart,
	//			},
	//			{
	//				Name:    "App stop",
	//				Aliases: []string{"stop"},
	//				Usage:   "Stop application binded to this sources",
	//				Action:  handlers.AppStop,
	//			},
	//			{
	//				Name:    "App restart",
	//				Aliases: []string{"restart"},
	//				Usage:   "Restart application binded to this sources",
	//				Action:  handlers.AppRestart,
	//			},
	//			{
	//				Name:    "App remove",
	//				Aliases: []string{"remove"},
	//				Usage:   "Remove application binded to this sources",
	//				Action:  handlers.AppRemove,
	//			},
	//		},
	//		Flags: []cli.Flag{
	//			&cli.StringFlag{
	//				Name:        "host",
	//				Usage:       "Adress of your host, where daemon is running",
	//				Value:       "api.deployit.co",
	//				Destination: &handlers.Host,
	//			},
	//			&cli.IntFlag{
	//				Name:        "port",
	//				Usage:       "Port of daemon host",
	//				Value:       3000,
	//				Destination: &handlers.Port,
	//			},
	//			&cli.BoolFlag{
	//				Name:        "ssl",
	//				Usage:       "HTTPS mode if your daemon uses ssl",
	//				Destination: &handlers.SSL,
	//			},
	//			&cli.BoolFlag{
	//				Name:        "log",
	//				Usage:       "Show build logs",
	//				Destination: &handlers.Log,
	//			}},
	//	},
	//	{
	//		Name:    "Services management",
	//		Aliases: []string{"service"},
	//		Usage:   "Deploy/start/stop/restart/remove service",
	//		Action:  handlers.ServiceDeploy,
	//		Subcommands: []cli.Command{
	//			{
	//				Name:    "Service start",
	//				Aliases: []string{"start"},
	//				Usage:   "Start service",
	//				Action:  handlers.ServiceStart,
	//			},
	//			{
	//				Name:    "Service stop",
	//				Aliases: []string{"stop"},
	//				Usage:   "Stop service",
	//				Action:  handlers.ServiceStop,
	//			},
	//			{
	//				Name:    "Service restart",
	//				Aliases: []string{"restart"},
	//				Usage:   "Restart service",
	//				Action:  handlers.ServiceRestart,
	//			},
	//			{
	//				Name:    "Service remove",
	//				Aliases: []string{"remove"},
	//				Usage:   "Remove service",
	//				Action:  handlers.ServiceRemove,
	//			},
	//		},
	//		Flags: []cli.Flag{
	//			&cli.StringFlag{
	//				Name:        "host",
	//				Usage:       "Adress of your host, where daemon is running",
	//				Value:       "api.deployit.co",
	//				Destination: &handlers.Host,
	//			},
	//			&cli.IntFlag{
	//				Name:        "port",
	//				Usage:       "Port of daemon host",
	//				Value:       3000,
	//				Destination: &handlers.Port,
	//			},
	//			&cli.BoolFlag{
	//				Name:        "ssl",
	//				Usage:       "HTTPS mode if your daemon uses ssl",
	//				Destination: &handlers.SSL,
	//			},
	//			&cli.BoolFlag{
	//				Name:        "log",
	//				Usage:       "Show build logs",
	//				Destination: &handlers.Log,
	//			},
	//			&cli.StringFlag{
	//				Name:        "name",
	//				Usage:       "Service name",
	//				Value:       "",
	//				Destination: &handlers.ServiceName,
	//			},
	//		},
	//	},
	//}
	//
	//app.Run(os.Args)
}
