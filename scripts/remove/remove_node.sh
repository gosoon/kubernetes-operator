#!/bin/bash

systemctl stop docker kubelet kube-proxy 
rm -f /usr/bin/{kubelet,kube-proxy}
rm -rf /var/lib/kubelet/
rm -rf /usr/lib/systemd/system/kubelet.service.d/
rm -rf /etc/kubernetes/
rm -f /var/log/deploy_node.log
