package protobuf

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"
)

// ClusterConvertToTypes is convert protobuf to kubernetesCluster golang types
func ClusterConvertToTypes(proto *installerv1.KubernetesClusterRequest) *ecsv1.KubernetesCluster {
	return &ecsv1.KubernetesCluster{}
}

// ClusterConvertToProtobuf is convert kubernetesCluster to protobuf
func ClusterConvertToProtobuf(cluster *ecsv1.KubernetesCluster) *installerv1.KubernetesClusterRequest {
	return &installerv1.KubernetesClusterRequest{TypeMeta: &installerv1.TypeMeta{Kind: "v1"}}
}
