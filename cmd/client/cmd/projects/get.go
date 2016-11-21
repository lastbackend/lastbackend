package projects

import (
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"encoding/json"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/structures"
	"github.com/lastbackend/lastbackend/libs/table"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/errors"
	"fmt"
)


func Get(p_name string, ctx *context.Context) {
	var jdata structures.Project
	fmt.Println("ENTER PROJECT ID: ")
	var id string
	fmt.Scan(&id)
	fmt.Println("YOUR ID: ", id)
	jbytes, status := httpClient.Get("http://localhost:3000/project/" + id, nil,
		"Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDg3NDE3ODk1LCJqdGkiOjE0ODc0MTc4OTUsIm9pZCI6IiIsInVpZCI6IjU2MmYwY2EwLTI2ZWEtNGFiNC1hZDBmLTU1N2NmYjJmYjgwNyIsInVzZXIiOiJtb2NrZWQifQ.VjHgKRqJCwf7TDphHPHhMl6njwL7agE1dzPVeGy5HFI")
	if errors.Process(status) {
		return
	}
	err := json.Unmarshal(jbytes, &jdata)
	if err != nil {
		return
	}
	var header []string = []string{"ID", "Name", "Created", "Updated"}
	var data [][]string
	d := []string{jdata.Id, jdata.Name, jdata.Created, jdata.Updated}

	data = append(data, d)
	d = d[:0]

	table.PrintTable(header, data, []string{})
}