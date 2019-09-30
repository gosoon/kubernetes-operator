package create

import (
	"errors"
	"strings"

	createtypes "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/types"
)

func validate(cluster *createtypes.ClusterOptions) error {
	errs := []string{}
	if cluster.Config == nil {
		errs = append(errs, "invalid config")
	}

	if cluster.NodeImage == "" {
		errs = append(errs, "invalid node image")

	}

	if cluster.NodeAddress == "" {
		errs = append(errs, "invalid local ip")
	}

	if cluster.Role == "" {
		errs = append(errs, "invalid role")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
