package main

import (
	"github.com/codegangsta/cli"
	"github.com/deployithq/deployit/daemon"
	"github.com/deployithq/deployit/drivers/docker"
	"github.com/deployithq/deployit/handlers"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "deployit"
	app.Usage = "Deploy it command line tool for deploying great apps!"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Debug mode",
			Destination: &handlers.Debug,
		},
		cli.StringFlag{
			Name:        "host",
			Usage:       "",
			Value:       "https://api.deployit.co", // TODO: change to hub.deployit.io
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
		},
		cli.StringFlag{
			Name:        "docker-uri",
			Usage:       "",
			Value:       "",
			Destination: &docker.DOCKER_URI,
			EnvVar:      "DOCKER_URI",
		},
		cli.StringFlag{
			Name:        "docker-cert",
			Usage:       "",
			Value:       "",
			Destination: &docker.DOCKER_CERT,
			EnvVar:      "DOCKER_CERT",
		},
		cli.StringFlag{
			Name:        "docker-ca",
			Usage:       "",
			Value:       "",
			Destination: &docker.DOCKER_CA,
			EnvVar:      "DOCKER_CA",
		},
		cli.StringFlag{
			Name:        "docker-key",
			Usage:       "",
			Value:       "",
			Destination: &docker.DOCKER_KEY,
			EnvVar:      "DOCKER_KEY",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "Deploy it command",
			Aliases:     []string{"it"},
			Usage:       "Use it when you want to deploy sources of current repository",
			Description: "This command deplos sources from current directory and sends it to Deployit servers for deploying",
			Action:      handlers.DeployIt,
		},
		{
			Name:        "Deploy it daemon",
			Aliases:     []string{"daemon"},
			Usage:       "",
			Description: "",
			Action:      daemon.Init,
		},
	}

	app.Run(os.Args)

}
