package agent

import (
	"net"
	"testing"

	"github.com/gosoon/glog"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/constants"

	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInstallCluster(t *testing.T) {
	flag := &flagpole{
		Port:     "10023",
		Registry: "registry.cn-hangzhou.aliyuncs.com/aliyun_kube_system",
	}
	// start grpc server
	l, err := net.Listen("tcp", ":"+"10023")
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	// create a cluster context and create the cluster
	ctx := cluster.NewContext(constants.DefaultClusterName, server, "10023")

	installer := NewInstaller(&Options{
		Flags:   flag,
		Context: ctx,
		Server:  server,
	})

	// register grpc server
	installerv1.RegisterInstallerServer(server, installer)

	go func() {
		glog.Fatal(server.Serve(l))
	}()

	kubernetesCluster := &ecsv1.KubernetesCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KubernetesCluster",
			APIVersion: "ecs.yun.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Spec: ecsv1.KubernetesClusterSpec{
			Cluster: ecsv1.Cluster{
				ClusterType:          ecsv1.KubernetesClusterType,
				PodCIDR:              "192.168.0.0/16",
				ServiceCIDR:          "10.233.0.0/18",
				MasterList:           []ecsv1.Node{{IP: "192.168.72.224", Role: ecsv1.ControlPlaneRole}},
				ExternalLoadBalancer: "127.0.0.1",
				Region:               "default",
				KubeVersion:          "v1.15.3",
			},
		},
	}

	err = installer.DoInstallCluster(kubernetesCluster)
	if err != nil {
		glog.Error(err)
	}

}
