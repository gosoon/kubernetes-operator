#!/bin/bash

# docker releases download docs:  https://kubernetes.io/docs/setup/production-environment/container-runtimes/
[ -e ./config.sh ] && . ./config.sh || exit

DOCKER_BIN_DIR="../bin/${DOCKER_VER}"
DOCKER_SYSTEMD_CONFIG_DIR="../systemd"
DEST_SYSTEMD_DIR="/usr/lib/systemd/system"

install_docker() {
    # Install Docker CE
    ## Set up the repository
    ### Install required packages.
    yum install yum-utils device-mapper-persistent-data lvm2

    ### Add docker repository.
    yum-config-manager \
        --add-repo \
        https://download.docker.com/linux/centos/docker-ce.repo

    ## Install docker ce.
    yum update && yum install ${DOCKER_VER}
}

if [ -f ${DOCKER_BIN_DIR}/dockerd ];then 
    cp ${DOCKER_BIN_DIR}/* /usr/bin/
    cp ${DOCKER_SYSTEMD_CONFIG_DIR}/docker.service  ${DEST_SYSTEMD_DIR}/
else
    install_docker  
fi

## Create /etc/docker directory.
mkdir /etc/docker

# Setup daemon.
cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2",
  "storage-opts": [
    "overlay2.override_kernel_check=true"
  ]
}
EOF

# Restart docker.
systemctl daemon-reload
systemctl restart docker
systemctl status docker

if [ $? -ne 0 ];then  
    echo "install ${DOCKER_VER} failed !!!" && exit 1
fi
