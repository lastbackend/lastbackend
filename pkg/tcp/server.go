package tcp

import (
  "fmt"
  "io"
  "net"
  "time"
)

func handleConnection(conn net.Conn) {
  defer conn.Close()
  notify := make(chan error)
  go func() {
    buf := make([]byte, 1024)
    for {
      n, err := conn.Read(buf)
      if err != nil {
        notify <- err
        return
      }
      if n > 0 {
        fmt.Println("unexpected data: %s", buf[:n])
      }
    }
  }()

  for {
    select {
    case err := <-notify:
      if io.EOF == err {
        fmt.Println("connection dropped message", err)
        return
      }
    case <-time.After(time.Second * 1):
      fmt.Println("timeout 1, still alive")
    }
  }
}

func RunTCPServer() {

  fmt.Println("Launching server...")

  ln, _ := net.Listen("tcp", ":9999")

  for {
    conn, _ := ln.Accept()
    go handleConnection(conn)
  }

}
