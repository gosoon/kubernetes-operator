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
	CreateCluster(region string, namespace string, name string, clusterInfo *types.EcsClient) error
	DeleteCluster(region string, namespace string, name string, clusterInfo *types.EcsClient) error

	// scale
	ScaleUp(region string, namespace string, name string, clusterInfo *types.EcsClient) error
	ScaleDown(region string, namespace string, name string, clusterInfo *types.EcsClient) error

	// callback
	CreateClusterCallback(region string, namespace string, name string, result *types.CallBack) error
	ScaleUpCallback(region string, namespace string, name string, result *types.CallBack) error
	ScaleDownCallback(region string, namespace string, name string, result *types.CallBack) error
	DeleteClusterCallback(region string, namespace string, name string, result *types.CallBack) error
}
