package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
)

type sessionCreateS struct {
	Login    *string `json:"login,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (s *sessionCreateS) decodeAndValidate(reader io.Reader) *e.Err {

	var err error
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return e.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return e.New("user").IncorrectJSON(err)
	}

	if s.Login == nil || *s.Login == "" {
		return e.New("user").BadParameter("login", err)
	}

	if s.Password == nil || *s.Password == "" {
		return e.New("user").BadParameter("password", err)
	}

	return nil
}

// SessionCreateH - create session handler
func SessionCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		ctx = context.Get()
	)

	ctx.Log.Debug("Create session handler")

	// request body struct
	rq := new(sessionCreateS)
	if er := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error(err.Error())
		er.Http(w)
		return
	}

	user, err := ctx.Storage.User().GetByUsername(*rq.Login)
	if err == nil && user == nil {
		user, err = ctx.Storage.User().GetByEmail(*rq.Login)
		if err == nil && user == nil {
			e.HTTP.Unauthorized(w)
			return
		}
	}
	if err != nil {
		ctx.Log.Error(err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	if err := user.ValidatePassword(*rq.Password); err != nil {
		e.HTTP.Unauthorized(w)
		return
	}

	sw := struct {
		Token string `json:"token"`
	}{}

	sw.Token, err = model.NewSession(user.ID, ``, user.Username, user.Email).Encode()
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	response, er := json.Marshal(sw)
	if er != nil {
		ctx.Log.Error(er)
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
