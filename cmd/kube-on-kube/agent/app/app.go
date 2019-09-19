package app

import (
	"net"
	"time"

	"github.com/gosoon/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/gosoon/kubernetes-operator/pkg/cluster"
	"github.com/gosoon/kubernetes-operator/pkg/cluster/create"
)

type flagpole struct {
	Name      string
	Config    string
	ImageName string
	Retain    bool
	Wait      time.Duration
	Port      int
}

// NewServerCommand returns a new cobra.Command for kube-on-kube server
func NewServerCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "cluster",
		Short: "Creates a local Kubernetes cluster",
		Long:  "Creates a local Kubernetes cluster using Docker container 'nodes'",
		Run: func(cmd *cobra.Command, args []string) error {
			run(flags)
		},
	}
	cmd.Flags().StringVar(&flags.Config, "config", "", "path to a kind config file")
	cmd.Flags().StringVar(&flags.ImageName, "image", "", "node docker image to use for booting the cluster")
	cmd.Flags().DurationVar(&flags.Wait, "wait", time.Duration(0), "Wait for control plane node to be ready (default 0s)")
	cmd.Flags().DurationVar(&flags.Port, "port", 10023, "installer agent grpc server port(default 10023)")
	return cmd
}

type Server struct{}

func run(flags *flagpole) {
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	// register grpc server
	pb.RegisterScpServer(server, &Server{})
	glog.Fatal(server.Serve(l))

	// TODO :receiver operator ecs config

	// create a cluster context and create the cluster
	ctx := cluster.NewContext(flags.Name)
	glog.Infof("Creating cluster %q ...\n", flags.Name)
	if err = ctx.Create(
		create.WithConfig(flags.Config),
		//create.WithNodeImage(flags.ImageName),
		//create.Retain(flags.Retain),
		create.WaitForReady(flags.Wait),
	); err != nil {
		//if utilErrors, ok := err.(util.Errors); ok {
		//for _, problem := range utilErrors.Errors() {
		//glog.Error(problem)
		//}
		//return errors.New("aborting due to invalid configuration")
		//}
		return errors.Wrap(err, "failed to create cluster")
	}

	return nil
}
