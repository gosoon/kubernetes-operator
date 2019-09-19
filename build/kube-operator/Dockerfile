FROM golang:1.12-alpine

RUN mkdir -p /home/kubernetes-operator/logs && 
    mkdir /home/kubernetes-operator/configs

ADD ./output/kubernetes-operator /home/kubernetes-operator
ADD ./build/operator/control.sh /home/kubernetes-operator \
    && chmod +x /home/kubernetes-operator/control.sh

ADD ./configs/kube-operator /home/kubernetes-operator/configs

WORKDIR /home/kubernetes-operator
ENTRYPOINT ["control.sh","start"]
