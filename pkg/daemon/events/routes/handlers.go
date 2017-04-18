//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package routes

import (
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/wss"
	"net/http"
)

func EventSubscribeH(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		log = context.Get().GetLogger()
		hub = context.Get().GetWssHub()
	)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conn, err := wss.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debug("New websockets connection")
	client := hub.NewConnection("lastbackend", conn)
	log.Debug("Websockets connection ready to receive data")
	client.WritePump()
}
