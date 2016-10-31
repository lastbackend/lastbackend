package handler

import (
	c "github.com/gorilla/context"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"net/http"
)

func UserGetH(w http.ResponseWriter, r *http.Request) {

	var err *e.Err
	var ctx = context.Get()

	ctx.Log.Debug("Get user handler")

	s, ok := c.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error(e.StatusAccessDenied)
		e.HTTP.AccessDenied(w)
		return
	}

	session := s.(*model.Session)

	user, err := ctx.Adapter.User.Get(ctx.Storage, session.Username)
	if err != nil {
		ctx.Log.Error(err)
		err.Http(w)
		return
	}

	response, errjson := user.View().ToJson()
	if errjson != nil {
		ctx.Log.Error(errjson)
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
}
