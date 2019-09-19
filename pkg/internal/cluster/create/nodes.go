package create

import "github.com/gosoon/kubernetes-operator/pkg/internal/apis/config"

// provisionNodes takes care of creating all the containers
// that copy kubeadm,kubectl,kubelet bin to locally
func provisionNodes(
	status *cli.Status, cfg *config.Cluster, clusterName, clusterLabel string,
) error {
	defer status.End(false)

	if err := copyBinaryToLocal(); err != nil {
		return err
	}

	status.End(true)
	return nil
}

func copyBinaryToLocal() {

	// start docker

	// copy binary

}
