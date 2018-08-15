//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package backend

import (
	"github.com/gorilla/websocket"
	"github.com/lastbackend/lastbackend/pkg/log"
	"net/http"
	"sync"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
  logLevel = 3
)

type Socket struct {
	sync.RWMutex

	token    string
	endpoint string

	conn *websocket.Conn

	write chan []byte
	read  chan []byte
	ping  chan []byte
	pong  chan []byte

	err   chan error
	close chan error

	online chan bool
	dial   chan int

	end chan error

	attempt int
}

func (s *Socket) manage() {
	for {
		select {
		case t := <-s.dial:

			// Check if need reconnect timer
			if t > 0 {
				timer := time.NewTimer(time.Second * time.Duration(t))
				select {
				case <-timer.C:
				}
				timer.Stop()
			}
			// Call connect
			if err := s.connect(); err != nil {
				log.Errorf("Socket connect err: %s", err)
				return
			}
		}
	}
}

func (s *Socket) reconnect() {
	s.conn.Close()
	s.dial <- 1
}

func (s *Socket) connect() error {
	var (
		err  error
		resp *http.Response
	)

	s.conn, resp, err = websocket.DefaultDialer.Dial(s.endpoint, nil)
	if err != nil {
		if err == websocket.ErrBadHandshake {
			log.V(6).Errorf("handshake failed with status %d", resp.StatusCode)
		} else {
			log.V(6).Errorf("Socket: stream dial error: %s", err)
		}

		if resp != nil {
			if resp.StatusCode == http.StatusNotFound {
				log.V(6).Error("Socket: error: stream not found")
				return err
			}
		}

		s.attempt++

		if s.attempt >= 5 {
			s.dial <- 5
		} else {
			s.dial <- s.attempt * 1
		}

		return nil
	}
	s.attempt = 0

	// start connection heartbeat

	s.listen()

	return nil
}

func (s *Socket) listen() {

	// Create listener to pipe message to hub
	pipe := make(chan msg)
	go func() {
		for {
			select {
			case p := <-pipe:

				s.conn.SetWriteDeadline(time.Now().Add(writeWait))

				s.Lock()
				err := s.conn.WriteMessage(p.MT, p.MD)
				s.Unlock()

				if err != nil {
					log.V(6).Errorf("Socket: stream: write message error: %s", err)
					s.reconnect()
					return
				}
			}
		}
	}()

	go func() {

		for {

			select {

			case m := <-s.write:
				pipe <- msg{websocket.TextMessage, m}

			case p := <-s.ping:
				pipe <- msg{websocket.PingMessage, p}

			case p := <-s.pong:
				pipe <- msg{websocket.PongMessage, p}

			case <-s.read:
				log.V(6).Debug("Socket: Incoming messaging system not implemented yet")

			case err := <-s.close:
				pipe <- msg{websocket.CloseMessage, []byte{}}
				log.V(6).Errorf("Socket: send: %v", err)
				s.disconnect()
				return
			}
		}
	}()

	s.conn.SetPongHandler(func(string) error {
		s.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go func() {
		for {

			m, _, err := s.conn.ReadMessage()

			if m == websocket.PingMessage {
				s.pong <- []byte{}
				continue
			}

			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					log.V(logLevel).Debug("Web socket closed by client:", err)
					s.end <- nil
					return
				}
			}

			if err != nil {
				log.Errorf("Unexpected socket error: %s", err)
				s.end <- err
				break
			}
		}

	}()

}

func (s *Socket) send(data []byte) error {
	s.write <- data
	return nil
}

func (s *Socket) disconnect() {

}

func (s *Socket) Disconnect() {
	s.disconnect()
}

func (s *Socket) End() error {
	select {
	case e := <-s.end:
		log.V(logLevel).Debug("Socket ended")
		return e
	}
}

func (s *Socket) Write(chunk []byte) {
	s.write <- chunk
}

func NewSocketBackend(endpoint string) StreamBackend {

	var s = new(Socket)

	s.ping = make(chan []byte)
	s.pong = make(chan []byte)
	s.write = make(chan []byte)
	s.read = make(chan []byte)

	s.err = make(chan error)
	s.close = make(chan error)

	s.online = make(chan bool)
	s.dial = make(chan int)

	s.end = make(chan error)
	s.attempt = 0

	s.endpoint = endpoint

	go s.manage()

	s.dial <- 0

	return s
}
