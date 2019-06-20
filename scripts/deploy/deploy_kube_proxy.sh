#!/bin/bash

# author: tianfeiyu

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./version.sh ] && . ./version.sh || exit


# TODO : if etcd not download and download etcd

KUBE_PROXY_BIN_DIR="../bin/${KUBERNETES_VER}"
KUBE_PROXY_SYSTEMD_CONFIG_DIR="../systemd"
KUBE_PROXY_CONFIG_DIR="../config/node"
KUBECONFIG_DIR="../kubeconfig/"

DEST_KUBECONFIG_DIR="/etc/kubernetes/kubeconfig/"
DEST_SYSTEMD_DIR="/usr/lib/systemd/system/"


cp ${KUBE_PROXY_BIN_DIR}/kube-proxy /usr/bin/

cp ${KUBE_PROXY_SYSTEMD_CONFIG_DIR}/kube-proxy.service  ${DEST_SYSTEMD_DIR}

cp ${KUBE_PROXY_CONFIG_DIR}/kube-proxy /etc/kubernetes/

# cp kubeconfig 
cp ${KUBECONFIG_DIR}/kube-proxy.kubeconfig  ${DEST_KUBECONFIG_DIR}

#TODO: update config and kubeconfig master ip
systemctl daemon-reload
systemctl enable kube-proxy
systemctl start kube-proxy
systemctl status kube-proxy
