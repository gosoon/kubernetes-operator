#!/bin/bash

systemctl stop kube-apiserver kube-controller-manager kube-scheduler
rm -f /usr/bin/{kube-apiserver,kube-controller-manager,kube-schedule,kubectl}
rm -rf /etc/kubernetes/
rm -rf /root/.kube/
rm -f /var/log/deploy_master.log
rm -rf /var/log/kubernetes/
