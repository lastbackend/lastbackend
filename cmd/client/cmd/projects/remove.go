package projects

import (
	"github.com/lastbackend/lastbackend/cmd/client/context"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/errors"
	"fmt"
)

func Remove(p_id string, ctx *context.Context) {
	fmt.Println("ENTER PROJECT ID: ")
	var id string
	fmt.Scan(&id)
	fmt.Println("YOUR ID: ", id)
	status := httpClient.Delete("http://localhost:3000/project/" + id, "Authorization",
		"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDg3NDE3ODk1LCJqdGkiOjE0ODc0MTc4OTUsIm9pZCI6IiIsInVpZCI6IjU2MmYwY2EwLTI2ZWEtNGFiNC1hZDBmLTU1N2NmYjJmYjgwNyIsInVzZXIiOiJtb2NrZWQifQ.VjHgKRqJCwf7TDphHPHhMl6njwL7agE1dzPVeGy5HFI")
	if errors.Process(status) {
		return
	}
}
