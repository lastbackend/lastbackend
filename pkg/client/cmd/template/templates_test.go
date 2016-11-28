package template

/*
	res["mock template 1"] = string{"first ver.", "last ver."}
	res["mock template 2"] = string{"first ver.","ver 0.0", "last ver."}
 */

import (
"github.com/lastbackend/lastbackend/pkg/client/context"
e "github.com/lastbackend/lastbackend/libs/errors"
tab "github.com/crackcomm/go-clitable"
)



func ViewTemplatesTest() error {

	var res = make(map[string][]string)
	req_err := new(e.Http)

	ctx := context.Mock()
	ctx.HTTP.
		GET("/jumpstart").
		Request(&res, req_err)

	res["mock template 1"] = string{"first ver.", "last ver."}
	res["mock template 2"] = string{"first ver.","ver 0.0", "last ver."}

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

	return nil
}