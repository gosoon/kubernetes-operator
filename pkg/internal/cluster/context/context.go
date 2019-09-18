package context

import (
	"fmt"
	"path/filepath"

	"github.com/gosoon/kubernetes-operator/pkg/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/internal/util/env"
)

// Context is the private shared context underlying pkg/cluster.Context
//
// NOTE: this is the internal one, it should contain reasonably trivial
// methods that are safe to share between various user facing methods
// pkg/cluster.Context is a superset of this, packages like create and delete
// consume this
type Context struct {
	name string
}

// NewContext returns a new internal cluster management context
// if name is "" the default name will be used
func NewContext(name string) *Context {
	if name == "" {
		name = constants.DefaultClusterName
	}
	return &Context{
		name: name,
	}
}

// Name returns the cluster's name
func (c *Context) Name() string {
	return c.name
}

// KubeConfigPath returns the path to where the Kubeconfig would be placed
// by kind based on the configuration.
func (c *Context) KubeConfigPath() string {
	// configDir matches the standard directory expected by kubectl etc
	configDir := filepath.Join(env.HomeDir(), ".kube")
	// note that the file name however does not, we do not want to overwrite
	// the standard config, though in the future we may (?) merge them
	fileName := fmt.Sprintf("config-%s", c.Name())
	return filepath.Join(configDir, fileName)
}
