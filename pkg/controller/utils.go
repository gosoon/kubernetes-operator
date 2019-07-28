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
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"

	mapset "github.com/deckarep/golang-set"
)

func convertNodesToString(nodes []ecsv1.Node) string {
	var nodeStr string
	l := len(nodes)
	for i, node := range nodes {
		nodeStr += node.IP
		if i != l-1 {
			nodeStr += " "
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
