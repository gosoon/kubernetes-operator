FROM ubuntu:19.04
MAINTAINER gosoon

# start a sshd service in image so that the container cannot exit
#

RUN apt-get update; apt-get upgrade -y \
    && apt-get install -y openssh-server 

RUN sed -i 's/UsePAM yes/UsePAM no/g' /etc/ssh/sshd_config 

RUN mkdir /var/run/sshd 

RUN mkdir -pv /kubernetes/bin \
    && mkdir -pv /kubernetes/manifests \
    && mkdir -pv /kubernetes/systemd

Add ./bin /kubernetes/bin/
Add ./manifests /kubernetes/manifests/
Add ./systemd /kubernetes/systemd/
Add ./version /kubernetes/
Add ./kubeadm.conf /kubernetes/

EXPOSE 22
ENTRYPOINT /usr/sbin/sshd -D
