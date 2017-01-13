package server

import (
	"log"
	"net"
)

const (
	daemon_host = "localhost"
	daemon_port = ":9999"
)

func StartProxyServer() {
	l, err := net.Listen("tcp", daemon_host+daemon_port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting daemon..." + daemon_host + daemon_port)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 10240)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		if n > 0 {
			_, err := conn.Write(buf)
			if err != nil {
				return
			}
		}
	}
}
