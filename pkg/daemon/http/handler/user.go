package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
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
		return e.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return e.New("user").IncorrectJSON(err)
	}

	if s.Email == nil {
		return e.New("user").BadParameter("email")
	}

	if !validator.IsEmail(*s.Email) {
		return e.New("user").BadParameter("email")
	}

	if s.Username == nil {
		return e.New("user").BadParameter("username")
	}

	if !validator.IsUsername(*s.Username) {
		return e.New("user").BadParameter("username")
	}

	if s.Password == nil {
		return e.New("user").BadParameter("password")
	}

	if !validator.IsPassword(*s.Password) {
		return e.New("user").BadParameter("password")
	}

	*s.Username = strings.ToLower(*s.Username)
	*s.Email = strings.ToLower(*s.Email)

	return nil
}

func UserCreateH(w http.ResponseWriter, r *http.Request) {

	var err error
	var ctx = c.Get()

	ctx.Log.Debug("Create user handler")

	// request body struct
	rq := new(userCreate)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		err.Http(w)
		return
	}

	salt, er := generator.GenerateSalt(*rq.Password)
	if er != nil {
		ctx.Log.Error("Error: generate salt", er)
		e.HTTP.InternalServerError(w)
		return
	}

	password, er := generator.GeneratePassword(*rq.Password, salt)
	if er != nil {
		ctx.Log.Error("Error: generate password", er)
		e.HTTP.InternalServerError(w)
		return
	}

	u := new(model.User)

	u.Username = *rq.Username
	u.Email = *rq.Email
	u.Gravatar = generator.GenerateGravatar(u.Email)
	u.Password = password
	u.Salt = salt

	user, err := ctx.Storage.User().GetByUsername(u.Username)
	if err == nil && user != nil {
		e.New("user").NotUnique("username").Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find user by username", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	user, err = ctx.Storage.User().GetByEmail(u.Email)
	if err == nil && user != nil {
		e.New("user").NotUnique("email").Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find user by email", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	user, err = ctx.Storage.User().Insert(u)
	if err != nil {
		ctx.Log.Error("Error: insert user to db", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	sw := struct {
		Token string `json:"token"`
	}{}

	sw.Token, err = model.NewSession(user.ID, ``, user.Username, user.Email).Encode()
	if err != nil {
		ctx.Log.Error("Error: create session token", err)
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := json.Marshal(sw)
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err)
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

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

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}
