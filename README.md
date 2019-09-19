# kubernetes-operator

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/gosoon/kubernetes-operator)](https://goreportcard.com/report/github.com/gosoon/kubernetes-operator)

kubernetes-operator is a control plane and manage all kubernetes cluster lifecycle.

<img src="doc/images/image-20190805195135765.png"></img>

## Introduce

**kubernetes-operator contains several large parts**：

- kubernetes-proxy: is a proxy and all requests pass through，look like a gateway.
- kube-operator: is a kubernetes operator deploy in meta kubernetes and manage all kubernetes clusters(create、scale、delete).
- cluster deploy: At present, the main focus is on kube-on-kube,and use ansible deploy kubernetes cluster and multiple kubernetes cluster applications, eg: metric-server、 promethus、log-polit、helm...
- node-operator: Manage cluster upgrade, fault self-healing, etc.

**The kubernetes-operator goal**：

- kube-on-kube：Similar to Ant Financial kube-on-kube-operator.
- kube-to-kube：The difference with `kube-on-kube`  is that the master component of the business cluster is not in the meta cluster.
- kube-to-cloud-kube : Manage clusters on multiple cloud environments.


## Getting started

First you need to make two images，one is kubernetes-operator，the other one is ansibleinit。And deploy kubernetes-operator in your kubernetes-cluster，if you don't have a kubernetes cluster，please see `scripts/REAEME.md` and deploy one.

```
$ make images

// deploy crd
$ kubectl create -f deploy/crds/ecs_v1_kubernetescluster_crd.yaml

// update your image address and deploy kubernetes-operator
$ kubectl create -f deploy/role.yaml  
$ kubectl create -f deploy/kube-operator.yaml

// CustomResources is KubernetesCluster,"ecs" for short in the kubernetes-operator
$ kubectl get crd
NAME                                  CREATED AT
kubernetesclusters.ecs.yun.com        2019-08-05T12:23:52Z

// update operator server in create_ecs.sh and create a cluster
$ bash test/create_ecs.sh  

$ kubectl get ecs
NAME           STATUS        AGE
test-cluster   Prechecking   45m
```



## Development Plan

1. Support deploy k3s、kubeedge cluster
2. support use kops deploy cluster
3. support for multiple version deploy
4. development node-operator 
5. support admission control

## Detailed instructions

- [kube-on-kube-operator 开发(一)](http://blog.tianfeiyu.com/2019/08/05/kube_on_kube_operator_1/)
- [kube-on-kube-operator 开发(二)](http://blog.tianfeiyu.com/2019/08/07/kube_on_kube_operator_2/)
- [kube-on-kube-operator 开发(三)](http://blog.tianfeiyu.com/2019/09/01/kube_on_kube_operator_3/)


