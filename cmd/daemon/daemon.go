package daemon

import (
	"encoding/json"
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/cmd/daemon/config"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	"github.com/lastbackend/lastbackend/libs/adapter/k8s"
	"github.com/lastbackend/lastbackend/libs/adapter/storage"
	"github.com/lastbackend/lastbackend/libs/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/1.5/pkg/api"
	"os"
	"os/signal"
	"syscall"
)

func Run(cmd *cli.Cmd) {
	var err error

	var ctx = context.Get()
	var cfg = config.Get()

	cmd.Spec = "[-c][-d]"

	var debug = cmd.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "Enable debug mode"})
	var configPath = cmd.String(cli.StringOpt{Name: "c config", Value: "./config.yaml", Desc: "Path to config file", HideValue: true})

	cmd.Before = func() {

		ctx.Log = new(log.Log)
		ctx.Log.Init()

		if *configPath != "" {

			// Parsing config file
			configBytes, err := ioutil.ReadFile(*configPath)
			if err != nil {
				ctx.Log.Panic(err)
			}

			err = yaml.Unmarshal(configBytes, cfg)
			if err != nil {
				ctx.Log.Panic(err)
			}
		}

		if *debug {
			cfg.Debug = *debug
			ctx.Log.SetDebugLevel()
			ctx.Log.Info("Logger debug mode enabled")
		}

		// Initializing database
		ctx.Log.Info("Initializing daemon")
		ctx.K8S, err = k8s.Get(config.GetK8S())
		if err != nil {
			ctx.Log.Panic(err)
		}

		ctx.Storage, err = storage.Get()
		if err != nil {
			ctx.Log.Panic(err)
		}

		if cfg.HttpServer.Port == 0 {
			cfg.HttpServer.Port = 3000
		}
	}

	cmd.Action = func() {

		CoreNodes := ctx.K8S.Core().Nodes()
		CoreNodesList, _ := CoreNodes.List(api.ListOptions{})
		buf, _ := json.Marshal(CoreNodesList)
		ctx.Log.Info(">> CoreNodes: ", string(buf))

		CorePods := ctx.K8S.Core().Pods("unloop")
		CorePodsList, _ := CorePods.List(api.ListOptions{})
		buf, _ = json.Marshal(CorePodsList)
		ctx.Log.Info(">> CorePods: ", string(buf))

		CorePodTemplates := ctx.K8S.Core().PodTemplates("unloop")
		CorePodTemplatesList, _ := CorePodTemplates.List(api.ListOptions{})
		buf, _ = json.Marshal(CorePodTemplatesList)
		ctx.Log.Info(">> CorePodTemplates: ", string(buf))

		CoreEndpoints := ctx.K8S.Core().Endpoints("unloop")
		CoreEndpointsList, _ := CoreEndpoints.List(api.ListOptions{})
		buf, _ = json.Marshal(CoreEndpointsList)
		ctx.Log.Info(">> CoreEndpoints: ", string(buf))

		CoreReplicationControllers := ctx.K8S.Core().ReplicationControllers("unloop")
		CoreReplicationControllersList, _ := CoreReplicationControllers.List(api.ListOptions{})
		buf, _ = json.Marshal(CoreReplicationControllersList)
		ctx.Log.Info(">> CoreReplicationControllers: ", string(buf))

		CoreServices := ctx.K8S.Core().Services("unloop")
		CoreServicesList, _ := CoreServices.List(api.ListOptions{})
		buf, _ = json.Marshal(CoreServicesList)
		ctx.Log.Info(">> CoreServices: ", string(buf))

		BatchJobs := ctx.K8S.Batch().Jobs("unloop")
		BatchJobsList, _ := BatchJobs.List(api.ListOptions{})
		buf, _ = json.Marshal(BatchJobsList)
		ctx.Log.Info(">> BatchJobs: ", string(buf))

		ExtensionsDeployments := ctx.K8S.Extensions().Deployments("unloop")
		ExtensionsDeploymentsList, _ := ExtensionsDeployments.List(api.ListOptions{})
		buf, _ = json.Marshal(ExtensionsDeploymentsList)
		ctx.Log.Info(">> ExtensionsDeployments: ", string(buf))

		ExtensionsJobs := ctx.K8S.Extensions().Jobs("unloop")
		ExtensionsJobsList, _ := ExtensionsJobs.List(api.ListOptions{})
		buf, _ = json.Marshal(ExtensionsJobsList)
		ctx.Log.Info(">> ExtensionsJobs: ", string(buf))

		ExtensionsReplicaSets := ctx.K8S.Extensions().ReplicaSets("unloop")
		ExtensionsReplicaSetsList, _ := ExtensionsReplicaSets.List(api.ListOptions{})
		buf, _ = json.Marshal(ExtensionsReplicaSetsList)
		ctx.Log.Info(">> ExtensionsReplicaSets: ", string(buf))

		go RunHttpServer(NewRouter(), cfg.HttpServer.Port)

		// Handle SIGINT and SIGTERM.
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		ctx.Log.Debug(<-ch)

		ctx.Log.Info("Handle SIGINT and SIGTERM.")
	}
}
