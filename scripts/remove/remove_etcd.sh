#!/bin/bash

[ -e ../deploy/config.sh ] && . ../deploy/config.sh || exit

CONFIG_DIR="/etc/etcd"
source ${CONFIG_DIR}/etcd.conf

systemctl stop etcd
userdel etcd
rm -f /usr/bin/etcd
rm -rf ${ETCD_DATA_DIR}
rm -rf ${CONFIG_DIR} 
rm -f /var/log/deploy_etcd.log
