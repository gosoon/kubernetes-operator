#!/bin/bash

region="default"
name="test-cluster"
request_body=$(cat<<EOF
{
  "name": "test-cluster",
  "clusterType": "kubernetes",
  "masterList": [
    {
      "ip": "192.168.1.10"
    },
    {
      "ip": "192.168.1.11"
    }
  ],
  "nodeList": [
    {
      "ip": "192.168.1.12"
    },
    {
      "ip": "192.168.1.12"
    }
  ],
  "etcdList": [
    {
      "ip": "192.168.1.12"
    },
    {
      "ip": "192.168.1.12"
    }
  ],
  "privateSSHKey": "",
  "serviceCIDR": "",
  "retry": false
}
EOF
)

curl -s -XPOST -d "${request_body}" \
    http://127.0.0.1:8080/api/v1/region/${region}/cluster/${name}/scaledown

