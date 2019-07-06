#!/bin/bash

set -x

#TODO : check os version and use root exec

swapoff -a
sed -i 's/SELINUX=permissive/SELINUX=disabled/' /etc/sysconfig/selinux 
setenforce 0

systemctl disable firewalld.service && systemctl stop firewalld.service

# use aliyun kubernetes yum source
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF

# add kube use 
[ -d "/home/kube" ] || useradd kube
