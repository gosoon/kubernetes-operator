FROM centos

RUN mkdir -p /home/kubernetes-operator/logs
ADD ./output/kubernetes-operator /home/kubernetes-operator
ADD ./build/operator/control.sh /home/kubernetes-operator
WORKDIR /home/kubernetes-operator
ENTRYPOINT ["/bin/bash","control.sh","start"]
