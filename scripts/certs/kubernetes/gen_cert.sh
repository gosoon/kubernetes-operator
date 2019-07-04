#!/bin/bash

[ -d output ] || mkdir output

cfssl gencert --initca=true k8s-root-ca-csr.json | cfssljson -bare output/k8s-root-ca

for targetName in kube-apiserver kube-controller-manager kube-scheduler kube-proxy kubelet-client admin; do
    cfssl gencert -ca=output/k8s-root-ca.pem -ca-key=output/k8s-root-ca-key.pem -config=k8s-gencert.json -profile=kubernetes $targetName-csr.json | cfssljson -bare output/$targetName
done
