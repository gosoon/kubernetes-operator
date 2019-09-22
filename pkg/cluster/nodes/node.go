package nodes

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"time"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/exec"
	"github.com/gosoon/kubernetes-operator/pkg/internal/util/env"

	"github.com/pkg/errors"
)

// ConfigNodeAddressAndRole is return the host ip and role,get host ip and role from kubernetes cluster
func ConfigNodeAddressAndRole(cluster *ecsv1.KubernetesCluster) (string, ecsv1.NodeRole) {
	var hosts []string
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				hosts = append(hosts, ipnet.IP.String())
			}
		}
	}

	var clusterNodeList []ecsv1.Node
	clusterNodeList = append(clusterNodeList, cluster.Spec.Cluster.NodeList...)
	clusterNodeList = append(clusterNodeList, cluster.Spec.Cluster.MasterList...)

	for _, node := range clusterNodeList {
		for _, host := range hosts {
			if node.IP == host {
				return node.IP, node.Role
			}
		}
	}
	return "", ""
}

// ConstructImage is return the node image,image representative kubernetes version and contain kubelet,
// kubectl,kubeadm binary
// defaule image is registry.cn-hangzhou.aliyuncs.com/aliyun_kube_system/kubernetes:v1.xxx
func ConstructImage(registry string, cluster *ecsv1.KubernetesCluster) string {
	return registry + "/kubernetes" + ":" + cluster.Spec.Cluster.KubeVersion
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

// KubeConfigPath returns the path to where the Kubeconfig would be placed
// by kind based on the configuration.
func KubeConfigPath(clusterName string) string {
	// configDir matches the standard directory expected by kubectl etc
	configDir := filepath.Join(env.HomeDir(), ".kube")
	// note that the file name however does not, we do not want to overwrite
	// the standard config, though in the future we may (?) merge them
	fileName := fmt.Sprintf("config-%s", clusterName)
	return filepath.Join(configDir, fileName)
}

// use grpc copy file
func CopyFrom(sourcePath string, destPath string) error {
	return nil
}

func CopyTo(sourcePath string, destPath string) error {
	return nil
}
