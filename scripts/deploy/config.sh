#!/bin/bash

# this section is glabel vars,do not edit
# deploy home  
DEPLOY_HOME_DIR="/home/kubernetes-operator"

# export certs bin
CERTS_BIN_DIR="scripts/bin/certs"
export PATH=$PATH:${DEPLOY_HOME_DIR}/${CERTS_BIN_DIR}

# ----------------
# this section is define by users
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

# if in docker and use env
if ! grep docker /proc/1/cgroup -qa; then
    [ -f hosts_env ] && source hosts_env || exit 1
fi
