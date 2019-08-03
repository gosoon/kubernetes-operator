#!/bin/bash

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./config.sh ] && . ./config.sh || exit

# etcd download releases doc: https://github.com/etcd-io/etcd/releases/

ETCD_BIN_DIR="../bin/${ETCD_VER}"
CERTS_DIR="../certs/"
ETCD_SYSTEMD_CONFIG_DIR="../systemd"
ETCD_CONFIG_DIR="../config/etcd"
GENERATE_CERTS_FILE="../certs/etcd"

DEST_ETCD_CONFIG_DIR="/etc/etcd/etcd.conf"
DEST_SYSTEMD_DIR="/usr/lib/systemd/system"
DEST_CONFIG_DIR="/etc/etcd"
DEST_CERTS_DIR="/etc/etcd/ssl"

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

# cp etcd config
cp ${ETCD_BIN_DIR}/{etcd,etcdctl} /usr/bin/ && chmod +x /usr/bin/{etcd,etcdctl}
cp ${ETCD_SYSTEMD_CONFIG_DIR}/etcd.service  ${DEST_SYSTEMD_DIR}/
[ -d ${DEST_CONFIG_DIR} ] || mkdir ${DEST_CONFIG_DIR}
cp ${ETCD_CONFIG_DIR}/* ${DEST_CONFIG_DIR}/

# generate ssl 
cd ${GENERATE_CERTS_FILE} && bash gen_cert.sh
[ $? -eq 0 ] && echo "generate certs success" || exit 1
cd -

# create ssl dir
[ -d ${DEST_CERTS_DIR} ] || mkdir ${DEST_CERTS_DIR}
cp ${GENERATE_CERTS_FILE}/output/* ${DEST_CERTS_DIR}/

# config etcd 
source ${ETCD_CONFIG_DIR}/*
etcd_num=$(echo ${ETCD_HOSTS} | awk -F ',' '{print NF}')
etcd_cluster=""
for i in `seq 1 ${etcd_num}`;do
	ip=$(echo ${ETCD_HOSTS} | awk -v idx=$i -F ',' '{print $idx}')
    cluster=$(echo "node$i=https://${ip}:2380")
    if [ $i -ne ${etcd_num} ];then
        cluster="${cluster},"
    fi
    etcd_cluster="${etcd_cluster}${cluster}"
    if [ "X${ip}" == "X${LOCAL_IP}" ];then
        sed -i -e "s@\(ETCD_NAME=\).*@\1\"node${i}\"@g" ${DEST_CONFIG_DIR}/etcd.conf
        break
    fi
done

# replace ip addr
sed -i -e "s@\(ETCD_INITIAL_ADVERTISE_PEER_URLS=\).*@\1\"https://${LOCAL_IP}:2380\"@g" ${DEST_CONFIG_DIR}/etcd.conf
sed -i -e "s@\(ETCD_ADVERTISE_CLIENT_URLS=\).*@\1\"https://${LOCAL_IP}:2379\"@g" ${DEST_CONFIG_DIR}/etcd.conf
sed -i -e "s@\(ETCD_INITIAL_CLUSTER=\).*@\1\"${etcd_cluster}\"@g" ${DEST_CONFIG_DIR}/etcd.conf

useradd etcd 
[ -d ${ETCD_DATA_DIR} ] || mkdir -pv ${ETCD_DATA_DIR}
chown -R etcd:etcd ${ETCD_DATA_DIR}

# start etcd service
systemctl daemon-reload
systemctl enable etcd
systemctl start etcd
systemctl status etcd

if [ $? -ne 0 ];then  
    echo "deploy etcd failed !!!" && exit 1
fi
