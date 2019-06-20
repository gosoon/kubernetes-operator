#!/bin/bash

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./version.sh ] && . ./version.sh || exit

# etcd download releases doc: https://github.com/etcd-io/etcd/releases/

ETCD_BIN_DIR="../bin/${ETCD_VER}"
CERTS_DIR="../certs/"
ETCD_SYSTEMD_CONFIG_DIR="../systemd"
ETCD_CONFIG_DIR="../config/etcd"


dest_etcd_config_dir="/etc/etcd/etcd.conf"
DEST_SYSTEMD_DIR="/usr/lib/systemd/system"

# TODO : if etcd not download and download etcd
download_etcd(){
    # choose either URL
    GOOGLE_URL=https://storage.googleapis.com/etcd
    GITHUB_URL=https://github.com/etcd-io/etcd/releases/download
    DOWNLOAD_URL=${GOOGLE_URL}

    rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
    rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test

    curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
    tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1
    rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz

    /tmp/etcd-download-test/etcd --version
    ETCDCTL_API=3 /tmp/etcd-download-test/etcdctl version
}

cp ${ETCD_BIN_DIR}/{etcd,etcdctl} /usr/bin/

cp ${ETCD_SYSTEMD_CONFIG_DIR}/{etcd.service}  ${DEST_SYSTEMD_DIR}

cp ${ETCD_CONFIG_DIR}/* /etc/etcd/

# cp ssl
cp ${CERTS_DIR}{etcd-client-key.pem,etcd-peer-key.pem,etcd-peer.pem,etcd-server-key.pem,etcd-server.pem,ca.pem}  /etc/kubernetes/ssl/

#TODO : replace etcd ip addr
systemctl daemon-reload
systemctl enable etcd
systemctl start etcd
systemctl status etcd
