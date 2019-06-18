#!/bin/bash

# author: tianfeiyu

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./version.sh ] && . ./version.sh || exit


# TODO : if etcd not download and download etcd




systemctl daemon-reload
systemctl enable kube-proxy
systemctl start kube-proxy
systemctl status kube-proxy
