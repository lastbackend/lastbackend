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
			Name:  "name",
			Usage: "",
			Value: "app",
		},
		cli.StringFlag{
			Name:  "tag",
			Usage: "",
			Value: "latest",
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
			Name:        "Deploy url command",
			Aliases:     []string{"url"},
			Usage:       "",
			Description: "",
			Action:      handlers.DeployURL,
		},
	}

	app.Run(os.Args)

}
