package controller

import (
	clientset "github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned"
	"github.com/gosoon/kubernetes-operator/pkg/server/service"
	"k8s.io/client-go/kubernetes"

	"github.com/gorilla/mux"
)

type Options struct {
	KubernetesClusterClientset clientset.Interface
	KubeClientset              kubernetes.Interface
	Service                    service.Interface
}

type Controller interface {
	Register(router *mux.Router)
}
