#!/bin/bash

# author:tianfeiyu

set -x

#TODO : check os version and use root exec

swapoff -a
sed -i 's/SELINUX=permissive/SELINUX=disabled/' /etc/sysconfig/selinux 
setenforce 0

systemctl disable firewalld.service && systemctl stop firewalld.service

cat << EOF >> /etc/sysctl.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
vm.swappiness=0
EOF 

sysctl -p

cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
