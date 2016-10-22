package context

import (
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"github.com/lastbackend/lastbackend/libs/interface/log"
)

var context Context

func Get() *Context {
	return &context
}

type Context struct {
	Info struct {
		Version    string
		ApiVersion string
	}
	Log log.ILogger
	K8S k8s.IK8S
	// Other info for HTTP handlers can be here, like user UUID
}
