package cluster

import (
	"github.com/gosoon/kubernetes-operator/pkg/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/cluster/create"
	internalcreate "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create"
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
func NewContext(name string) *Context {
	// wrap a new internal context
	return &Context{
		ic: internalcontext.NewContext(name),
	}
}

// Create provisions and starts a kubernetes-in-docker cluster
func (c *Context) Create(options ...create.ClusterOption) error {
	return internalcreate.Cluster(c.ic, options...)
}
