package handler

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/utils"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type accountCreateS struct {
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	// It is a struct for body data for account create route
	// Pointer is for data validating
}

func (s *accountCreateS) decodeAndValidate(reader io.Reader) *e.Err {

	var err error
	var ctx = context.Get()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return e.User.Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return e.User.IncorrectJSON(err)
	}

	if s.Email == nil {
		return e.User.BadParameter("email")
	}

	if !utils.IsEmail(*s.Email) {
		return e.User.BadParameter("email")
	}

	if s.Username == nil {
		return e.User.BadParameter("username")
	}

	if !utils.IsUsername(*s.Username) {
		return e.User.BadParameter("username")
	}

	if s.Password == nil {
		return e.User.BadParameter("password")
	}

	if !utils.IsPassword(*s.Password) {
		return e.User.BadParameter("password")
	}

	*s.Username = strings.ToLower(*s.Username)
	*s.Email = strings.ToLower(*s.Email)

	return nil
}

type SessionView struct {
	Token string `json:"token,omitempty"`
}

func AccountCreateH(w http.ResponseWriter, r *http.Request) {

	//var err *e.Err
	//var ctx = context.Get()

	//ctx.Log.Debug("Create account handler")
  //
	//// request body struct
	//rq := new(accountCreateS)
	//if err := rq.decodeAndValidate(r.Body); err != nil {
	//	ctx.Log.Error(err)
	//	err.Http(w)
	//	return
	//}
  //
	//salt, errsalt := utils.GenerateSalt(*rq.Password)
	//if errsalt != nil {
	//	ctx.Log.Error(errsalt)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}
  //
	//password, errpassword := utils.GeneratePassword(*rq.Password, salt)
	//if errpassword != nil {
	//	ctx.Log.Error(errpassword)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}
  //
	//gravatar := utils.GenerateGravatar(*rq.Email)
  //
	//accountID, err := ctx.Adapter.User.Insert(ctx.Storage, *rq.Username, *rq.Email, gravatar)
	//if err != nil {
	//	ctx.Log.Error(err)
	//	err.Http(w)
	//	return
	//}
  //
	//_, err = ctx.Adapter.Account.Insert(ctx.Storage, *rq.Username, *accountID, password, salt)
	//if err != nil {
	//	ctx.Log.Error(err)
	//	err.Http(w)
	//	return
	//}
  //
	//sw := new(SessionView)
	//var errencode error
	//sw.Token, errencode = model.NewSession(*accountID, ``, *rq.Username, *rq.Email).Encode()
	//if errencode != nil {
	//	ctx.Log.Error(errencode)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}
  //
	//response, errjson := json.Marshal(sw)
	//if errjson != nil {
	//	ctx.Log.Error(errjson)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	w.WriteHeader(200)
	w.Write([]byte{})
}

func AccountGetH(w http.ResponseWriter, r *http.Request) {

	//var err *e.Err
	//var ctx = context.Get()
  //
	//ctx.Log.Debug("Get account handler")
  //
	//s, ok := c.GetOk(r, `session`)
	//if !ok {
	//	ctx.Log.Error(e.StatusAccessDenied)
	//	e.HTTP.AccessDenied(w)
	//	return
	//}
  //
	//session := s.(*model.Session)
  //
	//account, err := ctx.Adapter.Account.GetByID(ctx.Storage, session.Uid)
	//if err != nil {
	//	ctx.Log.Error(err)
	//	err.Http(w)
	//	return
	//}
  //
	//response, errjson := account.View().ToJson()
	//if errjson != nil {
	//	ctx.Log.Error(errjson)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	w.WriteHeader(200)
	w.Write([]byte{})
}
