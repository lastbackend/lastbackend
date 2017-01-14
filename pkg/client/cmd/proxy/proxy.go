package proxy

import (
	ps "github.com/lastbackend/lastbackend/pkg/proxy/server"
	p "github.com/lastbackend/lastbackend/pkg/proxy/stream"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Proxy(port string) {
	go p.StartProxy(port)
	go ps.StartProxyServer()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
