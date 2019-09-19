package create

import (
	"fmt"
	"strings"

	"github.com/gosoon/kubernetes-operator/pkg/internal/apis/config"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/kind/pkg/container/docker"
)

// ensureNodeImages ensures that the node images used by the create
// configuration are present
func ensureNodeImages(status *cli.Status, cfg *config.Cluster) {
	// pull each required image
	for _, image := range requiredImages(cfg).List() {
		// prints user friendly message
		if strings.Contains(image, "@sha256:") {
			image = strings.Split(image, "@sha256:")[0]
		}
		status.Start(fmt.Sprintf("Ensuring node image (%s) ", image))

		// attempt to explicitly pull the image if it doesn't exist locally
		// we don't care if this errors, we'll still try to run which also pulls
		_, _ = docker.PullIfNotPresent(image, 3)
	}
}

// requiredImages returns the set of images specified by the config
func requiredImages(cfg *config.Cluster) sets.String {
	images := sets.NewString()
	for _, node := range cfg.Nodes {
		images.Insert(node.Image)
	}
	return images
}
