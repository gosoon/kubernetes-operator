package app

import (
	"net"

	"github.com/gosoon/glog"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type flagpole struct {
	Port string
}

// NewServerCommand returns a new cobra.Command for kube-on-kube server
func NewServerCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "cluster",
		Short: "Creates a local Kubernetes cluster",
		Long:  "Creates a local Kubernetes cluster",
		Run: func(cmd *cobra.Command, args []string) {
			run(flags)
		},
	}
	cmd.Flags().StringVar(&flags.Port, "port", "10022", "installer agent grpc server port(default 10022)")
	return cmd
}

func run(flags *flagpole) {
	// start grpc server
	l, err := net.Listen("tcp", ":"+flags.Port)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	installer := NewInstaller(&Options{
		Flags:  flags,
		Server: server,
	})
	// register grpc server
	installerv1.RegisterInstallerServer(server, installer)

	go server.Serve(l)

	master := []ecsv1.Node{{IP: "127.0.0.1"}}
	cluster := &ecsv1.KubernetesCluster{Spec: ecsv1.KubernetesClusterSpec{Cluster: ecsv1.Cluster{MasterList: master}}}
	installer.DispatchClusterConfig(cluster)
}
