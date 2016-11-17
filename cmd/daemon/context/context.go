package context

import (
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"github.com/lastbackend/lastbackend/libs/interface/log"
	"github.com/lastbackend/lastbackend/libs/interface/storage"

	m "github.com/lastbackend/lastbackend/libs/adapter/storage/mock"
	l "github.com/lastbackend/lastbackend/libs/log"
)

var context Context

func Get() *Context {
	return &context
}

func Mock() *Context {

	var err error

	context.Info.Version = "0.3.0"
	context.Info.ApiVersion = "0.3.0"

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
	Info struct {
		Version    string
		ApiVersion string
	}
	Log     log.ILogger
	K8S     k8s.IK8S
	Storage storage.IStorage
}
