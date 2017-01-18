package proxy

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type Proxy struct {
	localAddr  *net.TCPAddr
	remoteAddr *net.TCPAddr
	listener   net.Listener
	auth       *AuthOpts
	Closed     chan bool
}

type ProxyOpts struct {
	Auth *AuthOpts
}

type AuthOpts struct {
	Token string
}

func New(local, remote string, opts *ProxyOpts) (*Proxy, error) {

	var (
		err    error
		proxy = new(Proxy)
	)

	proxy.localAddr, err = net.ResolveTCPAddr("tcp", local)
	if err != nil {
		return nil, err
	}

	proxy.remoteAddr, err = net.ResolveTCPAddr("tcp", remote)
	if err != nil {
		return nil, err
	}

	if opts != nil {
		proxy.auth = opts.Auth
	}

	return proxy, nil
}

func (p *Proxy) Listen() error {

	var err error

	client, err := net.DialTCP("tcp", nil, p.remoteAddr)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to connect to %s, %v\n", client, err))
		return err
	}

	err = client.SetKeepAlive(true)
	if err != nil {
		return err
	}

	err = client.SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		return err
	}

	notify := make(chan error)

	go func() {
		buf := make([]byte, 1024)

		for {
			n, err := client.Read(buf)
			if err != nil {
				notify <- err
				if io.EOF == err {
					return
				}
			}

			if n > 0 {
				fmt.Println("unexpected data: %s", buf[:n])
			}
		}
	}()

	go func() {
		for {
			select {
			case err := <-notify:
				fmt.Println("connection dropped message", err)
				break
			case <-time.After(time.Second * 1):
				fmt.Println("timeout 1, still alive")
			}
		}
	}()

	p.listener, err = net.ListenTCP("tcp", p.localAddr)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to start listener, %v\n", err))
	}

	for {
		conn, err := p.listener.Accept()
		if err != nil {
			fmt.Printf("Accept failed, %v\n", err)
		} else {
			go p.proxy(conn, client)
		}
	}

	return nil
}

func (p *Proxy) Close() error {
	return p.listener.Close()
}

func (p *Proxy) proxy(conn, client net.Conn) {
	stream(conn, client)
	stream(client, conn)
}

func stream(from, to net.Conn) {
	go func() {
		defer from.Close()
		defer to.Close()
		io.Copy(from, to)
	}()
}
