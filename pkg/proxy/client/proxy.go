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
	project    string
	service    string
	authorized bool

	close chan int

	Ready chan bool
	Done  chan int
}

func New(address, token, project, service string) *Proxy {

	var proxy = new(Proxy)

	proxy.address = address
	proxy.token = token
	proxy.project = project
	proxy.service = service
	proxy.close = make(chan int)

	proxy.Done = make(chan int)
	proxy.Ready = make(chan bool)

	return proxy
}

func (p *Proxy) Start(port int) {

	server := NewTCPServer(port)

	client, err := NewTCPClient().Connect(p.address)
	if err != nil {
		fmt.Println(err)
		return
	}

	var rawData = []byte(fmt.Sprintf(`{"project":"%s","service":"%s","token":"%s"}`, p.project, p.service, p.token))

	_, err = client.Send(rawData)
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
				if !p.authorized {
					if `{"allow":true}` != string(msg) {
						break
					}

					if !server.Running {
						if err := server.Start(); err != nil {
							fmt.Println(err)
						}
						fmt.Printf("Listen proxy on %d port\n", port)
						//p.Ready <- true
					}

					server.Accept(func(conn net.Conn) {
						go p.copy(conn, client.connection)
					})

					p.authorized = true

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
