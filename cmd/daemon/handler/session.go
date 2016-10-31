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

	var err *e.Err
	var ctx = context.Get()

	ctx.Log.Debug("Create session handler")

	// request body struct
	rq := new(sessionCreateS)
	if err = rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error(err.Err())
		err.Http(w)
		return
	}

	user, err := ctx.Adapter.User.Get(ctx.Storage, *rq.Login)
	if err != nil {
		ctx.Log.Error(err.Err())
		err.Http(w)
		return
	}

	acc, err := ctx.Adapter.Account.Get(ctx.Storage, user.Username)
	if err != nil {
		ctx.Log.Error(err.Err())
		err.Http(w)
		return
	}

	if err := acc.ValidatePassword(*rq.Password); err != nil {
		e.HTTP.AccessDenied(w)
		return
	}

	var errsesion error
	sw := new(SessionView)
	sw.Token, errsesion = model.NewSession(user.UUID, ``, user.Username, user.Email).Encode()
	if errsesion != nil {
		ctx.Log.Error(errsesion)
		e.HTTP.InternalServerError(w)
		return
	}

	response, errjson := json.Marshal(sw)
	if errjson != nil {
		ctx.Log.Error(errjson)
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
}
