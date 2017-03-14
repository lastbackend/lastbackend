package context

import (
	m "github.com/lastbackend/lastbackend/libs/adapter/storage/mock"
	"github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"github.com/lastbackend/lastbackend/libs/interface/log"
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	l "github.com/lastbackend/lastbackend/libs/log"
)

var context Context

func Get() *Context {
	return &context
}

func Mock() *Context {
	var err error

	context.Log = new(l.Log)
	context.Log.Init()
	context.Log.SetDebugLevel()
	context.Log.Disabled()

	if err != nil {
		context.Log.Panic(err)
	}

	// TODO: Need create mocks for k8s
	//context.K8S, err = k8s.Get(config.GetK8S())
	//if err != nil {
	//	context.Log.Panic(err)
	//}

	context.Storage, err = m.Get()
	if err != nil {
		context.Log.Panic(err)
	}

	return &context
}

type Context struct {
	Log              log.ILogger
	K8S              k8s.IK8S
	TemplateRegistry *http.RawReq
	Storage          storage.IStorage
}

