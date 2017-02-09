package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	e "github.com/lastbackend/lastbackend/libs/errors"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/template"
)

type deployTemplateS struct {
	Project *string `json:"project,omitempty"`
	Target  *string `json:"target,omitempty"`
}

func (d *deployTemplateS) decodeAndValidate(reader io.Reader) *e.Err {

	var (
		err error
		ctx = c.Get()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return e.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, d)
	if err != nil {
		return e.New("service").IncorrectJSON(err)
	}

	if d.Project == nil {
		return e.New("service").BadParameter("project")
	}

	if d.Target == nil {
		return e.New("service").BadParameter("target")
	}

	return nil
}

func TemplateListH(w http.ResponseWriter, _ *http.Request) {

	var (
		er             error
		ctx            = c.Get()
		response_empty = func() {
			w.WriteHeader(http.StatusOK)
			_, er = w.Write([]byte("[]"))
			if er != nil {
				ctx.Log.Error("Error: write response", er.Error())
				return
			}
			return
		}
	)

	templates, err := template.List()
	if err != nil {
		ctx.Log.Error(err.Error())
		response_empty()
		return
	}

	if templates == nil {
		response_empty()
		return
	}

	response, err := templates.ToJson()
	if er != nil {
		ctx.Log.Error(err.Error())
		response_empty()
		return
	}

	w.WriteHeader(http.StatusOK)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
