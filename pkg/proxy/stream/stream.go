package stream

import (
	"log"
	"net"
)

const (
	proxy_host  = "localhost"
	remote_port = "9999"
	listen_port = "3333"
)

type Channel struct {
	from, to net.Conn
}

func pass_through(c *Channel) {
	buf := make([]byte, 10240)

	for {
		n, err := c.from.Read(buf)
		if err != nil {
			break
		}

		if n > 0 {
			_, err = c.to.Write(buf)
			if err != nil {
				break
			}
		}
	}

	c.from.Close()
	c.to.Close()
}

func process_connection(local net.Conn, target string) {
	remote, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Unable to connect to %s, %v\n", target, err)
	}

	go pass_through(&Channel{remote, local})
	go pass_through(&Channel{local, remote})
}

func StartProxy(port string) {
	target := net.JoinHostPort(proxy_host, remote_port)
	log.Printf("Start listening on port %s and forwarding data to %s\n",
		listen_port, target)

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Unable to start listener, %v\n", err)
	}

	for {
		if conn, err := l.Accept(); err == nil {
			go process_connection(conn, target)
		} else {
			log.Printf("Accept failed, %v\n", err)
		}
	}
}
