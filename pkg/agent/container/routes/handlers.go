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
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/agent/container"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"net/http"
)

func GetLogsH(w http.ResponseWriter, r *http.Request) {
	var (
		log      = context.Get().GetLogger()
		cid      = mux.Vars(r)["container"]
		notify   = w.(http.CloseNotifier).CloseNotify()
		doneChan = make(chan bool, 1)
	)

	log.Debug("Get container logs")

	go func() {
		<-notify
		log.Debug("HTTP connection just closed.")
		doneChan <- true
	}()

	if err := container.Logs(cid, true, w, doneChan); err != nil {
		log.Errorf("Error: get container logs err %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
}
