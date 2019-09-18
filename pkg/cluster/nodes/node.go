package nodes

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/gosoon/kubernetes-operator/pkg/exec"
	"github.com/pkg/errors"
)

// Node represents a handle to a kind node
// This struct must be created by one of: CreateControlPlane
// It should not be manually instantiated
// Node impleemnts exec.Cmder
type Node struct {
	IP    string
	Role  string
	Image string
}

// control plane nodes are "Ready".
func WaitForReady(until time.Time) bool {
	return tryUntil(until, func() bool {
		cmd := exec.Command(
			"kubectl",
			"--kubeconfig=/etc/kubernetes/admin.conf",
			"get",
			"nodes",
			"--selector=node-role.kubernetes.io/master",
			// When the node reaches status ready, the status field will be set
			// to true.
			"-o=jsonpath='{.items..status.conditions[-1:].status}'",
		)
		lines, err := exec.CombinedOutputLines(cmd)
		if err != nil {
			return false
		}

		// 'lines' will return the status of all nodes labeled as master. For
		// example, if we have three control plane nodes, and all are ready,
		// then the status will have the following format: `True True True'.
		status := strings.Fields(lines[0])
		for _, s := range status {
			// Check node status. If node is ready then this wil be 'True',
			// 'False' or 'Unkown' otherwise.
			if !strings.Contains(s, "True") {
				return false
			}
		}
		return true
	})
}

// helper that calls `try()`` in a loop until the deadline `until`
// has passed or `try()`returns true, returns wether try ever returned true
func tryUntil(until time.Time, try func() bool) bool {
	for until.After(time.Now()) {
		if try() {
			return true
		}
	}
	return false
}

// WriteFile writes content to dest on the node
func WriteFile(dest, content string) error {
	// create destination directory
	cmd := exec.Command("mkdir", "-p", filepath.Dir(dest))
	err := exec.RunLoggingOutputOnFail(cmd)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", dest)
	}

	return exec.Command("cp", "/dev/stdin", dest).SetStdin(strings.NewReader(content)).Run()
}

// use grpc copy file
func CopyFrom(sourcePath, destPath) {

}

func CopyTo(sourcePath, destPath) {

}
