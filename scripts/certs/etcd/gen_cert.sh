#!/bin/bash

[ -d output ] || mkdir output

# etcd server cert/key
cfssl gencert --initca=true etcd-root-ca-csr.json | cfssljson -bare output/ca

# etcd server
cfssl gencert \
  -ca=output/ca.pem \
  -ca-key=output/ca-key.pem \
  -config=etcd-gencert.json \
  -hostname=127.0.0.1,${ETCD_HOSTS} \
  etcd-csr.json | cfssljson -bare output/etcd-server

# etcd peer
cfssl gencert \
  -ca=output/ca.pem \
  -ca-key=output/ca-key.pem \
  -config=etcd-gencert.json \
  -hostname=127.0.0.1,${ETCD_HOSTS} \
  etcd-csr.json | cfssljson -bare output/etcd-peer
