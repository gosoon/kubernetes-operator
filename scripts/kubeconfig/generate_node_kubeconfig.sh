#!/bin/bash
# 指定 apiserver 地址
KUBE_APISERVER="https://<apiserver_ip>:6443"

# 生成 Bootstrap Token
BOOTSTRAP_TOKEN_ID=$(head -c 6 /dev/urandom | md5sum | head -c 6)
BOOTSTRAP_TOKEN_SECRET=$(head -c 16 /dev/urandom | md5sum | head -c 16)
BOOTSTRAP_TOKEN="${BOOTSTRAP_TOKEN_ID}.${BOOTSTRAP_TOKEN_SECRET}"

[ -d output ] || mkdir output

# 生成 kubelet 配置文件
echo "Create kubelet kubeconfig..."
kubectl config set-cluster kubernetes \
  --certificate-authority=k8s-root-ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER} \
  --kubeconfig=output/kubelet.kubeconfig
kubectl config set-credentials "system:masters" \
  --client-certificate=kubelet-client.pem \
  --client-key=kubelet-client-key.pem \
  --embed-certs=true \
  --kubeconfig=output/kubelet.kubeconfig
kubectl config set-context default \
  --cluster=kubernetes \
  --user=system:masters \
  --kubeconfig=output/kubelet.kubeconfig
kubectl config use-context default --kubeconfig=output/kubelet.kubeconfig

# 生成 kube-proxy 配置文件
echo "Create kube-proxy kubeconfig..."
kubectl config set-cluster kubernetes \
  --certificate-authority=k8s-root-ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER} \
  --kubeconfig=output/kube-proxy.kubeconfig
kubectl config set-credentials "system:kube-proxy" \
  --client-certificate=kube-proxy.pem \
  --client-key=kube-proxy-key.pem \
  --embed-certs=true \
  --kubeconfig=output/kube-proxy.kubeconfig
kubectl config set-context default \
  --cluster=kubernetes \
  --user=system:kube-proxy \
  --kubeconfig=output/kube-proxy.kubeconfig
kubectl config use-context default --kubeconfig=output/kube-proxy.kubeconfig
