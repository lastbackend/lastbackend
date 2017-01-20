package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

type Client struct {
	host string

	online     bool
	connection net.Conn

	connect   chan bool
	connected chan bool
	done      chan bool
	error     chan error

	Message chan []byte
}

func NewTCPClient(host string) *Client {

	var client = new(Client)

	client.host = host

	client.connect = make(chan bool)
	client.connected = make(chan bool)
	client.done = make(chan bool)
	client.error = make(chan error)

	client.Message = make(chan []byte)

	return client
}

func (c *Client) Connect(project, service string) (*Client, error) {

	hj, err := postHijack(c.host, project, service)
	if err != nil {
		return nil, err
	}

	c.connection = hj.Conn

	notify := make(chan error)

	go func() {
		buf := make([]byte, 1024)

		for {
			n, err := c.connection.Read(buf)
			if err != nil {
				notify <- err
				if io.EOF == err {
					return
				}
			}

			if n > 0 {
				fmt.Println("unexpected data: %s", buf[:n])
				c.Message <- buf[:n]
			}
		}
	}()

	go func() {
		for {
			select {
			case err := <-notify:
				fmt.Println("connection dropped message", err)
				c.connected <- false
				break
				//case <-time.After(time.Second * 1):
				//	fmt.Println("timeout 1, still alive")
			}
		}
	}()

	return c, nil
}

func (c *Client) Connected() chan bool {
	return c.connected
}

func (c *Client) Send(b []byte) (int, error) {
	return c.connection.Write(b)
}

func (c *Client) Close() error {
	return c.connection.Close()
}

type Hijack struct {
	Conn   net.Conn
	Reader *bufio.Reader
}

func postHijack(host, namespace, pod string) (*Hijack, error) {

	var url = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/attach", namespace, pod)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Host = host
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "tcp")

	conn, err := net.Dial("tcp", req.Host)

	// When we set up a TCP connection for hijack, there could be long periods
	// of inactivity (a long running command with no output) that in certain
	// network setups may cause ECONNTIMEOUT, leaving the client in an unknown
	// state. Setting TCP KeepAlive on the socket connection will prohibit
	// ECONNTIMEOUT unless the socket connection truly is broken
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}

	client := httputil.NewClientConn(conn, nil)
	defer client.Close()

	// Server hijacks the connection, error 'connection closed' expected
	_, err = client.Do(req)

	conn, reader := client.Hijack()

	return &Hijack{Conn: conn, Reader: reader}, err
}
