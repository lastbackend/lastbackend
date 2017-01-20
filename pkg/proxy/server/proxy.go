package server

import (
	"fmt"
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"io"
	"k8s.io/client-go/pkg/api/v1"
	"net"
)

type Proxy struct {
	ctx        k8s.IK8S
	port       string
	authorized bool
	close      chan int
	Ready      chan int
	Done       chan int
}

func New(ctx k8s.IK8S) *Proxy {

	var proxy = new(Proxy)

	proxy.ctx = ctx

	proxy.close = make(chan int)

	proxy.Done = make(chan int)
	proxy.Ready = make(chan int)

	return proxy
}

func (p *Proxy) Start(port int) {

	server := NewTCPServer(port)
	server.Start()

	server.Accept(func(conn net.Conn) {
		fmt.Printf("New connection")
		conn.Write([]byte(`{"allow":true}`))

		fmt.Println("::1")

		var otps = &v1.PodAttachOptions{
			Container: "redis",
			Stdin:     true,
			Stdout:    false,
			Stderr:    false,
			TTY:       false,
		}

		req := p.ctx.LB().Pods("aaaca8b4-6198-491c-8bb4-edb8f1740945").Attach("lb-redis-4065565212-a79sb", otps)

		//req := p.ctx.LB().RESTClient().Post().
		//	Resource("pods").
		//	Name("lb-redis-4065565212-a79sb").
		//	Namespace("aaaca8b4-6198-491c-8bb4-edb8f1740945").
		//	SubResource("attach")
		//req.VersionedParams(&api.PodAttachOptions{
		//	Container: "redis",
		//	Stdin:     true,
		//	Stdout:    true,
		//	Stderr:    true,
		//	TTY:       true,
		//}, api.ParameterCodec)

		fmt.Println(req.URL())

		readCloser, err := req.Stream()
		if err != nil {
			fmt.Println(err)
			return
		}

		defer readCloser.Close()

		io.Copy(conn, readCloser)

		notify := make(chan error)

		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := readCloser.Read(buf)
				if err != nil {
					notify <- err
					return
				}
				fmt.Println("::4")
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
				//case <-time.After(time.Second * 1):
				//  fmt.Println("timeout 1, still alive")
			}
		}

		//defer readCloser.Close()

		//cl := NewTCPClient(p.host)
		//
		//client, err := cl.Connect("aaaca8b4-6198-491c-8bb4-edb8f1740945", "lb-redis-4065565212-a79sb")
		//if err != nil {
		//	return
		//}

		//p.copy(client.connection, conn)
		//p.copy(conn, client.connection)

		//go func() {
		//	select {
		//	case <-p.close:
		//		client.Close()
		//		server.Close()
		//		close(p.Done)
		//		return
		//	}
		//}()
		//
		//go func() {
		//	for {
		//		select {
		//		case msg := <-client.Message:
		//			fmt.Println("---------------")
		//			fmt.Println(string(msg))
		//			fmt.Println("---------------")
		//		}
		//	}
		//}()
	})
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
