package create

import (
	"fmt"
	"path"
	"strings"

	"github.com/gosoon/glog"
	"github.com/gosoon/kubernetes-operator/pkg/container/docker"
	"github.com/gosoon/kubernetes-operator/pkg/exec"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/nodes"
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
	// TODO : if kubernetes_base is started and continue
	//startC,ntainer(ctx.ClusterOptions.NodeImage)

	defer func() {
		//_ = exec.Command("docker", "rm", "-f", "-v", containerID).Run()
	}()

	// copy kubernetes binary to local install path
	// docker cp kubernetes_base:/kubernetes/ ./tmp/installer/
	if err := docker.CopyFrom(containerName, "/kubernetes", constants.InstallPath); err != nil {
		return err
	}

	fmt.Println("copy bin ", constants.InstallPath+"kubernetes/bin/")
	// copy kubectl,kubelet,kubeadm to /usr/bin/
	err = exec.Command("cp", "-r", constants.InstallPath+"kubernetes/bin/*", "/usr/bin/").Run()
	if err != nil {
		glog.Error("copy bin failed :", err)
		return err
	}
	fmt.Println("copy bin success :", constants.InstallPath+"/kubernetes/bin/")

	// setup kubelet systemd
	// create the kubelet service
	kubeletService := nodes.KubeletServicePath
	if err := createFile(kubeletService, nodes.KubeletServiceContents); err != nil {
		return errors.Wrap(err, "failed to create kubelet service file")
	}

	// enable the kubelet service
	if err := exec.Command("systemctl", "enable", kubeletService).Run(); err != nil {
		return errors.Wrap(err, "failed to enable kubelet service")
	}

	// setup the kubelet dropin
	kubeletDropin := nodes.KubeletServiceConfigDir + "10-kubeadm.conf"
	if err := createFile(kubeletDropin, nodes.Kubeadm10conf); err != nil {
		return errors.Wrap(err, "failed to configure kubelet service")
	}

	// ensure we don't fail if swap is enabled on the host
	if err = exec.Command("/bin/sh", "-c",
		`echo "KUBELET_EXTRA_ARGS=--fail-swap-on=false" > /etc/default/kubelet`,
	).Run(); err != nil {
		glog.Errorf("Image build Failed! Failed to add kubelet extra args: %v", err)
		return err
	}

	status.End(true)
	return nil
}

func startContainer(image string) error {
	// start docker
	// docker run -d --name kubernetes_base registry.cn-hangzhou.aliyuncs.com/aliyun_kube_system/kubernetes:v1.15.3
	runArgs := []string{
		"-d",
		"--name",
		containerName,
	}
	err := docker.Run(
		image,
		docker.WithRunArgs(runArgs...),
	)

	return err
}

func createFile(filePath string, contents string) error {
	// ensure the directory first
	if err := exec.Command("mkdir", "-pv", path.Dir(filePath)).Run(); err != nil {
		return err
	}

	return exec.Command("cp", "/dev/stdin", filePath).SetStdin(
		strings.NewReader(contents),
	).Run()
}
