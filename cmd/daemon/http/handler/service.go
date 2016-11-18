package handler

import (
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	"k8s.io/client-go/1.5/pkg/api"
	"net/http"
)

func ServiceListH(w http.ResponseWriter, _ *http.Request) {

	var (
		er  error
		ctx = context.Get()
	)

	ctx.Log.Info("get nodes list")

	//ctx.K8S.LB().Accounts().Create()

	nodes, err := ctx.K8S.Core().Nodes().List(api.ListOptions{})
	if err != nil {
		ctx.Log.Panic(err.Error())
	}

	ctx.Log.Info(nodes)
	ctx.Log.Infof("There are %d pods in the cluster\n", len(nodes.Items))
	for _, node := range nodes.Items {
		ctx.Log.Infof("pod: %s: %s", node.Name, node.Status)
	}

	w.WriteHeader(200)
	_, er = w.Write([]byte(ctx.Info.Version))
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
