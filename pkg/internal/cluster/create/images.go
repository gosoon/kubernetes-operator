package create

import (
	"fmt"
	"strings"

	"github.com/gosoon/kubernetes-operator/pkg/container/docker"
	"github.com/gosoon/kubernetes-operator/pkg/internal/util/cli"
)

// ensureNodeImages ensures that the node images used by the create
// configuration are present
func ensureNodeImages(status *cli.Status, image string) {
	// prints user friendly message
	if strings.Contains(image, "@sha256:") {
		image = strings.Split(image, "@sha256:")[0]
	}
	status.Start(fmt.Sprintf("Ensuring node image (%s) ", image))

	fmt.Println("pull image ...")
	// attempt to explicitly pull the image if it doesn't exist locally
	// we don't care if this errors, we'll still try to run which also pulls
	_, _ = docker.PullIfNotPresent(image, 3)
}
