package proxy

import (
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/proxy/client"
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
	)

	proxy := client.New("127.0.0.1:9999", ctx.Token)

	ctx.Log.Info("Listen proxy on", port, "port")

	go proxy.Start(port)

	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		proxy.Shutdown()
	}()

	<-proxy.Done

	return nil
}
