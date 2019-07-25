#!/bin/bash

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./config.sh ] && . ./config.sh || exit

KUBE_NODE_BIN_DIR="../bin/${KUBERNETES_VER}"
KUBE_NODE_CONFIG_DIR="../config/node"
KUBE_NODE_SYSTEMD_CONFIG_DIR="../systemd"
CERTS_DIR="../certs"
KUBECONFIG_DIR="../kubeconfig/"
GENERATE_CERTS_FILE="../certs/node"
GENERATE_KUBECONFIG_FILE="../kubeconfig"

DEST_CERTS_DIR="/etc/kubernetes/ssl"
DEST_SYSTEMD_DIR="/usr/lib/systemd/system"
DEST_CONFIG_DIR="/etc/kubernetes"
KUBE_MASTER_LOG="/var/log/kubernetes"

cp ${KUBE_NODE_BIN_DIR}/kubelet /usr/bin/

cp ${KUBE_NODE_SYSTEMD_CONFIG_DIR}/kubelet.service  ${DEST_SYSTEMD_DIR}

[ -d ${DEST_CONFIG_DIR} ] || mkdir ${DEST_CONFIG_DIR}
cp ${KUBE_NODE_CONFIG_DIR}/{config,config.yaml,kubelet} ${DEST_CONFIG_DIR}/

# update config master ip
sed -i -e "s#--master=https://<apiserver_ip>:6443#--master=https://${MASTER_VIP}:6443#g" ${DEST_CONFIG_DIR}/config
sed -i -e "s#--hostname_override=<node_ip>#--hostname_override=${LOCAL_IP}#g" ${DEST_CONFIG_DIR}/kubelet
sed -i -e "s#hostnameOverride: <node_ip>#hostnameOverride: ${LOCAL_IP}#g" ${DEST_CONFIG_DIR}/config.yaml

# scp ssl from master
for master in `echo ${MASTER_HOSTS} | tr ',' ' '`;do
    scp root@${master}:${DEST_CERTS_DIR}/{ca.pem,ca-key.pem} ${DEST_CERTS_DIR}/
    [ $? -eq 0 ] && break
done

# generate ssl
#cd ${GENERATE_CERTS_FILE} && bash gen_cert.sh
#[ $? -eq 0 ] && echo "generate certs success" || exit 1
#cd -
#[ -d ${DEST_CERTS_DIR} ] || mkdir ${DEST_CERTS_DIR}
#cp ${GENERATE_CERTS_FILE}/output/{ca.pem,ca-key.pem,kube-proxy-key.pem,kube-proxy.pem,kubelet-client-key.pem, \
#kubelet-client.pem} ${DEST_CERTS_DIR}/

# scp kubeconfig from master
for master in `echo ${MASTER_HOSTS} | tr ',' ' '`;do
    scp root@${master}:/home/kubernetes-operator/scripts/kubeconfig/output/kubelet-${LOCAL_IP}.kubeconfig \
    ${DEST_CONFIG_DIR}/kubelet.kubeconfig
    [ $? -eq 0 ] && break
done

# generate kubeconfig
#cp ${GENERATE_KUBECONFIG_FILE}/output/{kubelet.kubeconfig,bootstrap.kubeconfig} ${DEST_CONFIG_DIR}/
#sed -i -e "s#https://<apiserver_ip>:6443#https://${MASTER_VIP}:6443#g" ${GENERATE_KUBECONFIG_FILE}/generate_node_kubeconfig.sh
#cd ${GENERATE_KUBECONFIG_FILE} && bash generate_node_kubeconfig.sh
#[ $? -eq 0 ] && echo "generate kubeconfig success" || exit 1
#cp ${GENERATE_KUBECONFIG_FILE}/output/* ${DEST_CONFIG_DIR}/
#cd -

# mkdir log dir
[ -d ${KUBE_MASTER_LOG} ] || mkdir -pv ${KUBE_MASTER_LOG}

# start service
systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet
systemctl status kubelet

if [ $? -ne 0 ];then  
    echo "deploy kubelet failed !!!" && exit 1
fi
