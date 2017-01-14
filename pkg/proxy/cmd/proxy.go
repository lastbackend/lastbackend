package cmd

import (
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/libs/db"
	e "github.com/lastbackend/lastbackend/libs/errors"
	a "github.com/lastbackend/lastbackend/libs/log"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/proxy/stream"
	"github.com/lastbackend/lastbackend/pkg/proxy/client"
	"time"
)

func Proxy(cmd *cli.Cmd) {

	var (
		ctx = context.Get()
		err error
	)

	ctx.Log = new(a.Log)
	ctx.Log.Init()

	ctx.Storage, err = db.Init()
	if err != nil {
		ctx.Log.Panic("Error: init local storage", err.Error())
		return
	}

	session := struct {
		Token string `json:"token,omitempty"`
	}{}

	ctx.Storage.Get("session", &session)
	ctx.Token = session.Token

	cmd.Spec = "[--port]"
	var port = cmd.String(cli.StringOpt{Name: "port", Value: "", Desc: "port for your proxy", HideValue: true})

	cmd.Action = func() {
		if len(*port) == 0 {
			cmd.PrintHelp()
			return
		}

		if len(ctx.Token) == 0 {
			ctx.Log.Panic(e.NotLoggedMessage)
			return
		}

		go stream.StartProxy(*port)
		time.Sleep(1 * time.Second)
		client.StartProxyClient(ctx.Token)
	}
}
