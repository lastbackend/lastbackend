package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func StartProxyClient() {
	conn, _ := net.Dial("tcp", ":3333")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')

		fmt.Fprintf(conn, text+"\n")

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
