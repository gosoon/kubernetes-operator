package nodes

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/internal/apis/config"

	"github.com/pkg/errors"
)

// BootstrapControlPlaneNode returns a handle to the bootstrap control plane node
func BootstrapControlPlaneNode(allNodes []config.Node) (*config.Node, error) {
	controlPlaneNode := &config.Node{}
	for _, node := range allNodes {
		if node.Role == ecsv1.ControlPlaneRole {
			controlPlaneNode.IP = node.IP
			controlPlaneNode.Role = node.Role
		}
	}

	if controlPlaneNode == nil {
		return nil, errors.Errorf(
			"expected at least one %s node",
			ecsv1.ControlPlaneRole,
		)
	}
	return controlPlaneNode, nil
}
