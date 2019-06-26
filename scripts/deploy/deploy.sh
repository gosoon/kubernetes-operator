#!/bin/bash

export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8

ROOT=`cd $(dirname $0); pwd`
DEPLOY_MASTER_LOG="/var/log/deploy_master.log"
DEPLOY_NODE_LOG="/var/log/deploy_node.log"
DEPLOY_ETCD_LOG="/var/log/deploy_etcd.log"

master() {
    echo -n "starting deploy master... "
    [ -f ${DEPLOY_MASTER_LOG} ] && { echo "master is deployed"; exit; } || touch ${DEPLOY_MASTER_LOG}
    [ -f base_env_01.sh ] && bash -x base_env_01.sh | tee -a ${DEPLOY_MASTER_LOG} || exit 1
    [ -f deploy_master.sh ] && bash -x deploy_master.sh | tee -a ${DEPLOY_MASTER_LOG} || exit 1
}

node() {
    echo -n "starting deploy node... "
    [ -f ${DEPLOY_NODE_LOG} ] && { echo "node is deployed"; exit; } || touch ${DEPLOY_NODE_LOG}
    [ -f base_env_01.sh ] && bash -x base_env_01.sh | tee -a ${DEPLOY_NODE_LOG} || exit 1
    [ -f deploy_kubelet.sh ] && bash -x deploy_kubelet.sh | tee -a ${DEPLOY_NODE_LOG} || exit 1
    [ -f deploy_kube_proxy.sh ] && bash -x deploy_kube_proxy.sh | tee -a ${DEPLOY_NODE_LOG} || exit 1
    [ -f deploy_coredns.sh ] && bash -x deploy_coredns.sh | tee -a ${DEPLOY_NODE_LOG} || exit 1
    [ -f deploy_calico.sh ] && bash -x deploy_calico.sh | tee -a ${DEPLOY_NODE_LOG} || exit 1
}

etcd() {
    echo -n "starting deploy etcd... "
    [ -f ${DEPLOY_ETCD_LOG} ] && { echo "etcd is deployed"; exit; } || touch ${DEPLOY_ETCD_LOG}
    [ -f base_env_01.sh ] && bash -x base_env_01.sh | tee -a ${DEPLOY_ETCD_LOG} || exit 1
    [ -f deploy_etcd.sh ] && bash -x deploy_etcd.sh | tee -a ${DEPLOY_ETCD_LOG} || exit 1
}


case "$1" in
master)
    master
    ;;

node)
    node
    ;;

etcd)
    etcd
    ;;

*)
    echo "Usage: $0 {master|node|etcd}"
    exit 1
esac
