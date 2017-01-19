package client

import (
	"fmt"
	"io"
	"net"
)

type Proxy struct {
	port       string
	address    string
	token      string
	authorized bool

	close chan int
	error chan error

	Done chan int
}

func New(address, token string) *Proxy {

	var proxy = new(Proxy)

	proxy.address = address
	proxy.token = token
	proxy.close = make(chan int)
	proxy.error = make(chan error)

	proxy.Done = make(chan int)

	return proxy
}

func (p *Proxy) Start(port int) {

	server := NewTCPServer(port)
	client, err := NewTCPClient().Connect(p.address)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = client.Send([]byte(p.token))
	if err != nil {
		fmt.Println(err)
		return
	}

	err = server.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		select {
		case <-p.close:
			client.Close()
			server.Close()
			close(p.Done)
			return
		}
	}()

	go func() {
		for {
			select {
			case msg := <-client.Message:
				fmt.Println("client.message", string(msg))

				if !p.authorized {
					p.authorized = true
					if err := server.Start(); err != nil {
						fmt.Println(err)
					}

					server.Accept(func(conn net.Conn) {
						go p.copy(conn, client.connection)
					})

				} else {
					server.Send(msg)
				}
			}
		}
	}()
}

func (p *Proxy) Shutdown() {
	close(p.close)
}

func (p *Proxy) copy(from, to net.Conn) {
	select {
	case <-p.close:
		return
	default:
		if _, err := io.Copy(to, from); err != nil {
			return
		}
	}
}
