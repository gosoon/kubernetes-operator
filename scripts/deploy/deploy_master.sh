#!/bin/bash

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./config.sh ] && . ./config.sh || exit


KUBE_MASTER_BIN_DIR="../bin/${KUBERNETES_VER}"
KUBE_MASTER_CONFIG_DIR="../config/master"
KUBE_MASTER_SYSTEMD_CONFIG_DIR="../systemd"
CERTS_DIR="../certs"

DEST_CONFIG_DIR="/etc/kubernetes"
DEST_SYSTEMD_DIR="/usr/lib/systemd/system"
DEST_CERTS_DIR="/etc/kubernetes/ssl"

cp ${KUBE_MASTER_BIN_DIR}/{kube-apiserver,kube-controller-manager,kube-scheduler,kubectl} /usr/bin/

cp ${KUBE_MASTER_SYSTEMD_CONFIG_DIR}/{kube-apiserver.service,kube-controller-manager.service,kube-scheduler.service} ${DEST_SYSTEMD_DIR}/

# cp config, apiserver config controller-manager scheduler 
cp ${KUBE_MASTER_CONFIG_DIR}/* ${DEST_CONFIG_DIR}/
# https://127.0.0.1:2379
etcd_num=$(echo ${ETCD_HOSTS} | awk -F ',' '{print NF}')
etcd_cluster=""
for i in `seq 1 ${etcd_num}`;do
	ip=$(echo ${ETCD_HOSTS} | awk -v idx=$i -F ',' '{print $idx}')
    cluster=$(echo "https://${ip}:2379")
    if [ $i -ne ${etcd_num} ];then
        cluster="${cluster},"
    fi
    etcd_cluster="${etcd_cluster}${cluster}"
done

sed -i -e "s#--etcd-servers=xxx#--etcd-servers=${etcd_cluster}#g" ${DEST_CONFIG_DIR}/apiserver
sed -i -e "s#--master=https#--master=https://${LOCAL_IP}:6443#g" ${DEST_CONFIG_DIR}/config

# cp ssl
cp ${CERTS_DIR}/{apiserver-client-key.pem,apiserver-client.csr,apiserver-client.pem,apiserver-server-key.pem,apiserver-server.csr,apiserver-server.pem,ca.csr,ca.pem,ca-key.pem}  ${DEST_CERTS_DIR}/

#TODO : copy kubeconfig

systemctl daemon-reload
systemctl enable kube-apiserver kube-controller-manager kube-scheduler 
systemctl start kube-apiserver kube-controller-manager kube-scheduler
systemctl status kube-apiserver kube-controller-manager kube-scheduler
