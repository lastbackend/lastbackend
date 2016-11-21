package projects

import (
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"fmt"
	"encoding/json"
	"github.com/lastbackend/lastbackend/libs/table"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/structures"
)

func List(ctx *context.Context) {
	var jdata structures.ProjList
	//Get(url string, json []byte, header string, headerType string) ([]byte, int)
	jbytes, status := httpClient.Get(config.Get().ProjectUrl, nil,
		"Authorization", "Bearer " /* + ctx.GetUserToken() */)
	if status != 200 {
		fmt.Println("something went wrong")
		return
	}
	json.Unmarshal(jbytes, jdata)

	var table_data [][]string
	var header []string = []string{"Title", "Description", "Created", "Updated"}

	for i := 0; i < len(jdata.Proj); i++ {
		var record []string
		record = append(record, jdata.Proj[i].Name)
		record = append(record, jdata.Proj[i].Description)
		record = append(record, jdata.Proj[i].Created)
		record = append(record, jdata.Proj[i].Updated)
		table_data = append(table_data, record)
	}
	table.PrintTable(header, table_data, []string{})
}
