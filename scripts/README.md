Scripts dir is use binary deploy kubernetes cluster,default version is v1.14.0.Kubernetes-operator use ansible call this scripts and deploy kubernetes cluster.

use scripts deploy kubernetes:

1. define version and hosts in scripts/deploy/config.sh

2. exec scripts to deploy
 ```
$ cd deploy/
$ bash deploy.sh etcd
$ bash deploy.sh master
$ bash deploy.sh node
 ```
