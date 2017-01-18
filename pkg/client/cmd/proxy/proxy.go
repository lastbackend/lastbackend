package proxy

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/proxy"
	"os"
	"os/signal"
)

func ProxyCmd(port int) {

	var (
		ctx = context.Get()
	)

	err := Proxy(port)
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func Proxy(port int) error {

	var (
		ctx  = context.Get()
		sigs = make(chan os.Signal, 1)
		done = make(chan bool, 1)
	)

	var local = fmt.Sprintf("127.0.0.1:%d", port)
	var remote = "127.0.0.1:9999"
	var opts = new(proxy.ProxyOpts)

	opts.Auth = new(proxy.AuthOpts)
	opts.Auth.Token = ctx.Token

	p, err := proxy.New(local, remote, opts)
	if err != nil {
		return err
	}

	ctx.Log.Info("Listen proxy on", port, "port")

	go p.Listen()

	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		p.Close()
		done <- true
	}()

	<-done

	return nil
}
