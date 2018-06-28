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

	"time"

	"github.com/gorilla/websocket"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
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

type Event struct {
	Entity string      `json:"entity"`
	Action string      `json:"action"`
	Name   string      `json:"name"`
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
		cm   = distribution.NewClusterModel(r.Context(), envs.Get().GetStorage())
		done = make(chan bool, 1)
	)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debugf("%s:subscribe:> set websocket upgrade err: %s", logPrefix, err.Error())
		return
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var serviceEvents = make(chan types.ServiceEvent)
	var namespaceEvents = make(chan types.NamespaceEvent)
	var clusterEvents = make(chan types.ClusterEvent)

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
				close(clusterEvents)
				return
			case e := <-clusterEvents:

				var data interface{}
				if e.Data == nil {
					data = nil
				} else {
					data = v1.View().Cluster().New(e.Data)
				}

				event := Event{
					Entity: "cluster",
					Action: e.Action,
					Name:   e.Name,
					Data:   data,
				}

				if err = conn.WriteJSON(event); err != nil {
					log.Errorf("%s:subscribe:> write cluster event to socket error.", logPrefix)
				}
			case e := <-serviceEvents:

				var data interface{}
				if e.Data == nil {
					data = nil
				} else {
					data = v1.View().Service().New(e.Data)
				}

				event := Event{
					Entity: "service",
					Action: e.Action,
					Name:   e.Name,
					Data:   data,
				}

				if err = conn.WriteJSON(event); err != nil {
					log.Errorf("%s:subscribe:> write service event to socket error.", logPrefix)
				}
			case e := <-namespaceEvents:

				var data interface{}
				if e.Data == nil {
					data = nil
				} else {
					data = v1.View().Namespace().New(e.Data)
				}

				event := Event{
					Entity: "namespace",
					Action: e.Action,
					Name:   e.Name,
					Data:   data,
				}

				if err = conn.WriteJSON(event); err != nil {
					log.Errorf("%s:subscribe:> write namespace event to socket error.", logPrefix)
				}
			}
		}
	}()

	go cm.Watch(clusterEvents)
	go sm.Watch(serviceEvents)
	go nm.Watch(namespaceEvents)

	go func() {
		for range ticker.C {
			if err := conn.WriteMessage(websocket.TextMessage, []byte{}); err != nil {
				log.Errorf("%s:subscribe:> writing to the client websocket err: %s", logPrefix, err.Error())
				done <- true
				break
			}
		}
	}()

	<-done
}
