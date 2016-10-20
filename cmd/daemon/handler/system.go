package handler

import (
	"github.com/deployithq/deployit/cmd/daemon/context"
	"net/http"
	"k8s.io/client-go/1.5/pkg/api"
)

func SystemVersionH(w http.ResponseWriter, _ *http.Request) {
	var ctx = context.Get()

	ctx.Log.Info("get nodes list")
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
	w.Write([]byte(ctx.Version))
}
