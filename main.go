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

	app.Commands = []*cli.Command{
		{
			Name:        "Deploy it daemon",
			Aliases:     []string{"daemon"},
			Usage:       "",
			Description: "",
			Action:      daemon.Init,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "debug",
					Usage:       "Debug mode",
					Destination: &daemon.Debug,
				},
				&cli.IntFlag{
					Name:        "port",
					Usage:       "Daemon port",
					Value:       3000,
					Destination: &daemon.Port,
				},
				&cli.StringFlag{
					Name:        "docker-uri",
					Usage:       "",
					Destination: &docker.DOCKER_URI,
					EnvVars:     []string{"DOCKER_URI"},
				},
				&cli.StringFlag{
					Name:        "docker-cert",
					Usage:       "",
					Destination: &docker.DOCKER_CERT,
					EnvVars:     []string{"DOCKER_CERT"},
				},
				&cli.StringFlag{
					Name:        "docker-ca",
					Usage:       "",
					Destination: &docker.DOCKER_CA,
					EnvVars:     []string{"DOCKER_CA"},
				},
				&cli.StringFlag{
					Name:        "docker-key",
					Usage:       "",
					Destination: &docker.DOCKER_KEY,
					EnvVars:     []string{"DOCKER_KEY"},
				}},
		},
		{
			Name:        "",
			Aliases:     []string{"it"},
			Usage:       "",
			Description: "",
			Action:      handlers.DeployIt,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "debug",
					Usage:       "Debug mode",
					Destination: &handlers.Debug,
				},
				&cli.StringFlag{
					Name:        "host",
					Usage:       "",
					Value:       "api.deployit.co",
					Destination: &handlers.Host,
				},
				&cli.IntFlag{
					Name:        "port",
					Usage:       "",
					Value:       3000,
					Destination: &handlers.Port,
				},
				&cli.BoolFlag{
					Name:        "ssl",
					Usage:       "",
					Destination: &handlers.SSL,
				},
				&cli.BoolFlag{
					Name:        "log",
					Usage:       "",
					Destination: &handlers.Log,
				},
				&cli.BoolFlag{
					Name:        "force",
					Usage:       "",
					Destination: &handlers.Force,
				},
				&cli.StringFlag{
					Name:        "tag",
					Usage:       "",
					Value:       "latest",
					Destination: &handlers.Tag,
				}},
		},
		{
			Name:        "",
			Aliases:     []string{"app"},
			Usage:       "",
			Description: "",
			Subcommands: []*cli.Command{
				{
					Name:        "",
					Aliases:     []string{"start"},
					Usage:       "",
					Description: "",
					Action:      handlers.AppStart,
				},
				{
					Name:        "",
					Aliases:     []string{"stop"},
					Usage:       "",
					Description: "",
					Action:      handlers.AppStop,
				},
				{
					Name:        "",
					Aliases:     []string{"restart"},
					Usage:       "",
					Description: "",
					Action:      handlers.AppRestart,
				},
				{
					Name:        "",
					Aliases:     []string{"remove"},
					Usage:       "",
					Description: "",
					Action:      handlers.AppRemove,
				},
			},
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "debug",
					Usage:       "Debug mode",
					Destination: &handlers.Debug,
				},
				&cli.StringFlag{
					Name:        "host",
					Usage:       "",
					Value:       "api.deployit.co",
					Destination: &handlers.Host,
				},
				&cli.IntFlag{
					Name:        "port",
					Usage:       "",
					Value:       3000,
					Destination: &handlers.Port,
				},
				&cli.BoolFlag{
					Name:        "ssl",
					Usage:       "",
					Destination: &handlers.SSL,
				},
				&cli.BoolFlag{
					Name:        "log",
					Usage:       "",
					Destination: &handlers.Log,
				}},
		},
	}

	app.Run(os.Args)
}
