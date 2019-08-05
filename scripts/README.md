Scripts dir is use binary deploy kubernetes cluster,default version is v1.14.0.Kubernetes-operator use ansible call this scripts and deploy kubernetes cluster.

use scripts deploy kubernetes:

1. clone kubernetes-operator scripts in localhost `/home/kubernetes-operator/scripts` path
2. define version and hosts in `deploy/config.sh`,default bin file in `bin/`,the `bin/` dir file is mock, and you need to replace with `https://github.com/gosoon/kubernetes-utils/tree/master/scripts/bin` 
3. define host list in `deploy/hosts_env`
4. configure host login using the private key，put the public key on all hosts，and save ssh private-key in `/home/kubernetes-operator/private-key`
5. exec scripts to deploy

```
$ cd deploy/
$ bash deploy.sh etcd
$ bash deploy.sh master
$ bash deploy.sh node
```
