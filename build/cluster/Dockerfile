FROM centos

RUN mkdir -p /home/kube
ADD ./kubernetes-operator /home/kubernetes-operator
ADD ./control.sh /home/kubernetes-operator
WORKDIR /home/kubernetes-operator
ENTRYPOINT ["/bin/bash","control.sh","start"]
