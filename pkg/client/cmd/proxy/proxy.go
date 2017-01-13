package proxy

import (
	"github.com/lastbackend/lastbackend/pkg/proxy/client"
	p "github.com/lastbackend/lastbackend/pkg/proxy/server/proxy"
	ps "github.com/lastbackend/lastbackend/pkg/proxy/server/server"
	"time"
)

func Proxy(port string) {
	go p.StartProxy(port)
	go ps.StartProxyServer()
	time.Sleep(1 * time.Millisecond)
	client.StartProxyClient()
}
