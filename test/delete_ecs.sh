#!/bin/bash

request_body=$(cat<<EOF
{
    "metadata": {
        "name":"test-cluster"
    },
    "spec": {
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
        "scaleDown": 0,
        "scaleUp": 0,
        "serviceCIDR": ""
    }
}
EOF
)

curl -s -XDELETE -d "${request_body}" \
    http://127.0.0.1:8080/api/v1/region/default/namespace/default/cluster/test-cluster/delete

