package service

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	clientset "github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned"
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
	CreateCluster(namespace string, name string, kubernetesCluster *ecsv1.KubernetesCluster) error
	DeleteCluster(namespace string, name string) error
}
