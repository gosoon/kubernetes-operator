#!/bin/bash

[ -d output ] || mkdir output

# etcd server cert/key
cfssl gencert --initca=true etcd-root-ca-csr.json | cfssljson -bare output/ca
hosts -component=etcd -csrfile=./etcd-csr.json | \
    cfssl gencert -ca=output/ca.pem -ca-key=output/ca-key.pem \
    -config=etcd-gencert.json etcd-csr.json | \
    cfssljson -bare output/etcd-server

# etcd peer cert/key
hosts -component=etcd -csrfile=./config-etcd-peer.json | \
cfssl gencert -ca=output/ca.pem -ca-key=output/ca-key.pem \
    -config=etcd-gencert.json \
    etcd-csr.json | \
    cfssljson -bare output/etcd-peer
