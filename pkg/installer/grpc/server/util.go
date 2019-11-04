package server

import installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"

// injectClusterConfig is set some config by server,eg:image registry and node role
// current only inject images registry
func (inst *installer) injectClusterConfig(cluster *installerv1.KubernetesClusterRequest) *installerv1.KubernetesClusterRequest {
	cluster.Spec.Cluster.ImagesRegistry = inst.opt.ImagesRegistry
	return cluster
}
