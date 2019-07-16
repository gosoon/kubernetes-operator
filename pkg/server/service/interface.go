package service

import (
	clientset "github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned"
	"github.com/gosoon/kubernetes-operator/pkg/types"
)

type Options struct {
	KubernetesClusterClientset clientset.Interface
}

type service struct {
	opt *Options
}

func New(opt *Options) Interface {
	return &service{opt: opt}
}

type Interface interface {
	// cluster
	CreateCluster(namespace string, name string, clusterInfo *types.EcsClient) error
	DeleteCluster(namespace string, name string, clusterInfo *types.EcsClient) error

	// scale
	ScaleUp(namespace string, name string, clusterInfo *types.EcsClient) error
	ScaleDown(namespace string, name string, clusterInfo *types.EcsClient) error

	// callback
	CreateClusterCallback(namespace string, name string, result *types.CallBack) error
	ScaleUpCallback(namespace string, name string, result *types.CallBack) error
	ScaleDownCallback(namespace string, name string, result *types.CallBack) error
	DeleteClusterCallback(namespace string, name string, result *types.CallBack) error
}
