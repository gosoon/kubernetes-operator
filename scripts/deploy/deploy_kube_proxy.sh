#!/bin/bash

# author: tianfeiyu

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./config.sh ] && . ./config.sh || exit


# TODO : if etcd not download and download etcd

KUBE_PROXY_BIN_DIR="../bin/${KUBERNETES_VER}"
KUBE_PROXY_SYSTEMD_CONFIG_DIR="../systemd"
KUBE_PROXY_CONFIG_DIR="../config/node"
KUBECONFIG_DIR="../kubeconfig/"

DEST_CONFIG_DIR="/etc/kubernetes"
DEST_SYSTEMD_DIR="/usr/lib/systemd/system/"
KUBE_MASTER_LOG="/var/log/kubernetes"

pre_ipvs() {

    cat >> /etc/sysctl.conf <<EOF
net.ipv4.ip_forward=1
net.bridge.bridge-nf-call-iptables=1
net.bridge.bridge-nf-call-ip6tables=1
EOF

    sysctl -p

    cat >> /etc/modules <<EOF
ip_vs
ip_vs_lc
ip_vs_wlc
ip_vs_rr
ip_vs_wrr
ip_vs_lblc
ip_vs_lblcr
ip_vs_dh
ip_vs_sh
ip_vs_fo
ip_vs_nq
ip_vs_sed
ip_vs_ftp
EOF

    yum install -y conntrack ipvsadm
}

pre_ipvs
cp ${KUBE_PROXY_BIN_DIR}/kube-proxy /usr/bin/

cp ${KUBE_PROXY_SYSTEMD_CONFIG_DIR}/kube-proxy.service  ${DEST_SYSTEMD_DIR}

cp ${KUBE_PROXY_CONFIG_DIR}/kube-proxy /etc/kubernetes/

# cp kubeconfig 
cp ${KUBECONFIG_DIR}/output/kube-proxy.kubeconfig  ${DEST_CONFIG_DIR}/

#TODO: update config and kubeconfig master ip
sed -i -e "s#--hostname_override=<node_ip>#--hostname_override=${LOCAL_IP}#g" ${DEST_CONFIG_DIR}/kube-proxy

[ -d ${KUBE_MASTER_LOG} ] || mkdir -pv ${KUBE_MASTER_LOG}

systemctl daemon-reload
systemctl enable kube-proxy
systemctl start kube-proxy
systemctl status kube-proxy
