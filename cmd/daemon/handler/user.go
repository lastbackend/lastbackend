package handler

import (
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/utils"
	"io"
	"io/ioutil"
	"k8s.io/client-go/1.5/pkg/api/v1"
	//"k8s.io/kubernetes/pkg/client/restclient"
	"net/http"
	"strings"
	"time"
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

type UserView struct {
	UUID         string    `json:"uuid,omitempty"`
	Username     string    `json:"username,omitempty"`
	Email        string    `json:"email,omitempty"`
	Gravatar     string    `json:"gravatar,omitempty"`
	Active       bool      `json:"active,omitempty"`
	Organization bool      `json:"organization,omitempty"`
	Balance      float64   `json:"balance,omitempty"`
	Created      time.Time `json:"created,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`
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

	ctx.Log.Info("generate password", password)

	conf := &v1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name:      *rq.Username,
			Namespace: *rq.Username,
			Labels: map[string]string{
				"name": "user",
			},
		},
	}

	_, err = ctx.K8S.Core().Namespaces().Create(conf)
	if err != nil {
		panic(err)
	}

	user := UserView{
		UUID:     utils.GetUUIDV4(),
		Username: *rq.Username,
		Email:    *rq.Email,
	}

	fmt.Printf("%#v", user)

	//rst := restclient.RESTClient{}
  //
	//fmt.Printf("%#v", rst)
  //
	//res := rst.Post().Resource("users").Body(user).Do()
  //
	//fmt.Printf("%#v", res)

	w.WriteHeader(200)
	w.Write([]byte(`{"status":"ok"}`))
}
