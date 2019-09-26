package context

import (
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/constants"
	createtypes "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/types"
	"google.golang.org/grpc"
)

// Context is the private shared context underlying pkg/cluster.Context
//
// NOTE: this is the internal one, it should contain reasonably trivial
// methods that are safe to share between various user facing methods
// pkg/cluster.Context is a superset of this, packages like create and delete
// consume this
type Context struct {
	Name           string
	ClusterOptions *createtypes.ClusterOptions
	Server         *grpc.Server
	Port           string
}

// NewContext returns a new internal cluster management context
// if name is "" the default name will be used
func NewContext(name string, server *grpc.Server, port string) *Context {
	if name == "" {
		name = constants.DefaultClusterName
	}
	return &Context{
		Name:   name,
		Server: server,
		Port:   port,
	}
}
