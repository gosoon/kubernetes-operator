package controller

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
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
