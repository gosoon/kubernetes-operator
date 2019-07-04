#!/bin/bash

request_body=$((cat <<EOF
{
	"name": "集群名称",
	"timeout_mins": 集群创建超时时间,
	"cluster_type": "集群类型，Kubernetes",
	"container_cidr": "容器POD CIDR",
	"service_cidr": "服务CIDR",
	"num_of_nodes": "Worker节点数",
}
EOF
))

curl -s -XPOST -d "$request_body" -H 'Authorization: Bearer xxxx' \
    http://xxxx/kube/api/v2/resource/region/hxy01/scene/all/namespace/com

