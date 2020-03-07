//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/lastbackend/lastbackend/internal/master/server/middleware"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/internal/util/socket"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:event"
)

// Handler represent the http handler for event
type Handler struct {
}

// NewEventHandler will initialize the event resources endpoint
func NewEventHandler(r *mux.Router, mw middleware.Middleware) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init event routes", logPrefix)

	handler := &Handler{
	}

	r.Handle("/events", h.Handle(handler.EventSubscribeH)).Methods(http.MethodGet)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// EventSubscribeH - realtime subscribe handler
func (handler Handler) EventSubscribeH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Debugf("%s:subscribe:> subscribe on subscribe", logPrefix)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Debugf("%s:subscribe:> watch all events", logPrefix)

	var (
		done  = make(chan bool, 1)
		leave = make(chan *socket.Socket)
		event = make(chan *socket.Message)
	)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debugf("%s:subscribe:> set websocket upgrade err: %s", logPrefix, err.Error())
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
}
