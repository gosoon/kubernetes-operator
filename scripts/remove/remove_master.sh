#!/bin/bash


[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ../deploy/config.sh ] && . ../deploy/config.sh || exit

systemctl stop kube-apiserver kube-controller-manager kube-scheduler
rm -rf /etc/kubernetes/
rm -f /var/log/deploy_master.log
rm -rf /var/log/kubernetes/
