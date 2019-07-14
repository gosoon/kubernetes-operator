#!/bin/bash

[ -d output ] || mkdir output

#cfssl gencert --initca=true k8s-root-ca-csr.json | cfssljson -bare output/ca
ca="/etc/kubernetes/ssl/"


for targetName in kube-proxy kubelet-client; do
    cat $targetName-csr.json | \
    cfssl gencert -ca=${ca}/ca.pem -ca-key=${ca}/ca-key.pem \
    -config=k8s-gencert.json -profile=kubernetes $targetName-csr.json | \
    cfssljson -bare output/$targetName
done
