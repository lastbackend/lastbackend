package server

import (
	"log"
	"net"
	"github.com/lastbackend/lastbackend/libs/model"
)

const (
	daemon_host = "localhost"
	daemon_port = ":9999"
)

func StartProxyServer() {

	l, err := net.Listen("tcp", daemon_host+daemon_port)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Starting daemon..." + daemon_host + daemon_port)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panic(err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 10240)

	i := 0

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Panic(err)
			return
		}

		if i == 0 {
			i++

			a := new(model.Session)
			err = a.Decode(string(buf[:248]))
			if err != nil {
				conn.Close()
				return
			}
		}

		log.Println("lal")

		if n > 0 {
			_, err := conn.Write(buf)
			if err != nil {
				log.Panic(err)
				return
			}
		}
	}
}
