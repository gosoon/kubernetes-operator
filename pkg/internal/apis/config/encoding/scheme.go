package encoding

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/internal/apis/config"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/kubeadm"
)

// Ecsv1ToInternalCluster is convert ecsv1.KubernetesCluster to internal cluster config
func Ecsv1ToInternalCluster(cluster *ecsv1.KubernetesCluster) *config.Cluster {
	out := &config.Cluster{
		ExternalLoadBalancer: cluster.Spec.Cluster.ExternalLoadBalancer,
		Networking: config.Networking{
			APIServerPort: kubeadm.APIServerPort,
			//APIServerAddress: "", // local IP
			PodSubnet:     cluster.Spec.Cluster.PodCIDR,
			ServiceSubnet: cluster.Spec.Cluster.ServiceCIDR,
		},
		KubeVersion: cluster.Spec.Cluster.KubeVersion,
	}

	var clusterNodeList []ecsv1.Node
	clusterNodeList = append(clusterNodeList, cluster.Spec.Cluster.NodeList...)
	clusterNodeList = append(clusterNodeList, cluster.Spec.Cluster.MasterList...)

	var workerList []config.Node
	for _, node := range clusterNodeList {
		workerList = append(workerList, config.Node{
			IP:   node.IP,
			Role: node.Role,
		})
	}
	out.Nodes = workerList
	return out
}
