package projects

import (
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"encoding/json"
	"github.com/lastbackend/lastbackend/libs/table"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/structures"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/errors"
)

func List(ctx *context.Context) {
	var jdata structures.ProjList

	jbytes, status := httpClient.Get("http://localhost:3000/project", nil,
		"Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDg3NDE3ODk1LCJqdGkiOjE0ODc0MTc4OTUsIm9pZCI6IiIsInVpZCI6IjU2MmYwY2EwLTI2ZWEtNGFiNC1hZDBmLTU1N2NmYjJmYjgwNyIsInVzZXIiOiJtb2NrZWQifQ.VjHgKRqJCwf7TDphHPHhMl6njwL7agE1dzPVeGy5HFI")
	if errors.Process(status) {
		return
	}
	json.Unmarshal(jbytes, &jdata)
	var header []string = []string{"ID", "Name", "Created", "Updated"}
	var data [][]string
	for i := 0; i < len(jdata); i++ {
		d := []string{jdata[i].Id, jdata[i].Name, jdata[i].Created, jdata[i].Updated}
		data = append(data, d)
	}
	table.PrintTable(header, data, []string{})

}
