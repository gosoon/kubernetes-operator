package agent

import (
	"sort"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
)

type nodeList []ecsv1.Node

func (n nodeList) Less(i, j int) bool { return n[i].IP < n[j].IP }
func (n nodeList) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n nodeList) Len() int           { return len(n) }

// injectClusterNodeRole is set role for cluster all nodes
// when installing a cluster, kubeadm first install a controlplane,then secondarycontrolplane should join controlplane,
// the workers node also join controlplane
// ps:
// this action is implementation in server at first,because proto3 not support type define alias,the NodeRole is a NodeRole type not string
// in order to reduce complexity so can do it in agent
func injectClusterNodeRole(cluster *ecsv1.KubernetesCluster) *ecsv1.KubernetesCluster {
	// set master role
	// grpc server select some node and set to ControlPlaneRole,SecondaryControlPlaneRole,WorkerRole
	if len(cluster.Spec.Cluster.MasterList) > 0 {
		masterList := cluster.Spec.Cluster.MasterList
		sort.Sort(nodeList(masterList))
		for idx, master := range cluster.Spec.Cluster.MasterList {
			cluster.Spec.Cluster.MasterList[idx].Role = ecsv1.SecondaryControlPlaneRole
			if master.IP == masterList[0].IP {
				cluster.Spec.Cluster.MasterList[idx].Role = ecsv1.ControlPlaneRole
			}
		}
	}

	// set node role,default all node is worker
	if len(cluster.Spec.Cluster.NodeList) > 0 {
		for idx := range cluster.Spec.Cluster.NodeList {
			cluster.Spec.Cluster.NodeList[idx].Role = ecsv1.WorkerRole
		}
	}

	return cluster
}
