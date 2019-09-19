package create

import (
	"time"

	internaltypes "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/types"
)

// ClusterOption is a cluster creation option
type ClusterOption func(*internaltypes.ClusterOptions) (*internaltypes.ClusterOptions, error)

// WithConfig configures creating the cluster using the config file at path
func WithConfig(path string) ClusterOption {
	return func(o *internaltypes.ClusterOptions) (*internaltypes.ClusterOptions, error) {
		// TODO :receiver operator ecs config
		//var err error
		//o.Config, err = internalencoding.Load(path)
		return o, nil
	}
}

// WaitForReady configures create to use interval as maximum wait time for the control plane node to be ready
func WaitForReady(interval time.Duration) ClusterOption {
	return func(o *internaltypes.ClusterOptions) (*internaltypes.ClusterOptions, error) {
		o.WaitForReady = interval
		return o, nil
	}
}
