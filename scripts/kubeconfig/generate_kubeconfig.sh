#!/bin/bash

OLDPATH=`pwd`
ROOT=`cd $(dirname $0); pwd`
cd ${ROOT}
function finish {
    cd ${OLDPATH}
}
trap finish EXIT
CERTS=${ROOT}/output
SCRIPT_NAME=$(basename $0)
OUTPUT=${ROOT}/kubeconfig

mkdir -p ${OUTPUT}


KUBE_APISERVER=""
if [[ "$#" -lt 2 ]]; then
    echo "need the apiserver secure endpoint, exec the script like './${SCRIPT_NAME} hxy01 https://10.0.2.15:6443'."
    exit 1
fi
REGION="$1"
KUBE_APISERVER="$2"
echo "use ${KUBE_APISERVER} as the apiserver endpoint."

kubectl config set-cluster kubernetes \
    --certificate-authority=root/ca.pem \
    --embed-certs=true \
    --server=${KUBE_APISERVER} \
    --kubeconfig=${OUTPUT}/client.kubeconfig

 kubectl config set-credentials client \
    --client-certificate=${CERTS}/apiserver-client.pem \
    --client-key=${CERTS}/apiserver-client-key.pem \
    --embed-certs=true \
    --kubeconfig=${OUTPUT}/client.kubeconfig

 kubectl config set-context default \
    --cluster=kubernetes \
    --user=client \
    --kubeconfig=${OUTPUT}/client.kubeconfig

 kubectl config use-context default --kubeconfig=kubeconfig/client.kubeconfig

cp ${OUTPUT}/client.kubeconfig ${OUTPUT}/${REGION}-kubectl.kubeconfig
cp ${OUTPUT}/client.kubeconfig ${OUTPUT}/${REGION}-kubelet.kubeconfig
mv ${OUTPUT}/client.kubeconfig ${OUTPUT}/${REGION}-kubeproxy.kubeconfig
chmod 0644 ${OUTPUT}/*
