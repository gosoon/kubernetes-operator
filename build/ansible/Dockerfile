#FROM golang:1.12-alpine

#RUN apk add --no-cache \
    #wget \
    #make \
    #git \
    #openssh-client \
    #ansible \
    #openssh \
    #vim
FROM golang:alpine

RUN mkdir -p /home/kubernetes-operator \
    && mkdir -p /home/kubernetes-operator/configs/ansible

ENV OPERATION creating

#ADD ./scripts/ /home/kubernetes-operator/scripts
ADD ./ansible/ /home/kubernetes-operator/ansible
ADD ./output/ansibleinit /home/kubernetes-operator/

ADD ./configs/ansible/ansibleinit.toml /home/kubernetes-operator/configs/ansible/

WORKDIR /home/kubernetes-operator
#ENTRYPOINT ["./ansibleinit","-config","configs/ansible/ansibleinit.toml"]
