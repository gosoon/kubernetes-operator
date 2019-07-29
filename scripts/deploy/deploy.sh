#!/bin/bash

export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8

ROOT=`cd $(dirname $0); pwd`
DEPLOY_MASTER_LOG="/var/log/deploy_master.log"
DEPLOY_NODE_LOG="/var/log/deploy_node.log"
DEPLOY_ETCD_LOG="/var/log/deploy_etcd.log"
DEPLOY_MASTER_FINISHED="/var/log/master_deployed"
DEPLOY_NODE_FINISHED="/var/log/node_deployed"
DEPLOY_ETCD_FINISHED="/var/log/etcd_deployed"

master() {
    echo -n "starting deploy master... "
    [ -f ${DEPLOY_MASTER_FINISHED} ] && { echo "master is deployed"; exit 1; } || > ${DEPLOY_MASTER_LOG}
    [ -f base_env_01.sh ] && bash -x base_env_01.sh &>> ${DEPLOY_MASTER_LOG} || exit 1
    [ -f deploy_master.sh ] && bash -x deploy_master.sh &>> ${DEPLOY_MASTER_LOG} || exit 1
    [ -f deploy_coredns.sh ] && bash -x deploy_coredns.sh &>> ${DEPLOY_MASTER_LOG} || exit 1
    [ -f deploy_calico.sh ] && bash -x deploy_calico.sh &>> ${DEPLOY_MASTER_LOG} || exit 1
    touch ${DEPLOY_MASTER_FINISHED}
}

node() {
    echo -n "starting deploy node... "
    [ -f ${DEPLOY_NODE_FINISHED} ] && { echo "node is deployed"; exit 1; } || > ${DEPLOY_NODE_LOG}
    [ -f base_env_01.sh ] && bash -x base_env_01.sh &>> ${DEPLOY_NODE_LOG} || exit 1
    [ -f deploy_docker.sh ] && bash -x deploy_docker.sh &>> ${DEPLOY_NODE_LOG} || exit 1
    [ -f deploy_kubelet.sh ] && bash -x deploy_kubelet.sh &>> ${DEPLOY_NODE_LOG} || exit 1
    [ -f deploy_kube_proxy.sh ] && bash -x deploy_kube_proxy.sh &>> ${DEPLOY_NODE_LOG} || exit 1
    touch ${DEPLOY_NODE_FINISHED}
}

etcd() {
    echo -n "starting deploy etcd... "
    [ -f ${DEPLOY_ETCD_FINISHED} ] && { echo "etcd is deployed"; exit 1; } || > ${DEPLOY_ETCD_LOG}
    [ -f base_env_01.sh ] && bash -x base_env_01.sh &>> ${DEPLOY_ETCD_LOG} || exit 1
    [ -f deploy_etcd.sh ] && bash -x deploy_etcd.sh &>> ${DEPLOY_ETCD_LOG} || exit 1
    touch ${DEPLOY_ETCD_FINISHED}
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
