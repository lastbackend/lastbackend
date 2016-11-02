package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	c "github.com/lastbackend/lastbackend/cmd/daemon/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/utils"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type userCreate struct {
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	// It is a struct for body data for account create route
	// Pointer is for data validating
}

func (s *userCreate) decodeAndValidate(reader io.Reader) *e.Err {

	var err error
	var ctx = c.Get()

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

func UserCreateH(w http.ResponseWriter, r *http.Request) {

	var err *e.Err
	var ctx = c.Get()

	ctx.Log.Debug("Create user handler")

	// request body struct
	rq := new(userCreate)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error(err)
		err.Http(w)
		return
	}

	salt, errsalt := utils.GenerateSalt(*rq.Password)
	if errsalt != nil {
		ctx.Log.Error(errsalt)
		e.HTTP.InternalServerError(w)
		return
	}

	password, errpassword := utils.GeneratePassword(*rq.Password, salt)
	if errpassword != nil {
		ctx.Log.Error(errpassword)
		e.HTTP.InternalServerError(w)
		return
	}

	u := new(model.User)

	u.Username = *rq.Username
	u.Email = *rq.Email
	u.Gravatar = utils.GenerateGravatar(u.Email)
	u.Password = password
	u.Salt = salt

	user, err := ctx.Storage.User().GetByUsername(u.Username)
	if err == nil && user != nil {
		err = e.User.UsernameExists()
	}
	if err != nil {
		ctx.Log.Error(err.Err())
		err.Http(w)
		return
	}

	user, err = ctx.Storage.User().GetByEmail(u.Email)
	if err == nil && user != nil {
		err = e.User.EmailExists()
	}
	if err != nil {
		ctx.Log.Error(err.Err())
		err.Http(w)
		return
	}

	user, err = ctx.Storage.User().Insert(u)
	if err != nil {
		ctx.Log.Error(err.Err())
		err.Http(w)
		return
	}

	sw := struct {
		Token string `json:"token"`
	}{}

	var errencode error
	sw.Token, errencode = model.NewSession(user.ID, ``, user.Username, user.Email).Encode()
	if errencode != nil {
		ctx.Log.Error(errencode)
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

func UserGetH(w http.ResponseWriter, r *http.Request) {

	var err *e.Err
	var ctx = c.Get()

	ctx.Log.Debug("Get user handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error(e.StatusAccessDenied)
		e.HTTP.AccessDenied(w)
		return
	}

	session := s.(*model.Session)

	user, err := ctx.Storage.User().GetByID(session.Uid)
	if err != nil {
		ctx.Log.Error(err)
		err.Http(w)
		return
	}

	response, err := user.ToJson()
	if err != nil {
		ctx.Log.Error(err.Err())
		err.Http(w)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
}
