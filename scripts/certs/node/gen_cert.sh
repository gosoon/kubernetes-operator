#!/bin/bash

[ -d output ] || mkdir output

cfssl gencert --initca=true k8s-root-ca-csr.json | cfssljson -bare output/ca

for targetName in kube-proxy kubelet-client; do
    cat $targetName-csr.json | \
    cfssl gencert -ca=output/ca.pem -ca-key=output/ca-key.pem \
    -config=k8s-gencert.json -profile=kubernetes $targetName-csr.json | \
    cfssljson -bare output/$targetName
done
