#!/bin/bash

[ -e ./config.sh ] && . ./config.sh || exit

CALICO_CONFIG="../yaml/calico_${CALICO_VER}/calico-etcd.yaml"

download_calico(){
    wget -O ${CALICO_CONFIG}  https://docs.projectcalico.org/${CALICO_VER}/manifests/calico-etcd.yaml
}

# config etcd server
etcd_num=$(echo ${ETCD_HOSTS} | awk -F ',' '{print NF}')
etcd_cluster=""
for i in `seq 1 ${etcd_num}`;do
	ip=$(echo ${ETCD_HOSTS} | awk -v idx=$i -F ',' '{print $idx}')
    cluster=$(echo "https://${ip}:2379")
    if [ $i -ne ${etcd_num} ];then
        cluster="${cluster},"
    fi
    etcd_cluster="${etcd_cluster}${cluster}"
done

# update secrets etcd info  
etcd_key_base64=$(cat /etc/etcd/ssl/etcd-server-key.pem | base64 -w 0)
etcd_cert_base64=$(cat /etc/etcd/ssl/etcd-server.pem | base64 -w 0)
etcd_ca_base64=$(cat /etc/etcd/ssl/ca.pem | base64 -w 0)
sed -i -e "s@# etcd-key:.*@etcd-key: ${etcd_key_base64}@g" ${CALICO_CONFIG}
sed -i -e "s@# etcd-cert:.*@etcd-cert: ${etcd_cert_base64}@g" ${CALICO_CONFIG}
sed -i -e "s@# etcd-ca:.*@etcd-ca: ${etcd_ca_base64}@g" ${CALICO_CONFIG}

# update configmap etcd info
sed -i -e 's#etcd_endpoints:.*#etcd_endpoints: \"${etcd_cluster}\"#g' ${CALICO_CONFIG}
sed -i -e 's#etcd_ca:.*#etcd_ca: \"/calico-secrets/etcd-ca\"#g' ${CALICO_CONFIG}
sed -i -e 's#etcd_cert:.*#etcd_cert: \"/calico-secrets/etcd-cert\"#g' ${CALICO_CONFIG}
sed -i -e 's#etcd_key:.*#etcd_key: \"/calico-secrets/etcd-key\"#g' ${CALICO_CONFIG}

# TODO : check download_calico
kubectl apply -f ${CALICO_CONFIG}
if [ $? -ne 0 ];then  
    echo "deploy calico failed !!!" && exit 1
fi
