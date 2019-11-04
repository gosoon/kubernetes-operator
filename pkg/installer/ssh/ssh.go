package ssh

import (
	clientset "github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

// Options is ssh installer must flags.
type Options struct {
	Kubeclientset              kubernetes.Interface
	KubernetesClusterClientset clientset.Interface
	Recorder                   record.EventRecorder
}

type installer struct {
	opt *Options
}

// NewSSHInstaller is new a ssh installer object.
func NewSSHInstaller(o *Options) *installer {
	return &installer{opt: o}
}
