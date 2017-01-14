package cmd

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/pkg/proxy/server"
	"github.com/lastbackend/lastbackend/pkg/proxy/stream"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Proxy(cmd *cli.Cmd) {

	cmd.Spec = "[--port]"
	var port = cmd.String(cli.StringOpt{Name: "port", Value: "", Desc: "port for your proxy", HideValue: true})

	cmd.Action = func() {
		if len(*port) == 0 {
			cmd.PrintHelp()
			return
		}

		go server.StartProxyServer()
		go stream.StartProxy(*port)

		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		log.Println(<-ch)
	}
}
