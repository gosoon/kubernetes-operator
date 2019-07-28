FROM golang:1.12-alpine

RUN mkdir -p /home/kubernetes-operator/logs
ADD ./output/kubernetes-operator /home/kubernetes-operator
ADD ./build/operator/control.sh /home/kubernetes-operator \
    && chmod +x /home/kubernetes-operator/control.sh
WORKDIR /home/kubernetes-operator
ENTRYPOINT ["control.sh","start"]
