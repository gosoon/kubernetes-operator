#!/bin/bash

[ -e ../../deploy/config.sh ] && . ../../deploy/config.sh || exit

CERTS_DIR="/etc/kubernetes/ssl"
[ -d output ] || mkdir output

# kubelet
cfssl gencert \
  -ca=${CERTS_DIR}/ca.pem \
  -ca-key=${CERTS_DIR}/ca-key.pem \
  -config=ca-config.json \
  -hostname=${NODE_HOSTNAME},${NODE_HOSTS} \
  -profile=kubernetes \
  kubelet-csr.json | cfssljson -bare output/kubelet

# other component
cfssl gencert \
  -ca=${CERTS_DIR}/ca.pem \
  -ca-key=${CERTS_DIR}/ca-key.pem \
  -config=ca-config.json \
  -profile=kubernetes \
  kube-proxy-csr.json | cfssljson -bare output/kube-proxy
