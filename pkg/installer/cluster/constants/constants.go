package constants

// DefaultClusterName is the default cluster Context name
const DefaultClusterName = "default"

const (
	// ExternalLoadBalancerNodeRoleValue identifies a node that hosts an
	// external load balancer for the API server in HA configurations.
	//
	// Please note that `kind` nodes hosting external load balancer are not
	// kubernetes nodes
	ExternalLoadBalancerNodeRoleValue string = "external-load-balancer"

	// ExternalEtcdNodeRoleValue identifies a node that hosts an external-etcd
	// instance.
	//
	// WARNING: this node type is not yet implemented!
	//
	// Please note that `kind` nodes hosting external etcd are not
	// kubernetes nodes
	ExternalEtcdNodeRoleValue string = "external-etcd"

	// InstallPath is write kubeadm config default path
	InstallPath string = "/tmp/install/"
)
