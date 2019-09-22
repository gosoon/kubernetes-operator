package create

import (
	"strings"

	"github.com/gosoon/glog"
	"github.com/gosoon/kubernetes-operator/pkg/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/container/docker"
	"github.com/gosoon/kubernetes-operator/pkg/exec"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/context"
	"github.com/gosoon/kubernetes-operator/pkg/internal/util/cli"
	"github.com/pkg/errors"
)

const (
	containerName = "kubernetes_base"
)

// provisionNodes takes care of creating all the containers
// that copy kubeadm,kubectl,kubelet binary and config to locally
func provisionNodes(
	status *cli.Status, ctx *context.Context,
) error {
	defer status.End(false)

	// create local install path
	cmd := exec.Command("mkdir", "-pv", constants.InstallPath)
	lines, err := exec.CombinedOutputLines(cmd)
	glog.Info(strings.Join(lines, "\n"))
	if err != nil {
		return errors.Wrap(err, "failed to create locally installer path")
	}

	// start base container
	if err := startContainer(ctx.ClusterOptions.NodeImage); err != nil {
		return err
	}

	// copy kubernetes binary to local install path
	// docker cp kubernetes_base:/kubernetes/ ./tmp/installer/
	if err := docker.CopyFrom(containerName, "/kubernetes", constants.InstallPath); err != nil {
		return err
	}

	status.End(true)
	return nil
}

func startContainer(image string) error {
	// start docker
	// docker run -it --privileged --name kubernetes_base registry.cn-hangzhou.aliyuncs.com/aliyun_kube_system/kubernetes:v1.15.3
	runArgs := []string{
		"--it",
		"--privileged",
		"--name",
		containerName,
	}
	err := docker.Run(
		image,
		docker.WithRunArgs(runArgs...),
	)

	return err
}
