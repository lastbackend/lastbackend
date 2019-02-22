//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package events

import (
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/socket"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	logLevel  = 2
	logPrefix = "api:handler:event"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//EventSubscribeH - realtime subscribe handler
func EventSubscribeH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:subscribe:> subscribe on subscribe", logPrefix)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.V(logLevel).Debugf("%s:subscribe:> watch all events", logPrefix)

	var (
		done  = make(chan bool, 1)
		leave = make(chan *socket.Socket)
		event = make(chan *socket.Message)
	)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.V(logLevel).Debugf("%s:subscribe:> set websocket upgrade err: %s", logPrefix, err.Error())
		return
	}

	skt := socket.NewSocket(r.Context(), conn, leave, event)

	var (
		es = make(chan *types.Event)
	)

	go func() {
		for {
			select {

			case <-leave:
				done <- true
				return

			case e := <-es:

				event := v1.View().Event().New(e)
				msg, err := event.ToJson()
				if err != nil {
					log.Errorf("err: %s", err.Error())
					continue
				}

				skt.Write(msg)
			}
		}
	}()

	envs.Get().GetMonitor().Subscribe(es, done)
}
