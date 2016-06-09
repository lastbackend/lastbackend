package main

import (
	"github.com/deployithq/deployit/daemon"
	"github.com/deployithq/deployit/drivers/docker"
	"github.com/deployithq/deployit/handlers"
	"gopkg.in/urfave/cli.v2"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Name = "deployit"
	app.Usage = "Deploy it command line tool for deploying great apps!"
	app.Version = "0.1"

	app.Commands = []*cli.Command{
		{
			Name:        "Deploy it daemon",
			Aliases:     []string{"daemon"},
			Usage:       "Building and deploying application to host",
			Description: "Deploy it daemon is a server-side component for building and deploying applications to host where it is ran.",
			Action:      daemon.Init,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "debug",
					Usage:       "Shows you debug logs",
					Destination: &daemon.Debug,
				},
				&cli.IntFlag{
					Name:        "port",
					Usage:       "Port, which daemon will listen",
					Value:       3000,
					Destination: &daemon.Port,
				},
				&cli.StringFlag{
					Name:        "docker-uri",
					Usage:       "Docker daemon adress",
					Destination: &docker.DOCKER_URI,
					EnvVars:     []string{"DOCKER_URI"},
				},
				&cli.StringFlag{
					Name:        "docker-cert",
					Usage:       "Docker client certificate",
					Destination: &docker.DOCKER_CERT,
					EnvVars:     []string{"DOCKER_CERT"},
				},
				&cli.StringFlag{
					Name:        "docker-ca",
					Usage:       "Docker certificate authority that signed the registry certificate",
					Destination: &docker.DOCKER_CA,
					EnvVars:     []string{"DOCKER_CA"},
				},
				&cli.StringFlag{
					Name:        "docker-key",
					Usage:       "Docker client key",
					Destination: &docker.DOCKER_KEY,
					EnvVars:     []string{"DOCKER_KEY"},
				}},
		},
		{
			Name:    "Deploy it",
			Aliases: []string{"it"},
			Usage:   "Deploying sources from current directory",
			Action:  handlers.It,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "debug",
					Usage:       "Shows you debug logs",
					Destination: &handlers.Debug,
				},
				&cli.StringFlag{
					Name:        "host",
					Usage:       "Adress of your host, where daemon is running",
					Value:       "api.deployit.co",
					Destination: &handlers.Host,
				},
				&cli.IntFlag{
					Name:        "port",
					Usage:       "Port of daemon host",
					Value:       3000,
					Destination: &handlers.Port,
				},
				&cli.BoolFlag{
					Name:        "ssl",
					Usage:       "HTTPS mode if your daemon uses ssl",
					Destination: &handlers.SSL,
				},
				&cli.BoolFlag{
					Name:        "log",
					Usage:       "Show build logs",
					Destination: &handlers.Log,
				},
				&cli.BoolFlag{
					Name:        "force",
					Usage:       "",
					Destination: &handlers.Force,
				},
				&cli.StringFlag{
					Name:        "tag",
					Usage:       `Version of your app, examples: "latest", "master", "0.3", "1.9.9", etc.`,
					Value:       "latest",
					Destination: &handlers.Tag,
				}},
		},
		{
			Name:    "App management",
			Aliases: []string{"app"},
			Usage:   "App management command which allows to stop/start/restart/remove application and see its logs",
			Subcommands: []*cli.Command{
				{
					Name:    "App start",
					Aliases: []string{"start"},
					Usage:   "Start application binded to this sources",
					Action:  handlers.AppStart,
				},
				{
					Name:    "App stop",
					Aliases: []string{"stop"},
					Usage:   "Stop application binded to this sources",
					Action:  handlers.AppStop,
				},
				{
					Name:    "App restart",
					Aliases: []string{"restart"},
					Usage:   "Restart application binded to this sources",
					Action:  handlers.AppRestart,
				},
				{
					Name:    "App remove",
					Aliases: []string{"remove"},
					Usage:   "Remove application binded to this sources",
					Action:  handlers.AppRemove,
				},
			},
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "debug",
					Usage:       "Shows you debug logs",
					Destination: &handlers.Debug,
				},
				&cli.StringFlag{
					Name:        "host",
					Usage:       "Adress of your host, where daemon is running",
					Value:       "api.deployit.co",
					Destination: &handlers.Host,
				},
				&cli.IntFlag{
					Name:        "port",
					Usage:       "Port of daemon host",
					Value:       3000,
					Destination: &handlers.Port,
				},
				&cli.BoolFlag{
					Name:        "ssl",
					Usage:       "HTTPS mode if your daemon uses ssl",
					Destination: &handlers.SSL,
				},
				&cli.BoolFlag{
					Name:        "log",
					Usage:       "Show build logs",
					Destination: &handlers.Log,
				}},
		},
	}

	app.Run(os.Args)
}
