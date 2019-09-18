package actions

import (
	"github.com/gosoon/kubernetes-operator/pkg/cluster/nodes"
	"github.com/gosoon/kubernetes-operator/pkg/internal/apis/config"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/context"
	"github.com/gosoon/kubernetes-operator/pkg/internal/util/cli"
)

type Action interface {
	Execute(ctx *ActionContext) error
}

// ActionContext is data supplied to all actions
type ActionContext struct {
	Status         *cli.Status
	Config         *config.Cluster
	ClusterContext *context.Context
	//cache          *cachedData
	LocalIP string
	Role    string
}

func NewActionContext(
	cfg *config.Cluster,
	ctx *context.Context,
	status *cli.Status,
	localIP string,
	role string,
) *ActionContext {
	return &ActionContext{
		Status:         status,
		Config:         cfg,
		ClusterContext: ctx,
		LocalIP:        localIP,
		Role:           role,
	}
}

// Nodes returns the list of cluster nodes, this is a cached call
func (ac *ActionContext) Nodes() []nodes.Node {
	allNodes := ac.Config.Nodes
	var res []nodes.Node
	for _, node := range allNodes {
		res = append(res, nodes.Node{
			IP:    node.IP,
			Role:  string(node.Role),
			Image: node.Image,
		})
	}
	return res
}
