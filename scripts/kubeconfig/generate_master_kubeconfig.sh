#!/bin/bash
# 指定 apiserver 地址
KUBE_APISERVER="https://<apiserver_ip>:6443"

# 生成 Bootstrap Token
BOOTSTRAP_TOKEN_ID=$(head -c 6 /dev/urandom | md5sum | head -c 6)
BOOTSTRAP_TOKEN_SECRET=$(head -c 16 /dev/urandom | md5sum | head -c 16)
BOOTSTRAP_TOKEN="${BOOTSTRAP_TOKEN_ID}.${BOOTSTRAP_TOKEN_SECRET}"

[ -d output ] || mkdir output

# 生成 kubectl 配置文件
echo "Create kubectl kubeconfig..."
kubectl config set-cluster kubernetes \
  --certificate-authority=k8s-root-ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER} \
  --kubeconfig=output/kubectl.kubeconfig
kubectl config set-credentials "system:masters" \
  --client-certificate=admin.pem \
  --client-key=admin.pem \
  --embed-certs=true \
  --kubeconfig=output/kubectl.kubeconfig
kubectl config set-context default \
  --cluster=kubernetes \
  --user=system:masters \
  --kubeconfig=output/kubectl.kubeconfig
kubectl config use-context default --kubeconfig=output/kubectl.kubeconfig

# 生成 kubelet tls bootstrap 配置
echo "Create kubelet bootstrapping kubeconfig..."
kubectl config set-cluster kubernetes \
  --certificate-authority=k8s-root-ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER} \
  --kubeconfig=output/bootstrap.kubeconfig
kubectl config set-credentials "system:bootstrap:${BOOTSTRAP_TOKEN_ID}" \
  --token=${BOOTSTRAP_TOKEN} \
  --kubeconfig=output/bootstrap.kubeconfig
kubectl config set-context default \
  --cluster=kubernetes \
  --user="system:bootstrap:${BOOTSTRAP_TOKEN_ID}" \
  --kubeconfig=output/bootstrap.kubeconfig
kubectl config use-context default --kubeconfig=output/bootstrap.kubeconfig

# 生成 kube-controller-manager 配置文件
echo "Create kube-controller-manager kubeconfig..."
kubectl config set-cluster kubernetes \
  --certificate-authority=k8s-root-ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER} \
  --kubeconfig=output/kube-controller-manager.kubeconfig
kubectl config set-credentials "system:kube-controller-manager" \
  --client-certificate=kube-controller-manager.pem \
  --client-key=kube-controller-manager-key.pem \
  --embed-certs=true \
  --kubeconfig=output/kube-controller-manager.kubeconfig
kubectl config set-context default \
  --cluster=kubernetes \
  --user=system:kube-controller-manager \
  --kubeconfig=output/kube-controller-manager.kubeconfig
kubectl config use-context default --kubeconfig=output/kube-controller-manager.kubeconfig

# 生成 kube-scheduler 配置文件
echo "Create kube-scheduler kubeconfig..."
kubectl config set-cluster kubernetes \
  --certificate-authority=k8s-root-ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER} \
  --kubeconfig=output/kube-scheduler.kubeconfig
kubectl config set-credentials "system:kube-scheduler" \
  --client-certificate=kube-scheduler.pem \
  --client-key=kube-scheduler-key.pem \
  --embed-certs=true \
  --kubeconfig=output/kube-scheduler.kubeconfig
kubectl config set-context default \
  --cluster=kubernetes \
  --user=system:kube-scheduler \
  --kubeconfig=output/kube-scheduler.kubeconfig
kubectl config use-context default --kubeconfig=output/kube-scheduler.kubeconfig
