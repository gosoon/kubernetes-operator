package cluster

import (
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/create"
	internalcontext "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/context"
	internalcreate "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create"
	"google.golang.org/grpc"
)

// DefaultName is the default cluster name
const DefaultName = constants.DefaultClusterName

// Context is used to create / manipulate kubernetes-in-docker clusters
// See: NewContext()
type Context struct {
	// the internal context type, shared between implementations of more
	// advanced methods like create
	ic *internalcontext.Context
}

// NewContext returns a new cluster management context
// if name is "" the default name will be used (constants.DefaultClusterName)
func NewContext(name string, server *grpc.Server, port string) *Context {
	// wrap a new internal context
	return &Context{
		ic: internalcontext.NewContext(name, server, port),
	}
}

// Create provisions and starts a kubernetes-in-docker cluster
func (c *Context) Create(options ...create.ClusterOption) error {
	return internalcreate.Cluster(c.ic, options...)
}
