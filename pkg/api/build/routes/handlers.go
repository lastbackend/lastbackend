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
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"net/http"
)

func BuildListH(w http.ResponseWriter, r *http.Request) {

	var (
		log     = context.Get().GetLogger()
		storage = context.Get().GetStorage()
	)

	log.Debug("Get boold list handler")

	// TODO: replace to valid image uuid
	var uuid string

	builds, err := storage.Build().ListByImage(r.Context(), uuid)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	buf, err := json.Marshal(builds)
	if err != nil {
		log.Error(err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(buf); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}
