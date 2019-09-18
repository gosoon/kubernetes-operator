package app

import (
	"time"

	"github.com/gosoon/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/gosoon/kubernetes-operator/pkg/cluster"
	"github.com/gosoon/kubernetes-operator/pkg/cluster/create"
)

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
	return cmd
}

func run(flags *flagpole) {

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
