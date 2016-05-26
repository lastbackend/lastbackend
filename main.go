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

	app.Commands = []cli.Command{
	{
		Name:        "Deploy it command",
		Aliases:     []string{"it"},
		Usage:       "Use it when you want to deploy sources of current repository",
		Description: "This command deplos sources from current directory and sends it to Deployit servers for deploying",
		Action:      handlers.DeployIt,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "debug",
				Usage:       "Debug mode",
				Destination: &handlers.Debug,
			},
			cli.StringFlag{
				Name:        "host",
				Usage:       "",
				Value:       "https://api.deployit.co",
				Destination: &handlers.Host,
			},
			cli.StringFlag{
				Name:        "name",
				Usage:       "",
				Value:       "app",
				Destination: &handlers.AppName,
			},
			cli.StringFlag{
				Name:        "tag",
				Usage:       "",
				Value:       "latest",
				Destination: &handlers.Tag,
			}},
	}, {
		Name:        "Deploy it daemon",
		Aliases:     []string{"daemon"},
		Usage:       "",
		Description: "",
		Action:      daemon.Init,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "debug",
				Usage:       "Debug mode",
				Destination: &daemon.Debug,
			},
			cli.IntFlag{
				Name:        "port",
				Usage:       "Daemon port",
				Value:       3000,
				Destination: &daemon.Port,
			},
			cli.StringFlag{
				Name:        "docker-uri",
				Usage:       "",
				Value:       "",
				Destination: &docker.DOCKER_URI,
				EnvVars:     []string{"DOCKER_URI"},
			},
			cli.StringFlag{
				Name:        "docker-cert",
				Usage:       "",
				Value:       "",
				Destination: &docker.DOCKER_CERT,
				EnvVars:     []string{"DOCKER_CERT"},
			},
			cli.StringFlag{
				Name:        "docker-ca",
				Usage:       "",
				Value:       "",
				Destination: &docker.DOCKER_CA,
				EnvVars:     []string{"DOCKER_CA"},
			},
			cli.StringFlag{
				Name:        "docker-key",
				Usage:       "",
				Value:       "",
				Destination: &docker.DOCKER_KEY,
				EnvVars:     []string{"DOCKER_KEY"},
			}},
		},
	}

	app.Run(os.Args)

}
