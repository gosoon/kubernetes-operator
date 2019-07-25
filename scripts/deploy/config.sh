#!/bin/bash

# this section is glabel vars,do not edit
# deploy home  
DEPLOY_HOME_DIR="/home"

# export certs bin
CERTS_BIN_DIR="kubernetes-operator/scripts/bin/certs"
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

# ----------------
# this section is define by users,if in docker do not update hosts list
host_list() {
    # etcd hosts, eg : "10.0.2.15,10.0.2.16"   
    ETCD_HOSTS="10.0.4.15"  

    # master hosts, eg : "10.0.2.15,10.0.2.16"
    # if not use vip,please assign the one master ip
    MASTER_HOSTS="10.0.4.15"  
    MASTER_VIP="10.0.4.15"     

    # node hosts, eg : "10.0.2.15,10.0.2.16"
    NODE_HOSTS="10.0.4.15"  

    # host ip
    LOCAL_IP="10.0.4.15"
}

# if in docker and use env
if ! grep docker /proc/1/cgroup -qa; then
    host_list
fi
