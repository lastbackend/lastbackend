package template

import (
	"encoding/json"
	"fmt"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"io/ioutil"
)

func Get(name, version string) (*model.Template, *e.Err) {

	var (
		er  error
		ctx = context.Get()
		httperr = new(e.Http)
		tpl = new(model.Template)
	)

	_, _, er = ctx.TemplateRegistry.
		GET(fmt.Sprintf("/template/%s/%s", name, version)).
		Request(tpl, httperr)
	if er != nil {
		return nil, e.New("template").Unknown(er)
	}

	if httperr.Code != 0 {
		switch httperr.Status {
		case "NOT_FOUND":
			return nil, nil
		default:
			return nil, e.New("template").Unknown(er)
		}
	}

	return nil, nil
}

func List() (*model.TemplateList, *e.Err) {

	var (
		er error
		ctx       = context.Get()
		templates = new(model.TemplateList)
	)

	_, resp, er := ctx.TemplateRegistry.GET("/template").Do()
	if er != nil {
		return nil, e.New("template").Unknown(er)
	}

	buf, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		return nil, e.New("template").Unknown(er)
	}

	er = json.Unmarshal(buf, templates)
	if er != nil {
		return nil, e.New("template").Unknown(er)
	}

	return templates, nil
}
