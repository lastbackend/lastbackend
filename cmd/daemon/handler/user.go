package handler

import (
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/account"
	"github.com/lastbackend/lastbackend/pkg/user"
	"github.com/lastbackend/lastbackend/utils"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type userCreateS struct {
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	// It is a struct for body data for user create route
	// Pointer is for data validating
}

func (s *userCreateS) decodeAndValidate(reader io.Reader) *e.Err {

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

func UserCreateH(w http.ResponseWriter, r *http.Request) {

	var err error
	// var cfg = config.Get()
	var ctx = context.Get()

	// request body struct
	rq := new(userCreateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error(err)
		err.Http(w)
		return
	}

	salt, err := utils.GenerateSalt(*rq.Password)
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	password, err := utils.GeneratePassword(*rq.Password, salt)
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	u1, err := user.Create(*rq.Username, *rq.Email)
	if err != nil {
		ctx.Log.Error(":: 1 >> ", err)
		e.HTTP.InternalServerError(w)
		return
	}

	fmt.Sprintf("1: %#v", u1)

	u2, err := user.Get(*rq.Username)
	if err != nil {
		ctx.Log.Error(":: 2 >> ", err)
		e.HTTP.InternalServerError(w)
		return
	}

	fmt.Sprintf("2: %#v", u2)

	session, err := account.Create(u2.UUID, password)
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := json.Marshal(session)
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
}

func UserGetH(w http.ResponseWriter, r *http.Request) {

	//var err error
	//var ctx = context.Get()
	//
	//s, ok := c.GetOk(r, `session`)
	//if !ok {
	//	ctx.Log.Error(e.StatusAccessDenied)
	//	e.HTTP.AccessDenied(w)
	//	return
	//}

	//session := s.(*model.Session)

	//namespace, err := ctx.K8S.Core().Namespaces().Get(session.Username)
	//if err != nil {
	//panic(err)
	//}
	//
	//fmt.Printf("%#v", namespace)
	//
	//var client = &v.CoreClient{}
	//res := client.Get().Resource("users")
	//fmt.Printf("%#v", res)

	w.WriteHeader(200)
	w.Write([]byte(`{"status":"ok"}`))
}
