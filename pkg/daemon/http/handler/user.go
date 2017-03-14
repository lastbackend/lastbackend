package handler

import (
	"github.com/gorilla/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"net/http"
)

func UserGetH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		ctx = c.Get()
	)

	ctx.Log.Debug("Get user handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		e.HTTP.Unauthorized(w)
		return
	}

	session := s.(*model.Session)

	user, err := ctx.Storage.User().GetByID(session.Uid)
	if err == nil && user == nil {
		e.New("user").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find user by id", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := user.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	ctx.Log.Info(string(response))

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}
