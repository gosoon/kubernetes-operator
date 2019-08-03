package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:method=GetScale,verb=get,subresource=scale,result=k8s.io/kubernetes/pkg/apis/autoscaling.Scale
// +genclient:method=UpdateScale,verb=update,subresource=scale,input=k8s.io/kubernetes/pkg/apis/autoscaling.Scale,result=k8s.io/kubernetes/pkg/apis/autoscaling.Scale

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
	TimeoutMins   string     `json:"timeout_mins,omitempty"`
	ClusterType   string     `json:"clusterType,omitempty"`
	ContainerCIDR string     `json:"containerCIDR,omitempty"`
	ServiceCIDR   string     `json:"serviceCIDR,omitempty"`
	MasterList    []Node     `json:"masterList" tag:"required"`
	MasterVIP     string     `json:"masterVIP,omitempty"`
	NodeList      []Node     `json:"nodeList" tag:"required"`
	EtcdList      []Node     `json:"etcdList,omitempty"`
	Region        string     `json:"region,omitempty"`
	AuthConfig    AuthConfig `json:"authConfig,omitempty"`
}

// AuthConfig defines the nodes peer authentication
type AuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	PrivateSSHKey string `json:"privateSSHKey,omitempty"`
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

// users
// "None,Creating,Running,Failed,Scaling"
type KubernetesOperatorPhase string

type Node struct {
	IP string `json:"ip,omitempty"`
}
