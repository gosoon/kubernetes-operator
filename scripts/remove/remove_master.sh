#!/bin/bash


[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ../deploy/config.sh ] && . ../deploy/config.sh || exit

source ${ETCD_CONFIG_DIR}/*

CONFIG_DIR="/etc/etcd"

userdel etcd
rm -rf ${ETCD_DATA_DIR}
rm -rf ${CONFIG_DIR} 
rm -f /var/log/deploy_master.log
