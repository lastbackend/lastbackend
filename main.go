package main

import ()
import (
	"github.com/deployithq/deployit/daemon"
	"github.com/deployithq/deployit/handlers"
	"github.com/mitchellh/cli"
	"log"
	"os"
)

func main() {

	c := cli.NewCLI("deploy it", "0.1.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"daemon": func() (cli.Command, error) {
			return new(daemon.DaemonCommand), nil
		},
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
}
