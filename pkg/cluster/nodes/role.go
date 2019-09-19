package nodes

import (
	"sort"
	"strings"

	"github.com/gosoon/kubernetes-operator/pkg/cluster/constants"
	"github.com/pkg/errors"
)

// SelectNodesByRole returns a list of nodes with the matching role
func SelectNodesByRole(allNodes []Node, role string) []Node {
	out := []Node{}
	for _, node := range allNodes {
		if node.Role == role {
			out = append(out, node)
		}
	}
	return out
}

// ExternalLoadBalancerNode returns a node handle for the external control plane
// loadbalancer node or nil if there isn't one
func ExternalLoadBalancerNode(allNodes []Node) (*Node, error) {
	// identify and validate external load balancer node
	loadBalancerNodes := SelectNodesByRole(
		allNodes,
		constants.ExternalLoadBalancerNodeRoleValue,
	)
	if len(loadBalancerNodes) < 1 {
		return nil, nil
	}
	if len(loadBalancerNodes) > 1 {
		return nil, errors.Errorf(
			"unexpected number of %s nodes %d",
			constants.ExternalLoadBalancerNodeRoleValue,
			len(loadBalancerNodes),
		)
	}
	return &loadBalancerNodes[0], nil
}

// ControlPlaneNodes returns all control plane nodes such that the first entry
// is the bootstrap control plane node
func ControlPlaneNodes(allNodes []Node) ([]Node, error) {
	controlPlaneNodes := SelectNodesByRole(
		allNodes,
		constants.ControlPlaneNodeRoleValue,
	)
	// pick the first by sorting
	// TODO(bentheelder): perhaps in the future we should mark this node
	// specially at container creation time
	sort.Slice(controlPlaneNodes, func(i, j int) bool {
		return strings.Compare(controlPlaneNodes[i].IP, controlPlaneNodes[j].IP) < 0
	})
	return controlPlaneNodes, nil
}

// BootstrapControlPlaneNode returns a handle to the bootstrap control plane node
func BootstrapControlPlaneNode(allNodes []Node) (*Node, error) {
	controlPlaneNodes, err := ControlPlaneNodes(allNodes)
	if err != nil {
		return nil, err
	}
	if len(controlPlaneNodes) < 1 {
		return nil, errors.Errorf(
			"expected at least one %s node",
			constants.ControlPlaneNodeRoleValue,
		)
	}
	return &controlPlaneNodes[0], nil
}

// SecondaryControlPlaneNodes returns handles to the secondary
// control plane nodes and NOT the bootstrap control plane node
func SecondaryControlPlaneNodes(allNodes []Node) ([]Node, error) {
	controlPlaneNodes, err := ControlPlaneNodes(allNodes)
	if err != nil {
		return nil, err
	}
	if len(controlPlaneNodes) < 1 {
		return nil, errors.Errorf(
			"expected at least one %s node",
			constants.ControlPlaneNodeRoleValue,
		)
	}
	return controlPlaneNodes[1:], nil
}
