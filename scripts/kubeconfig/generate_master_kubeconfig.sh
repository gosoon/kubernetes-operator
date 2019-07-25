#!/bin/bash

[ -e ../deploy/config.sh ] && . ../deploy/config.sh || exit 1

KUBE_APISERVER="https://${MASTER_VIP}:6443"
CERTS_DIR="${DEPLOY_HOME_DIR}/kubernetes-operator/scripts/certs/master/output"

[ -d output ] || mkdir output

# 生成 kubectl 配置文件
echo "Create kubectl kubeconfig..."
kubectl config set-cluster kubernetes \
  --certificate-authority=${CERTS_DIR}/ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER} \
  --kubeconfig=output/kubectl.kubeconfig
kubectl config set-credentials "system:masters" \
  --client-certificate=${CERTS_DIR}/admin.pem \
  --client-key=${CERTS_DIR}/admin-key.pem \
  --embed-certs=true \
  --kubeconfig=output/kubectl.kubeconfig
kubectl config set-context default \
  --cluster=kubernetes \
  --user=system:masters \
  --kubeconfig=output/kubectl.kubeconfig
kubectl config use-context default --kubeconfig=output/kubectl.kubeconfig

# 生成 kube-controller-manager 配置文件
echo "Create kube-controller-manager kubeconfig..."
kubectl config set-cluster kubernetes \
  --certificate-authority=${CERTS_DIR}/ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER} \
  --kubeconfig=output/kube-controller-manager.kubeconfig
kubectl config set-credentials "system:kube-controller-manager" \
  --client-certificate=${CERTS_DIR}/kube-controller-manager.pem \
  --client-key=${CERTS_DIR}/kube-controller-manager-key.pem \
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
  --certificate-authority=${CERTS_DIR}/ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER} \
  --kubeconfig=output/kube-scheduler.kubeconfig
kubectl config set-credentials "system:kube-scheduler" \
  --client-certificate=${CERTS_DIR}/kube-scheduler.pem \
  --client-key=${CERTS_DIR}/kube-scheduler-key.pem \
  --embed-certs=true \
  --kubeconfig=output/kube-scheduler.kubeconfig
kubectl config set-context default \
  --cluster=kubernetes \
  --user=system:kube-scheduler \
  --kubeconfig=output/kube-scheduler.kubeconfig
kubectl config use-context default --kubeconfig=output/kube-scheduler.kubeconfig

#for node in `echo ${NODE_HOSTS} | tr ',' ' '`;do
    ## 生成 kubelet 配置文件
    #echo "Create kubelet kubeconfig..."
    #kubectl config set-cluster kubernetes \
      #--certificate-authority=${CERTS_DIR}/ca.pem \
      #--embed-certs=true \
      #--server=${KUBE_APISERVER} \
      #--kubeconfig=output/kubelet-${node}.kubeconfig

    #kubectl config set-credentials system:node:${node} \
      #--client-certificate=${CERTS_DIR}/kubelet.pem \
      #--client-key=${CERTS_DIR}/kubelet-key.pem \
      #--embed-certs=true \
      #--kubeconfig=output/kubelet-${node}.kubeconfig

    #kubectl config set-context default \
      #--cluster=kubernetes \
      #--user=system:node:${node} \
      #--kubeconfig=output/kubelet-${node}.kubeconfig
    #kubectl config use-context default --kubeconfig=output/kubelet-${node}.kubeconfig
#done

# 生成 kube-proxy 配置文件
#echo "Create kube-proxy kubeconfig..."
#kubectl config set-cluster kubernetes \
  #--certificate-authority=${CERTS_DIR}/ca.pem \
  #--embed-certs=true \
  #--server=${KUBE_APISERVER} \
  #--kubeconfig=output/kube-proxy.kubeconfig
#kubectl config set-credentials "system:kube-proxy" \
  #--client-certificate=${CERTS_DIR}/kube-proxy.pem \
  #--client-key=${CERTS_DIR}/kube-proxy-key.pem \
  #--embed-certs=true \
  #--kubeconfig=output/kube-proxy.kubeconfig
#kubectl config set-context default \
  #--cluster=kubernetes \
  #--user=system:kube-proxy \
  #--kubeconfig=output/kube-proxy.kubeconfig
#kubectl config use-context default --kubeconfig=output/kube-proxy.kubeconfig
