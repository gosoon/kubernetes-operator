#!/bin/bash

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./version.sh ] && . ./version.sh || exit

# k8s version >= v1.9
# docs: https://github.com/coredns/deployment/tree/master/kubernetes

# TODO : download coredns 
download_coredns(){
    git clone  https://github.com/coredns/deployment
}

COERDNS_CONFIG="../yaml/coredns_${COREDNS_VER}/coredns.yaml"

kubectl apply -f ${COERDNS_CONFIG}
