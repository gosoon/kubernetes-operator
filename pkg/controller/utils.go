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

package controller

import (
	"bytes"
	"encoding/json"
	"time"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/gosoon/kubernetes-operator/pkg/types"
	"sigs.k8s.io/yaml"

	mapset "github.com/deckarep/golang-set"
)

// Returns 0 for resyncPeriod in case resyncing is not needed.
func NoResyncPeriodFunc() time.Duration {
	return 0
}

func convertNodesToString(nodes []ecsv1.Node) string {
	var nodeStr string
	l := len(nodes)
	for i, node := range nodes {
		nodeStr += node.IP
		if i != l-1 {
			nodeStr += ","
		}
	}
	return nodeStr
}

// diffNode is diff oldNodeList and newNodeList
func diffNodeList(oldNodeList []ecsv1.Node, newNodeList []ecsv1.Node, operation string) []ecsv1.Node {
	oldNodeListSet := mapset.NewSet()
	for _, node := range oldNodeList {
		oldNodeListSet.Add(node.IP)
	}

	newNodeListSet := mapset.NewSet()
	for _, node := range newNodeList {
		newNodeListSet.Add(node.IP)
	}

	var diff mapset.Set
	if operation == enum.KubeScalingUp {
		diff = newNodeListSet.Difference(oldNodeListSet)
	} else {
		diff = oldNodeListSet.Difference(newNodeListSet)
	}
	return setToEcsV1Node(diff)
}

func setToEcsV1Node(set mapset.Set) []ecsv1.Node {
	var nodeList []ecsv1.Node
	for t := range set.Iterator().C {
		ip := t.(string)
		node := ecsv1.Node{IP: ip}
		nodeList = append(nodeList, node)
	}
	return nodeList
}

// hosts yaml by ansible used,it is like:
/*
all:
  children:
    calico-rr:
      hosts: {}
    etcd:
      hosts: {}
    kube-master:
      hosts:
        192.168.1.1:
    kube-node:
      hosts: {}
  hosts:
    192.168.1.1:
      access_ip: 192.168.1.1
      ansible_host: 192.168.1.1
      ip: 192.168.1.1
*/
func compressHostsYAML(cluster *ecsv1.KubernetesCluster) string {
	masterList := map[string]*types.Host{}
	nodeList := map[string]*types.Host{}
	etcdList := map[string]*types.Host{}
	hostList := map[string]*types.Host{}

	for _, node := range cluster.Spec.Cluster.MasterList {
		ip := node.IP
		masterList[ip] = nil
		hostList[ip] = &types.Host{
			AnsibleHost: ip,
			IP:          ip,
			AccessIP:    ip,
		}
	}

	for _, node := range cluster.Spec.Cluster.NodeList {
		ip := node.IP
		nodeList[ip] = nil
		hostList[ip] = &types.Host{
			AnsibleHost: ip,
			IP:          ip,
			AccessIP:    ip,
		}
	}

	for _, node := range cluster.Spec.Cluster.EtcdList {
		ip := node.IP
		etcdList[ip] = nil
		hostList[ip] = &types.Host{
			AnsibleHost: ip,
			IP:          ip,
			AccessIP:    ip,
		}
	}

	config := types.HostsYamlFormat{
		All: types.AllHosts{
			Hosts: hostList,
			Children: types.Children{
				KubeMaster: map[string]map[string]*types.Host{"hosts": masterList},
				KubeNode:   map[string]map[string]*types.Host{"hosts": nodeList},
				Etcd:       map[string]map[string]*types.Host{"hosts": etcdList},
				Calico:     map[string]map[string]*types.Host{"hosts": map[string]*types.Host{}},
			},
		},
	}

	var configBytes []byte
	configJSON, _ := json.MarshalIndent(config, "", "")

	configBytes, _ = yaml.JSONToYAML(configJSON)

	configBytes = bytes.Replace(configBytes, []byte("null"), []byte(""), -1)

	return string(configBytes)
}
