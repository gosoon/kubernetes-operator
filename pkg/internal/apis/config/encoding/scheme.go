/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package encoding

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/internal/apis/config"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/kubeadm"
)

// Ecsv1ToInternalCluster is convert ecsv1.KubernetesCluster to internal cluster config
func Ecsv1ToInternalCluster(cluster *ecsv1.KubernetesCluster, nodeAddress string) *config.Cluster {
	out := &config.Cluster{
		ExternalLoadBalancer: cluster.Spec.Cluster.ExternalLoadBalancer,
		Networking: config.Networking{
			APIServerPort:    kubeadm.APIServerPort,
			APIServerAddress: nodeAddress,
			PodSubnet:        cluster.Spec.Cluster.PodCIDR,
			ServiceSubnet:    cluster.Spec.Cluster.ServiceCIDR,
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
