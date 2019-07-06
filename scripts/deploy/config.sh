#!/bin/bash

# docker version
DOCKER_VER="docker-ce-18.06.1.ce"

# k8s version
KUBERNETES_VER="kubernetes_v1.14.0"
KUBERNETES_DOWNLOAD_URL="https://dl.k8s.io/v1.14.0/kubernetes.tar.gz"

# etcd version
ETCD_VER="etcd_v3.3.13"

# calico version 
CALICO_VER="v3.7"

# coredns version 
COREDNS_VER="v1.4.0"

# etcd hosts, eg : "10.0.2.15,10.0.2.16"   
ETCD_HOSTS="10.0.2.15"  

# master hosts, eg : "10.0.2.15,10.0.2.16"
MASTER_HOSTS="10.0.2.15"  

# host ip
LOCAL_IP=$(ip route get 1 | awk '{print $NF;exit}')

# deploy home  
DEPLOY_HOME_DIR="/home/kube"

# export certs bin
CERTS_BIN_DIR="kubernetes-operator/scripts/bin/certs"
export PATH=$PATH:${DEPLOY_HOME_DIR}/${CERTS_BIN_DIR}
