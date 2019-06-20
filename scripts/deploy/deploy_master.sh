#!/bin/bash

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./version.sh ] && . ./version.sh || exit


kube_master_bin_dir="../bin/${KUBERNETES_VER}"
kube_master_config_dir="../config/master"
kube_master_systemd_config_dir="../systemd"
cert_dir="../certs"

DEST_SYSTEMD_DIR="/usr/lib/systemd/system/"

cp ${kube_master_bin_dir}/{kube-apiserver,kube-controller-manager,kube-scheduler,kubectl} /usr/bin/

cp ${kube_master_systemd_config_dir}/{kube-apiserver.service,kube-controller-manager.service,kube-scheduler.service} ${DEST_SYSTEMD_DIR}

cp ${kube_master_config_dir}/* /etc/kubernetes/

# cp ssl
cp ${cert_dir}/{apiserver-client-key.pem,apiserver-client.csr,apiserver-client.pem,apiserver-server-key.pem,apiserver-server.csr,apiserver-server.pem,ca.csr,ca.pem,ca-key.pem}  /etc/kubernetes/ssl/


#TODO : copy kubeconfig

systemctl daemon-reload
systemctl enable kube-apiserver kube-controller-manager kube-scheduler 
systemctl start kube-apiserver kube-controller-manager kube-scheduler
systemctl status kube-apiserver kube-controller-manager kube-scheduler
