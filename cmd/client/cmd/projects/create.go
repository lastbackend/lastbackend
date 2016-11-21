package projects

import (
	"github.com/lastbackend/lastbackend/cmd/client/context"
	//"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/structures"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"github.com/lastbackend/lastbackend/cmd/client/config"
)

type CreateTemplate struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func Create(p_name string, desc string, ctx *context.Context) {
	/*
	fmt.Println("DESCRIPTION: ", *desc)
	var project = new(structures.Project)
	local_time := time.Now().String()
	project.Name, project.Description, project.Created, project.Updated = *p_name, *desc, local_time, local_time
	if err != nil {
		fmt.Println(err)
	}
	*/
	jData := []byte("{ \"name\": \"" + p_name + "\", \"description\": \"" + desc + "\" }")
	/*resp, status := */httpClient.Post(config.Get().ProjectUrl, jData,
		"Authorization", "Bearer " /* + ctx.GetUserToken() */)

}