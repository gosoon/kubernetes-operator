#!/bin/bash

[ -d output ] || mkdir output

cfssl gencert --initca=true k8s-root-ca-csr.json | cfssljson -bare output/k8s-root-ca

# TODO: modified bin path 
for targetName in kube-apiserver kube-controller-manager kube-scheduler kube-proxy kubelet-client admin; do
    bin/hosts -component=kubernetes -csrfile=$targetName-csr.json | \
    bin/cfssl gencert -ca=output/k8s-root-ca.pem -ca-key=output/k8s-root-ca-key.pem -config=k8s-gencert.json -profile=kubernetes | \
    bin/cfssljson -bare output/$targetName  
done
