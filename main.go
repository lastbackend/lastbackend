package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
)

func main() {
	log.Println("Start cli")

	app := cli.NewApp()
	app.Name = "deployit"
	app.Usage = ""

	app.Action = Action

	app.Run(os.Args)

}

func Action(c *cli.Context) error {
	log.Println("Start")

	return nil
}
