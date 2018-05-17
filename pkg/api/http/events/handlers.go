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

package events

import (
	"net/http"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/gorilla/websocket"
	"time"
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

type Event struct {
	Entity string      `json:"entity"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
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
		sm   = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		nm   = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		done = make(chan bool, 1)
	)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debugf("%s:subscribe:> set websocket upgrade err: %s", logPrefix, err.Error())
		return
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := conn.WriteMessage(websocket.TextMessage, []byte{}); err != nil {
			log.Errorf("%s:subscribe:> writing to the client websocket err: %s", logPrefix, err.Error())
			done <- true
			break
		}
	}

	var serviceEvents = make(chan *types.Event)
	var namespaceEvents = make(chan *types.Event)

	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		log.Debugf("%s:subscribe:> HTTP connection just closed.", logPrefix)
		done <- true
	}()

	go func() {
		for {
			select {
			case <-done:
				close(serviceEvents)
				close(namespaceEvents)
				return
			case e := <-serviceEvents:

				if e.Data == nil {
					continue
				}

				event := Event{
					Entity: "service",
					Action: e.Action,
					Data:   v1.View().Service().New(e.Data.(*types.Service)),
				}

				if err = conn.WriteJSON(event); err != nil {
					log.Errorf("%s:subscribe:> write service event to socket error.", logPrefix)
				}
			case e := <-namespaceEvents:

				if e.Data == nil {
					continue
				}

				event := Event{
					Entity: "namespace",
					Action: e.Action,
					Data:   v1.View().Namespace().New(e.Data.(*types.Namespace)),
				}

				if err = conn.WriteJSON(event); err != nil {
					log.Errorf("%s:subscribe:> write namespace event to socket error.", logPrefix)
				}
			}
		}
	}()

	go sm.Watch(serviceEvents)
	go nm.Watch(namespaceEvents)

	<-done
}
