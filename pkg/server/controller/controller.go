package controller

import (
	clientset "github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned"
	"github.com/gosoon/kubernetes-operator/pkg/server/service"

	"github.com/gorilla/mux"
)

type Options struct {
	KubernetesClusterClientset clientset.Interface
	Service                    service.Interface
}

type Controller interface {
	Register(router *mux.Router)
}
