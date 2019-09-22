package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesCluster is the Schema for the kubernetesclusters API
type KubernetesCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubernetesClusterSpec   `json:"spec,omitempty"`
	Status KubernetesClusterStatus `json:"status,omitempty"`
}

// KubernetesClusterSpec defines the desired state of KubernetesCluster
type KubernetesClusterSpec struct {
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Cluster Cluster `json:"cluster,omitempty"`

	// Addons are some of the applications that need to be pre-installed in the cluster,eg: helm,promethus,logpolit...
	Addons Addons `json:"addons,omitempty"`
}

// Cluster xxx
type Cluster struct {
	TimeoutMins string `json:"timeout_mins,omitempty"`

	// ClusterType is a specified cluster,eg: kubernetes,k3s
	ClusterType ClusterType `json:"clusterType,omitempty"`

	// PodCIDR
	PodCIDR string `json:"podCIDR,omitempty"`

	// ServiceCIDR is apiserver and controller-manager flag `--service-cluster-ip-range`
	ServiceCIDR string `json:"serviceCIDR,omitempty"`

	// MasterList
	MasterList []Node `json:"masterList" tag:"required"`

	// ExternalLoadBalancer is a vip by lvs,haproxy or etc
	ExternalLoadBalancer string `json:"externalLoadBalancer,omitempty"`

	// NodeList
	NodeList []Node `json:"nodeList" tag:"required"`

	// EtcdList
	EtcdList []Node `json:"etcdList,omitempty"`

	// Region
	Region string `json:"region,omitempty"`

	// login destination host used authConfig
	AuthConfig AuthConfig `json:"authConfig,omitempty"`

	// kubernetes version
	KubeVersion string `json:"kubeVersion"`
}

// ClusterType is a specified cluster,eg: kubernetes,k3s,kind...
type ClusterType string

const (
	// KubernetesClusterType
	KubernetesClusterType ClusterType = "kubernetes"

	// K3sClusterType
	K3sClusterType ClusterType = "k3s"

	// kubeedge
	KubeedgeClusterType ClusterType = "kubeedge"

	// kind
	KindClusterType ClusterType = "kind"
)

// AuthConfig defines the nodes peer authentication
type AuthConfig struct {
	// Username
	Username string `json:"username,omitempty"`

	// Password
	Password string `json:"password,omitempty"`

	// PrivateSSHKey, use base64 encode
	PrivateSSHKey string `json:"privateSSHKey,omitempty"`
}

// Addons are some of the applications that need to be pre-installed in the Cluster
type Addons struct {
	// TODO
}

// KubernetesClusterStatus defines the observed state of KubernetesCluster
type KubernetesClusterStatus struct {
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Phase KubernetesOperatorPhase `json:"phase,omitempty"`

	// when job failed callback or job timeout used
	Reason string `json:"reason,omitempty"`

	// JobName is store each job name
	JobName string `json:"jobName,omitempty"`

	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesClusterList contains a list of KubernetesCluster
type KubernetesClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubernetesCluster `json:"items"`
}

// Node xxx
type Node struct {
	// IP
	IP string `json:"ip,omitempty"`

	// Role is used in kubeadm installer
	Role NodeRole `json:"role:omitempty"`
}

// NodeRole defines possible role for nodes in a Kubernetes cluster managed by `kind`
type NodeRole string

const (
	// ControlPlaneRole identifies a node that hosts a Kubernetes control-plane.
	// NOTE: in single node clusters, control-plane nodes act also as a worker
	// nodes, in which case the taint will be removed. see:
	// https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#control-plane-node-isolation
	ControlPlaneRole NodeRole = "control-plane"

	SecondaryControlPlaneRole NodeRole = "secondary-control-plane"

	// WorkerRole identifies a node that hosts a Kubernetes worker
	WorkerRole NodeRole = "worker"
)

// "None,Creating,Running,Failed,Scaling"
type KubernetesOperatorPhase string
