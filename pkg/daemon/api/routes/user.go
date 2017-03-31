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
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/daemon/api/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"net/http"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/gorilla/context"
)

func UserGetH(w http.ResponseWriter, r *http.Request) {

	var (
		ctx = c.Get()
		session *types.Session
	)

	ctx.Log.Debug("Get user handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	session = s.(*types.Session)

	user, err := ctx.Storage.User().GetByUsername(session.Username)
	if err == nil && user == nil {
		errors.New("user").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find user by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewUser(user).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}
