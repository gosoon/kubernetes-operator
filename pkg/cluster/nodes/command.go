package nodes

import (
	"github.com/gosoon/kubernetes-operator/pkg/exec"
)

// Command returns a new exec.Cmd that will run on the node
func Command(command string, args ...string) exec.Cmd {
	return exec.Command(command, args...)
}
