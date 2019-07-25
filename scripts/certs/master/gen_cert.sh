#!/bin/bash

[ -e ../../deploy/config.sh ] && . ../../deploy/config.sh || exit
[ -d output ] || mkdir output

cfssl gencert -initca ca-csr.json | cfssljson -bare output/ca

# apiserver
cfssl gencert \
  -ca=output/ca.pem \
  -ca-key=output/ca-key.pem \
  -config=ca-config.json \
  -hostname=10.250.0.1,${MASTER_HOSTS},${MASTER_VIP},127.0.0.1,kubernetes,kubernetes.default,kubernetes.default.svc \
  -profile=kubernetes \
  kube-apiserver-csr.json | cfssljson -bare output/kube-apiserver

# kubelet
for node in `echo ${NODE_HOSTS} | tr ',' ' '`;do
    sed -i -e "s#<nodeName>#${node}#g" kubelet-csr.json
    cfssl gencert \
      -ca=output/ca.pem \
      -ca-key=output/ca-key.pem \
      -config=ca-config.json \
      -hostname=${NODE_HOSTS} \
      -profile=kubernetes \
      kubelet-csr.json | cfssljson -bare output/kubelet
done

# other component
for component in kube-controller-manager kube-scheduler kube-proxy apiserver-kubelet-client admin service-account;do
    cfssl gencert \
      -ca=output/ca.pem \
      -ca-key=output/ca-key.pem \
      -config=ca-config.json \
      -profile=kubernetes \
      ${component}-csr.json | cfssljson -bare output/${component}
done
