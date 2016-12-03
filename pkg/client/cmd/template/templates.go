package template

import (
	"errors"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	tab "github.com/crackcomm/go-clitable"
)

func ViewTemplatesCmd() {
	ctx := context.Get()

	err := ViewTemplates()
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func ViewTemplates() error {

	var res = make(map[string][]string)
	req_err := new(e.Http)

	ctx := context.Mock()
	_, _, err := ctx.HTTP.
		GET("/template").
		Request(&res, req_err)

	if err != nil {
		return err
	}

	if req_err.Code != 0 {
		return errors.New(e.Message(req_err.Status))
	}

	table := tab.New([]string{"Name", "Version"})
	keys := make([]string, 0, len(res))
	for k := range res {
		keys = append(keys, k)
	}
	for i := 0; i < len(res); i++ {

		table.AddRow(map[string]interface{}{
			"Name":    keys[i],
			"Version": res[keys[i]][0],
		})
		table.Markdown = true
		for ii := 1; ii < len(res[keys[i]]); ii++ {
			table.AddRow(map[string]interface{}{
				"Name":    " ",
				"Version": res[keys[i]][ii],
			})
			table.Markdown = true
		}
	}
	table.Print()

	return err
}