#!/bin/bash

[ -d output ] || mkdir output

cfssl gencert --initca=true k8s-root-ca-csr.json | cfssljson -bare output/ca

for targetName in kube-apiserver kube-controller-manager kube-scheduler kubelet-api-admin admin kube-proxy kubelet-client; do
    hosts -component=kubernetes -csrfile=$targetName-csr.json | \
    cfssl gencert -ca=output/ca.pem -ca-key=output/ca-key.pem \
    -config=k8s-gencert.json -profile=kubernetes $targetName-csr.json | \
    cfssljson -bare output/$targetName
done
