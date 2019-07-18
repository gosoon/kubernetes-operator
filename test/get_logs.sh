#!/bin/bash

region="default"
name="test-cluster"

curl -s -XGET -d "${request_body}" \
    http://127.0.0.1:8080/api/v1/region/${region}/cluster/${name}/logs
