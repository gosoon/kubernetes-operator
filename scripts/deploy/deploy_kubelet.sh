#!/bin/bash

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./config.sh ] && . ./config.sh || exit

KUBE_NODE_BIN_DIR="../bin/${KUBERNETES_VER}"
KUBE_NODE_CONFIG_DIR="../config/node"
KUBE_NODE_SYSTEMD_CONFIG_DIR="../systemd"
CERTS_DIR="../certs"
KUBECONFIG_DIR="../kubeconfig/"

DEST_CERTS_DIR="/etc/kubernetes/ssl"
DEST_SYSTEMD_DIR="/usr/lib/systemd/system"
DEST_CONFIG_DIR="/etc/kubernetes"
KUBE_MASTER_LOG="/var/log/kubernetes"

cp ${KUBE_NODE_BIN_DIR}/kubelet /usr/bin/

cp ${KUBE_NODE_SYSTEMD_CONFIG_DIR}/kubelet.service  ${DEST_SYSTEMD_DIR}

cp ${KUBE_NODE_CONFIG_DIR}/{config,config.yaml,kubelet} /etc/kubernetes/

# dest dir use ansible confirm
# cp ssl 
cp ${CERTS_DIR}/{apiserver-client-key.pem,apiserver-client.csr,apiserver-client.pem} ${DEST_CERTS_DIR}

# cp kubeconfig 
cp ${KUBECONFIG_DIR}/kubelet.kubeconfig  ${DEST_CONFIG_DIR}/

sed -i -e "s#--master=https://<apiserver_ip>:6443#--master=https://${MASTER_HOSTS}:6443#g" ${DEST_CONFIG_DIR}/config
#TODO: update config and kubeconfig master ip
sed -i -e "s#--hostname_override=<node_ip>#--hostname_override=${LOCAL_IP}#g" ${DEST_CONFIG_DIR}/kubelet
sed -i -e "s#hostnameOverride: <node_ip>#hostnameOverride: ${LOCAL_IP}#g" ${DEST_CONFIG_DIR}/config.yaml

[ -d ${KUBE_MASTER_LOG} ] || mkdir -pv ${KUBE_MASTER_LOG}

systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet
systemctl status kubelet
