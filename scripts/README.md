kubernetes 使用二进制部署，默认版本为 v1.14.0


### kube-master


### kube-node

##### kubelet 

--network-plugin=cni  设置启用CNI网络插件，因为后续是使用Calico网络，所以需要配置

##### kube-proxy
--cluster-cidr=10.1.0.0/16指定pod在kubernetes启动的虚拟IP网段(CNI网络),提供后续calico使用参数


##### coerdns
因为kuberntes中的所有pod都是基于service域名解析后，再负载均衡分发到service后端的各个pod服务中，那么如果没有DNS解析，则无法查到各个服务对应的service服务，以下举个例子。


首先第一步要知道集群使用的DNS的IP地址, kubelet config `--cluster-dns=10.0.6.200`,

### etcd

指定 etcd 的工作目录为 /var/lib/etcd，数据目录为 /var/lib/etcd，需在启动服务前创建这两个目录；




