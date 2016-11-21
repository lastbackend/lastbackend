package projects

import (
	"github.com/lastbackend/lastbackend/cmd/client/context"
	//"github.com/lastbackend/lastbackend/cmd/client/cmd/projects/structures"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"github.com/lastbackend/lastbackend/cmd/client/config"
)

func Remove(p_id string, ctx *context.Context) {
	//httpClient.Post(config.Get().ProjectUrl, jData,
	//	"Authorization", "Bearer " /* + ctx.GetUserToken() */)
	httpClient.Delete(config.Get().ProjectUrl, "Authorization", "Bearer " /* + ctx.GetUserToken() */)
}
