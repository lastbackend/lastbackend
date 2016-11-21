package projects

import (
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/structures"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/errors"
	"fmt"
	"encoding/json"
)

type CreateTemplate struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func Create(p_name string, desc string, ctx *context.Context) {
	jData := []byte("{ \"name\": \"" + p_name + "\", \"description\": \"" + desc + "\" }")
	fmt.Println("\n\nP_NAME: ", p_name)
	fmt.Println("DESC: ", desc, "\n\n")
	resp, status := httpClient.Post("http://localhost:3000/project", jData,
		"Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDg3NDE3ODk1LCJqdGkiOjE0ODc0MTc4OTUsIm9pZCI6IiIsInVpZCI6IjU2MmYwY2EwLTI2ZWEtNGFiNC1hZDBmLTU1N2NmYjJmYjgwNyIsInVzZXIiOiJtb2NrZWQifQ.VjHgKRqJCwf7TDphHPHhMl6njwL7agE1dzPVeGy5HFI")
	if errors.Process(status) {
		return
	}
	var proj structures.Project
	json.Unmarshal(resp, &proj)
	fmt.Println("PROJECT ID: ", proj.Id)

}