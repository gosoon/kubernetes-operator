FROM centos:centos7

RUN mkdir -p /home/data1/kubernetes-operator/logs
ADD ./kubernetes-operator /home/data1/kubernetes-operator
ADD ./control.sh /home/data1/kubernetes-operator
WORKDIR /home/data1/kubernetes-operator
ENTRYPOINT ["/bin/bash","control.sh","start"]
