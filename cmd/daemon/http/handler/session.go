package handler

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"io"
	"io/ioutil"
	"net/http"
)

type sessionCreateS struct {
	Login    *string `json:"login,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (s *sessionCreateS) decodeAndValidate(reader io.Reader) *e.Err {

	var err error
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return e.Session.Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return e.Session.IncorrectJSON(err)
	}

	if s.Login == nil || *s.Login == "" {
		return e.Session.BadParameter("login", err)
	}

	if s.Password == nil || *s.Password == "" {
		return e.Session.BadParameter("password", err)
	}

	return nil
}

// SessionCreateH - create session handler
func SessionCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		er  error
		err *e.Err
		ctx = context.Get()
	)

	ctx.Log.Debug("Create session handler")

	// request body struct
	rq := new(sessionCreateS)
	if err = rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error(err.Err())
		err.Http(w)
		return
	}

	user, err := ctx.Storage.User().GetByUsername(*rq.Login)
	if err == nil && user != nil {
		user, err := ctx.Storage.User().GetByEmail(*rq.Login)
		if err == nil && user != nil {
			err = e.User.NotFound()
		}
	}
	if err != nil {
		ctx.Log.Error(err.Err())
		err.Http(w)
		return
	}

	if err := user.ValidatePassword(*rq.Password); err != nil {
		e.HTTP.AccessDenied(w)
		return
	}

	sw := struct {
		Token string `json:"token"`
	}{}

	sw.Token, er = model.NewSession(user.ID, ``, user.Username, user.Email).Encode()
	if er != nil {
		ctx.Log.Error(er)
		e.HTTP.InternalServerError(w)
		return
	}

	response, er := json.Marshal(sw)
	if er != nil {
		ctx.Log.Error(er)
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
