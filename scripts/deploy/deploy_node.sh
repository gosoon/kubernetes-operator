#!/bin/bash

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./version.sh ] && . ./version.sh || exit

kube_node_bin_dir="../bin/${KUBERNETES_VER}"
kube_node_config_dir="../config/node"
kube_node_systemd_config_dir="../systemd"
cert_dir="../certs"
kubeconfig_dir="../kubeconfig/"

dest_cert_dir="/etc/kubernetes/ssl/"
dest_kubeconfig_dir="/etc/kubernetes/kubeconfig/"
DEST_SYSTEMD_DIR="/usr/lib/systemd/system/"


cp ${kube_node_bin_dir}/{kubelet,kube-proxy} /usr/bin/

cp ${kube_node_systemd_config_dir}/{kubelet.service,kube-proxy.service}  ${DEST_SYSTEMD_DIR}

cp ${kube_node_config_dir}/* /etc/kubernetes/

# dest dir use ansible confirm
# cp ssl 
cp ${cert_dir}/{apiserver-client-key.pem,apiserver-client.csr,apiserver-client.pem} ${dest_cert_dir}

# cp kubeconfig 
cp ${kubeconfig_dir}/kubelet.kubeconfig  ${dest_kubeconfig_dir}


#TODO: update config and kubeconfig master ip
