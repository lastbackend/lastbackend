package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"log"
)

func StartProxyClient(token string) {
	conn, _ := net.Dial("tcp", ":3333")

	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(conn, token+"\n")
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)

	for {
		fmt.Print("Text to send: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Panic(err)
		}

		fmt.Fprintf(conn, text+"\n")

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Panic(err)
		}
		fmt.Print("Message from server: " + message)
	}
}
