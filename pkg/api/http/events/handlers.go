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
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/sockets"
	"net/http"
)

const (
	logLevel      = 2
	defaultClient = "lastbackend"
)

func EventSubscribeH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Event: subscribe on events")

	var (
		err error
		hub = envs.Get().GetWssHub()
	)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conn, err := sockets.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Event: upgrade socker err: %s", err)
		return
	}

	log.V(logLevel).Debug("Handler: Event: new websocket connection")

	client := hub.NewConnection(defaultClient, conn)

	log.V(logLevel).Debug("Handler: Event: connection ready to receive data")

	client.WritePump()
}
