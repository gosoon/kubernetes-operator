package kubeadminit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gosoon/glog"
	"github.com/gosoon/kubernetes-operator/pkg/exec"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/kubeadm"
	"github.com/pkg/errors"
)

// kubeadmInitAction implements action for executing the kubadm init
// and a set of default post init operations like e.g. install the
// CNI network plugin.
type action struct{}

// NewAction returns a new action for kubeadm init
func NewAction() actions.Action {
	return &action{}
}

// Execute runs the action
func (a *action) Execute(ctx *actions.ActionContext) error {
	ctx.Status.Start("Starting control-plane 🕹️ ")
	defer ctx.Status.End(false)

	//allNodes, err := ctx.Nodes()
	//if err != nil {
	//return err
	//}

	// get the target node for this task
	//node, err := nodes.BootstrapControlPlaneNode(allNodes)
	//if err != nil {
	//return err
	//}

	// run kubeadm
	cmd := exec.Command(
		// init because this is the control plane node
		"kubeadm", "init",
		// preflight errors are expected, in particular for swap being enabled
		// TODO(bentheelder): limit the set of acceptable errors
		"--ignore-preflight-errors=all",
		// specify our generated config file
		"--config=/tmp/kubeadm.conf",
		"--skip-token-print",
		// increase verbosity for debugging
		"--v=6",
	)

	lines, err := exec.CombinedOutputLines(cmd)
	glog.Info(strings.Join(lines, "\n"))
	if err != nil {
		return errors.Wrap(err, "failed to init node with kubeadm")
	}

	kubeConfigPath := ctx.ClusterContext.KubeConfigPath()
	hostPort := kubeadm.APIServerPort
	if err := writeKubeConfig(kubeConfigPath, ctx.Config.Networking.APIServerAddress, hostPort); err != nil {
		return errors.Wrap(err, "failed to get kubeconfig from node")
	}

	// if we are only provisioning one node, remove the master taint
	// https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#master-isolation

	// mark success
	ctx.Status.End(true)
	return nil
}

// matches kubeconfig server entry like:
//    server: https://172.17.0.2:6443
// which we rewrite to:
//    server: https://$ADDRESS:$PORT
var serverAddressRE = regexp.MustCompile(`^(\s+server:) https://.*:\d+$`)

// writeKubeConfig writes a fixed KUBECONFIG to dest
// this should only be called on a control plane node
// While copyng to the host machine the control plane address
// is replaced with local host and the control plane port with
// a randomly generated port reserved during node creation.
func writeKubeConfig(dest string, hostAddress string, hostPort int) error {
	cmd := exec.Command("cat", "/etc/kubernetes/admin.conf")
	lines, err := exec.CombinedOutputLines(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to get kubeconfig from node")
	}

	// fix the config file, swapping out the server for the forwarded localhost:port
	var buff bytes.Buffer
	for _, line := range lines {
		match := serverAddressRE.FindStringSubmatch(line)
		if len(match) > 1 {
			addr := net.JoinHostPort(hostAddress, fmt.Sprintf("%d", hostPort))
			line = fmt.Sprintf("%s https://%s", match[1], addr)
		}
		buff.WriteString(line)
		buff.WriteString("\n")
	}

	// create the directory to contain the KUBECONFIG file.
	// 0755 is taken from client-go's config handling logic: https://github.com/kubernetes/client-go/blob/5d107d4ebc00ee0ea606ad7e39fd6ce4b0d9bf9e/tools/clientcmd/  loader.go#L412
	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create kubeconfig output directory")
	}

	return ioutil.WriteFile(dest, buff.Bytes(), 0600)
}
