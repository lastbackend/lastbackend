package client

import (
	"fmt"
	"io"
	"net"
	"time"
)

type Client struct {
	online     bool
	connection *net.TCPConn

	connect   chan bool
	connected chan bool
	done      chan bool
	error     chan error

	Message chan []byte
}

func NewTCPClient() *Client {

	var client = new(Client)

	client.connect = make(chan bool)
	client.connected = make(chan bool)
	client.done = make(chan bool)
	client.error = make(chan error)
	client.Message = make(chan []byte)

	return client
}

func (c *Client) Connect(address string) (*Client, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return c, err
	}

	c.connection, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return c, err
	}

	err = c.connection.SetKeepAlive(true)
	if err != nil {
		return c, err
	}

	err = c.connection.SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		return c, err
	}

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
				//fmt.Println("unexpected data: %s", buf[:n])
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
