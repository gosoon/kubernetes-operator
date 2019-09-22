package actions

import (
	createtypes "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/types"
	"github.com/gosoon/kubernetes-operator/pkg/internal/util/cli"
	"google.golang.org/grpc"
)

type Action interface {
	Execute(ctx *ActionContext) error
}

// ActionContext is data supplied to all actions
type ActionContext struct {
	Status  *cli.Status
	Cluster *createtypes.ClusterOptions
	Server  *grpc.Server
	Port    string
}

func NewActionContext(
	cluster *createtypes.ClusterOptions,
	server *grpc.Server,
	port string,
	status *cli.Status,
) *ActionContext {
	return &ActionContext{
		Status:  status,
		Cluster: cluster,
		Server:  server,
		Port:    port,
	}
}
