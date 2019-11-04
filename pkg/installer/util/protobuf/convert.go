package protobuf

import (
	"encoding/json"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterConvertToTypes is convert protobuf to kubernetesCluster golang types
// typeMeta and objectMeta is a nested struct in KubernetesCluster struct,
// and proto not support struct custom json tag and embed struct,
// so typeMeta and objectMeta is assigned separately and spec is assigned by json unmarshal
func ClusterConvertToTypes(clusterRequest *installerv1.KubernetesClusterRequest) (*ecsv1.KubernetesCluster, error) {
	typeMeta := metav1.TypeMeta{
		Kind:       clusterRequest.Kind,
		APIVersion: clusterRequest.APIVersion,
	}
	objectMeta := metav1.ObjectMeta{
		Name:        clusterRequest.Name,
		Namespace:   clusterRequest.Namespace,
		Labels:      clusterRequest.Labels,
		Annotations: clusterRequest.Annotations,
		Finalizers:  clusterRequest.Finalizers,
	}

	cluster := &ecsv1.KubernetesCluster{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
	}

	bytes, err := json.Marshal(clusterRequest)
	if err != nil {
		//glog.Error("marshal kubernetes cluster request failed with:", err)
		return cluster, err
	}

	err = json.Unmarshal(bytes, cluster)
	if err != nil {
		//glog.Error("unmarshal kubernetes cluster request bytes failed with:", err)
		return cluster, err
	}

	return cluster, nil
}

// ClusterConvertToProtobuf is convert kubernetesCluster to protobuf
func ClusterConvertToProtobuf(cluster *ecsv1.KubernetesCluster) (*installerv1.KubernetesClusterRequest, error) {
	typeMeta := installerv1.TypeMeta{
		Kind:       cluster.Kind,
		APIVersion: cluster.APIVersion,
	}
	objectMeta := installerv1.ObjectMeta{
		Name:        cluster.Name,
		Namespace:   cluster.Namespace,
		Labels:      cluster.Labels,
		Annotations: cluster.Annotations,
		Finalizers:  cluster.Finalizers,
	}

	clusterRequest := &installerv1.KubernetesClusterRequest{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
	}

	bytes, err := json.Marshal(cluster)
	if err != nil {
		//glog.Error("marshal kubernetes cluster types failed with:", err)
		return clusterRequest, err
	}

	err = json.Unmarshal(bytes, clusterRequest)
	if err != nil {
		//glog.Error("unmarshal kubernetes cluster bytes failed with:", err)
		return clusterRequest, err
	}

	return clusterRequest, nil
}
