package projects

import (
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"fmt"
	"encoding/json"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/structures"
	"github.com/lastbackend/lastbackend/libs/table"
	"github.com/lastbackend/lastbackend/cmd/client/context"
)

func Get(p_name string, ctx *context.Context) {
	var jdata structures.Project
	//Get(url string, json []byte, header string, headerType string) ([]byte, int)
	jbytes, status := httpClient.Get(config.Get().ProjectUrl + "/" + p_name, nil,
		"Authorization", "Bearer " /* + ctx.GetUserToken() */)
	if status != 200 {
		fmt.Println("something went wrong")
		return
	}
	json.Unmarshal(jbytes, jdata)

	var header []string = []string{"Title", "Description", "Created", "Updated"}
	var record []string
	record = append(record, jdata.Name)
	record = append(record, jdata.Description)
	record = append(record, jdata.Created)
	record = append(record, jdata.Updated)
	var data [][]string
	data = append(data, record, []string{})
	table.PrintTable(header, data, []string{})
}