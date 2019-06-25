#!/bin/bash

# author: tianfeiyu

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./config.sh ] && . ./config.sh || exit

CALICO_CONFIG="../yaml/calico_${CALICO_VER}/calico-etcd.yaml"

download_calico(){
    wget -O ${CALICO_CONFIG}  https://docs.projectcalico.org/${CALICO_VER}/manifests/calico-etcd.yaml
}

# TODO : check download_calico
kubectl apply -f ${CALICO_CONFIG}
