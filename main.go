package main

import (
	"github.com/codegangsta/cli"
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
	}

	app.Run(os.Args)

}
