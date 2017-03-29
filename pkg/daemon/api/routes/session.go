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
	"io"
	"io/ioutil"
	"net/http"

	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/api/types"
)

type sessionCreateS struct {
	Login    *string `json:"login,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (s *sessionCreateS) decodeAndValidate(reader io.Reader) *errors.Err {

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("user").IncorrectJSON(err)
	}

	if s.Login == nil || *s.Login == "" {
		return errors.New("user").BadParameter("login", err)
	}

	if s.Password == nil || *s.Password == "" {
		return errors.New("user").BadParameter("password", err)
	}

	return nil
}

// SessionCreateH - create session handler
func SessionCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		ctx = context.Get()
	)

	ctx.Log.Debug("Create session handler")

	// request body struct
	rq := new(sessionCreateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error(err)
		errors.New("incomming data invalid").Unknown().Http(w)
		return
	}

	user, err := ctx.Storage.User().GetByUsername(*rq.Login)
	if err == nil && user == nil {
		user, err = ctx.Storage.User().GetByEmail(*rq.Login)
		if err == nil && user == nil {
			errors.HTTP.Unauthorized(w)
			return
		}
	}
	if err != nil {
		ctx.Log.Error(err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if err := user.Security.Pass.ValidatePassword(*rq.Password); err != nil {
		errors.HTTP.Unauthorized(w)
		return
	}

	sw := struct {
		Token string `json:"token"`
	}{}

	sw.Token, err = types.NewSession(user.Username, user.Emails.GetDefault()).Encode()
	if err != nil {
		ctx.Log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, er := json.Marshal(sw)
	if er != nil {
		ctx.Log.Error(er)
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
