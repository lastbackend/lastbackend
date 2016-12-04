package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/deployer"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type deployS struct {
	Project *string `json:"project,omitempty"`
	Target  *string `json:"target,omitempty"`
}

func (d *deployS) decodeAndValidate(reader io.Reader) *e.Err {

	var (
		err error
		ctx = c.Get()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return e.User.Unknown(err)
	}

	err = json.Unmarshal(body, d)
	if err != nil {
		return e.IncorrectJSON(err)
	}

	if d.Project == nil {
		return e.BadParameter("project")
	}

	if d.Target == nil {
		return e.BadParameter("target")
	}

	return nil
}

func DeployH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		ctx     = c.Get()
		session *model.Session
	)

	ctx.Log.Debug("Deploy handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.User.AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	// request body struct
	rq := new(deployS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err.Err())
		err.Http(w)
		return
	}

	if validator.IsGitUrl(*rq.Target) {
		// TODO: deploy from git repo url
		ctx.Log.Error("Error: not implement")
		e.HTTP.NotImplemented(w)
		return
	}

	// Below this template deployment

	parts := strings.Split(*rq.Target, ":")

	var name = parts[0]
	var version = "latest"

	if len(parts) == 2 {
		version = parts[1]
	}

	var httperr = new(e.Http)
	var tpl = new(model.Template)

	_, _, er = ctx.TemplateRegistry.
		GET(fmt.Sprintf("/template/%s/%s", name, version)).
		Request(tpl, httperr)

	ctx.Log.Debug("RESPONSE: GET TEMPALTE", er, httperr)
	if er != nil {
		ctx.Log.Error("Error: get tempalte", er.Error())
		e.HTTP.InternalServerError(w)
		return
	}
	if httperr != nil && httperr.Code != 0 {
		// TODO: Switch errors for client
		ctx.Log.Error("Error: get tempalte", httperr.Message)
		e.HTTP.InternalServerError(w)
		return
	}

	d := deployer.Get()
	err := d.DeployFromTemplate(session.Uid, *rq.Project, *tpl)
	if err != nil {
		ctx.Log.Error("Error: deploy from tempalte", err.Err())
		err.Http(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
