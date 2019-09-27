package installer

// Installer is a interface define all installer implementation operation
type Installer interface {
	// ClusterNew is creating a new cluster
	ClusterNew(cluster *ecsv1.KubernetesCluster) error

	// ClusterScaleUp is scale up a cluster node
	ClusterScaleUp(cluster *ecsv1.KubernetesCluster, scaleUpNodeList []ecsv1.Node) error

	// ClusterScaleDown is scale down a cluster node
	ClusterScaleDown(cluster *ecsv1.KubernetesCluster, scaleDonwNodeList []ecsv1.Node) error

	// ClusterTerminating is delete a cluster
	ClusterTerminating(cluster *ecsv1.KubernetesCluster) error
}
