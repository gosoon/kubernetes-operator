package kubeadmjoin

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cloudflare/cfssl/log"
	"github.com/gosoon/kubernetes-operator/pkg/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/cluster/nodes"
	"github.com/gosoon/kubernetes-operator/pkg/exec"
	"github.com/gosoon/kubernetes-operator/pkg/fs"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions"
	"github.com/pkg/errors"
)

// Action implements action for creating the kubeadm join
// and deployng it on the bootrap control-plane node.
type Action struct{}

// NewAction returns a new action for creating the kubeadm jion
func NewAction() actions.Action {
	return &Action{}
}

// Execute runs the action
func (a *Action) Execute(ctx *actions.ActionContext) error {
	if ctx.Role == constants.SecondaryControlPlaneRole {
		allNodes := ctx.Nodes()
		if err := joinSecondaryControlPlanes(
			ctx, allNodes,
		); err != nil {
			return err
		}
	}

	// then join worker nodes if any
	if ctx.Role == constants.WorkerNodeRoleValue {
		if err := joinWorkers(ctx); err != nil {
			return err
		}
	}

	return nil
}

func joinSecondaryControlPlanes(
	ctx *actions.ActionContext,
	allNodes []nodes.Node,
) error {
	ctx.Status.Start("Joining more control-plane nodes 🎮")
	defer ctx.Status.End(false)

	// TODO(bentheelder): it's too bad we can't do this concurrently
	// (this is not safe currently)
	if err := runKubeadmJoinControlPlane(allNodes); err != nil {
		return err
	}

	ctx.Status.End(true)
	return nil
}

func joinWorkers(ctx *actions.ActionContext) error {
	ctx.Status.Start("Joining worker nodes 🚜")
	defer ctx.Status.End(false)

	// create the workers concurrently
	//fns := []func() error{}
	//for _, node := range workers {
	//node := node // capture loop variable
	//fns = append(fns, func() error {
	runKubeadmJoin()

	//if err := concurrent.UntilError(fns); err != nil {
	//return err
	//}

	ctx.Status.End(true)
	return nil
}

// runKubeadmJoinControlPlane executes kubadm join --control-plane command
func runKubeadmJoinControlPlane(allNodes []nodes.Node) error {
	// creates the folder tree for pre-loading necessary cluster certificates
	// on the joining node
	if err := exec.Command("mkdir", "-p", "/etc/kubernetes/pki/etcd").Run(); err != nil {
		return errors.Wrap(err, "failed to join node with kubeadm")
	}

	// define the list of necessary cluster certificates
	fileNames := []string{
		"ca.crt", "ca.key",
		"front-proxy-ca.crt", "front-proxy-ca.key",
		"sa.pub", "sa.key",
		// TODO(someone): if we gain external etcd support these will be
		// handled differently
		"etcd/ca.crt", "etcd/ca.key",
	}

	// TODO(gosoon): get ca files from controlPlane

	// creates a temporary folder on the host that should acts as a transit area
	// for moving necessary cluster certificates
	tmpDir, err := fs.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	err = os.MkdirAll(filepath.Join(tmpDir, "/etcd"), os.ModePerm)
	if err != nil {
		return err
	}
	// get the handle for the bootstrap control plane node (the source for necessary cluster certificates)
	controlPlaneHandle, err := nodes.BootstrapControlPlaneNode(allNodes)
	if err != nil {
		return err
	}

	// TODO(gosoon): use grpc transfer file
	// copies certificates from the bootstrap control plane node to the joining node
	for _, fileName := range fileNames {
		// sets the path of the certificate into a node
		containerPath := path.Join("/etc/kubernetes/pki", fileName)
		// set the path of the certificate into the tmp area on the host
		tmpPath := filepath.Join(tmpDir, fileName)
		// copies from bootstrap control plane node to tmp area
		if err := nodes.CopyFrom(containerPath, tmpPath); err != nil {
			return errors.Wrapf(err, "failed to copy certificate %s", fileName)
		}
		// copies from tmp area to joining node
		if err := nodes.CopyTo(tmpPath, containerPath); err != nil {
			return errors.Wrapf(err, "failed to copy certificate %s", fileName)
		}
	}
	return runKubeadmJoin()
}

// runKubeadmJoin executes kubadm join command
func runKubeadmJoin() error {
	// run kubeadm join
	// TODO(bentheelder): this should be using the config file
	cmd := exec.Command(
		"kubeadm", "join",
		// the join command uses the config file generated in a well known location
		"--config", "/tmp/kubeadm.conf",
		// preflight errors are expected, in particular for swap being enabled
		// TODO(bentheelder): limit the set of acceptable errors
		"--ignore-preflight-errors=all",
		// increase verbosity for debugging
		"--v=6",
	)
	lines, err := exec.CombinedOutputLines(cmd)
	log.Debug(strings.Join(lines, "\n"))
	if err != nil {
		return errors.Wrap(err, "failed to join node with kubeadm")
	}
	return nil
}
