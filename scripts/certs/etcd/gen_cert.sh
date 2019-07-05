#!/bin/bash

[ -d output ] || mkdir output

cfssl gencert --initca=true etcd-root-ca-csr.json | cfssljson -bare output/etcd-root-ca 
bin/hosts -component=etcd -csrfile=./etcd-csr.json | \
cfssl gencert -ca=output/etcd-root-ca.pem -ca-key=output/etcd-root-ca-key.pem -config=etcd-gencert.json etcd-csr.json | \
cfssljson -bare output/etcd
