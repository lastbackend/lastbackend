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
		done = make(chan bool, 1)
	)

	signal.Notify(sigs, os.Interrupt)

	proxy := client.New("127.0.0.1:9999", ctx.Token, "aaaca8b4-6198-491c-8bb4-edb8f1740945", "lb-redis-4065565212-a79sb")

	go proxy.Start(port)

	go func() {
		select {
		case <-proxy.Ready:
			ctx.Log.Info("Listen proxy on", port, "port")
		case <-proxy.Done:
			done <- true
			return
		case <-sigs:
			proxy.Shutdown()
			done <- true
			return
		}
	}()

	<-done

	return nil
}
