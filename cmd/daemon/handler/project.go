package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	c "github.com/lastbackend/lastbackend/cmd/daemon/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"io"
	"io/ioutil"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"net/http"
)

func ProjectListH(w http.ResponseWriter, _ *http.Request) {
	var ctx = c.Get()
	ctx.Log.Info("get projects list")

	projects, err := ctx.Storage.Project().GetByUser("test")
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	buf, er := json.Marshal(projects)
	if er != nil {
		ctx.Log.Error(er.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write(buf)
}

func ProjectInfoH(w http.ResponseWriter, r *http.Request) {
	var ctx = c.Get()
	ctx.Log.Info("get project info")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		e.User.AccessDenied().Http(w)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	session := s.(*model.Session)

	projects, err := ctx.Storage.Project().GetByID(session.Uid, id)
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	buf, er := json.Marshal(projects)
	if er != nil {
		ctx.Log.Error(er.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write(buf)
}

type projectCreateS struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *projectCreateS) decodeAndValidate(reader io.Reader) *e.Err {
	var ctx = c.Get()
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return e.Project.Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return e.Project.IncorrectJSON(err)
	}

	return nil
}

func ProjectCreateH(w http.ResponseWriter, r *http.Request) {
	var ctx = c.Get()
	ctx.Log.Info("create project")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		e.User.AccessDenied().Http(w)
		return
	}
	session := s.(*model.Session)

	// request body struct
	rq := new(projectCreateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error(err)
		err.Http(w)
		return
	}

	p := new(model.Project)
	p.Name = *rq.Name
	p.Description = *rq.Description
	p.User = session.Uid

	project, err := ctx.Storage.Project().Insert(p)
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	namespace := &v1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name:      project.ID,
			Namespace: project.ID,
			Labels: map[string]string{
				"user": session.Username,
			},
		},
	}

	_, er := ctx.K8S.Core().Namespaces().Create(namespace)
	if er != nil {
		ctx.Log.Error(er.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	buf, er := json.Marshal(project)
	if er != nil {
		ctx.Log.Error(er.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write(buf)
}

func ProjectDeleteH(w http.ResponseWriter, r *http.Request) {
	var ctx = c.Get()
	ctx.Log.Info("delete project")

	//params := mux.Vars(r)
	//id := params["id"]

	w.WriteHeader(200)
	w.Write([]byte{})
}
