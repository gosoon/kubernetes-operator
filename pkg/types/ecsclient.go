package types

import ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"

// EcsClient xxx
type EcsClient struct {
	Name          string       `json:"name"`
	Namespace     string       `json:"namespace"`
	Region        string       `json:"region"`
	TimeoutMins   string       `json:"timeoutMins"`
	ClusterType   string       `json:"clusterType"`
	ContainerCIDR string       `json:"containerCIDR"`
	ServiceCIDR   string       `json:"serviceCIDR"`
	MasterList    []ecsv1.Node `json:"masterList"`
	NodeList      []ecsv1.Node `json:"nodeList"`
	EtcdList      []ecsv1.Node `json:"etcdList"`
	Retry         bool         `json:"retry"`
}

// CallBack xxx
type CallBack struct {
	Name       string       `json:"name"`
	Namespace  string       `json:"namespace"`
	Region     string       `json:"region"`
	MasterList []ecsv1.Node `json:"masterList"`
	NodeList   []ecsv1.Node `json:"nodeList"`
	EtcdList   []ecsv1.Node `json:"etcdList"`
	KubeConfig string       `json:"kubeconfig"`
	Success    bool         `json:"success"`
	Message    string       `json:"message"`
}
