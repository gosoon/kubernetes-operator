package types

import (
	"time"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/internal/apis/config"
)

// ClusterOptions holds cluster creation options
// NOTE: this is only exported for usage by the parent package and the options
// package
// See ClusterOption instead
type ClusterOptions struct {
	Name                 string
	Config               *config.Cluster
	NodeImage            string
	WaitForReady         time.Duration
	SetupKubernetes      bool // if kind should setup kubernetes after creating nodes
	NodeAddress          string
	Role                 ecsv1.NodeRole
	ExternalLoadBalancer string
	KubeConfigPath       string
}
