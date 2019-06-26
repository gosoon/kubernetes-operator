#!/bin/bash


[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ../deploy/config.sh ] && . ../deploy/config.sh || exit

systemctl stop kubelet kube-proxy 
rm -rf /var/lib/kubelet/*
rm -f /var/log/deploy_node.log
