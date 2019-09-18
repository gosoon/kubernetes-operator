package types

import (
	"time"

	"github.com/gosoon/kubernetes-operator/pkg/internal/apis/config"
)

// ClusterOptions holds cluster creation options
// NOTE: this is only exported for usage by the parent package and the options
// package
// See ClusterOption instead
type ClusterOptions struct {
	Config *config.Cluster
	// NodeImage overrides the nodes' images in Config if non-zero
	NodeImage string
	//Retain       bool
	WaitForReady time.Duration
	//TODO: Refactor this. It is a temporary solution for a phased breakdown of different
	//      operations, specifically create. see https://github.com/kubernetes-sigs/kind/issues/324
	SetupKubernetes bool // if kind should setup kubernetes after creating nodes
	LocalIP         string
	Role            string
}
