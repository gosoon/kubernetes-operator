package docker

import (
	"github.com/gosoon/kubernetes-operator/pkg/exec"
)

// CopyTo copies the file at hostPath to the container at destPath
func CopyTo(hostPath, containerNameOrID, destPath string) error {
	cmd := exec.Command(
		"docker", "cp",
		hostPath,                       // from the source file
		containerNameOrID+":"+destPath, // to the node, at dest
	)
	return cmd.Run()
}

// CopyFrom copies the file or dir in the container at srcPath to the host at hostPath
func CopyFrom(containerNameOrID, srcPath, hostPath string) error {
	cmd := exec.Command(
		"docker", "cp",
		containerNameOrID+":"+srcPath, // from the node, at src
		hostPath,                      // to the host
	)
	return cmd.Run()
}
