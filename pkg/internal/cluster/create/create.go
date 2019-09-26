package create

import (
	"fmt"
	"os"

	"github.com/gosoon/glog"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/create"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/context"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions"
	configaction "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions/config"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions/installcni"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions/installstorage"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions/kubeadminit"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions/kubeadmjoin"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions/waitforready"
	createtypes "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/types"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/delete"
	"github.com/gosoon/kubernetes-operator/pkg/internal/util/cli"
)

// Cluster creates a cluster
func Cluster(ctx *context.Context, options ...create.ClusterOption) error {
	// apply options, do defaulting etc.
	opts, err := collectOptions(options...)
	if err != nil {
		return err
	}

	fmt.Printf("clusterOptions:%+v \n", opts)

	if err := validate(opts); err != nil {
		return err
	}

	ctx.ClusterOptions = opts

	status := cli.NewStatus(os.Stdout)
	// pull docker image
	// attempt to explicitly pull the required node images if they doesn't exist locally
	// we don't care if this errors, we'll still try to run which also pulls
	ensureNodeImages(status, opts.NodeImage)

	// prepare node, copy docker image bin to local path and  为 node 和 master 生成配置文件
	// Create node containers implementing defined config Nodes
	if err := provisionNodes(status, ctx); err != nil {
		// In case of errors nodes are deleted (except if retain is explicitly set)
		glog.Error(err)

		// 执行失败清理
		// if exec failed and cleanup
		//if !opts.Retain {
		// TODO
		_ = delete.Cluster(ctx)
		//}
		return err
	}

	actionsToRun := []actions.Action{
		configaction.NewAction(), // setup kubeadm config  // 创建 kubeadm 配置 /kind/kubeadm.conf
	}

	if opts.SetupKubernetes {
		actionsToRun = append(actionsToRun,
			kubeadminit.NewAction(), // run kubeadm init   // start control plane
		)
		// this step might be skipped, but is next after init
		if !opts.Config.Networking.DisableDefaultCNI {
			actionsToRun = append(actionsToRun,
				installcni.NewAction(), // install CNI
			)
		}
		// add remaining steps
		actionsToRun = append(actionsToRun,
			installstorage.NewAction(),                // install StorageClass
			kubeadmjoin.NewAction(),                   // run kubeadm join
			waitforready.NewAction(opts.WaitForReady), // wait for cluster readiness
		)
	}

	actionsContext := actions.NewActionContext(opts, ctx.Server, ctx.Port, status)
	for _, action := range actionsToRun {
		if err := action.Execute(actionsContext); err != nil {
			//if !opts.Retain {
			_ = delete.Cluster(ctx)
			//}
			return err
		}
	}

	return nil

}

func collectOptions(options ...create.ClusterOption) (*createtypes.ClusterOptions, error) {
	// apply options
	opts := &createtypes.ClusterOptions{
		SetupKubernetes: true,
	}
	for _, option := range options {
		newOpts, err := option(opts)
		if err != nil {
			return nil, err
		}
		opts = newOpts
	}

	return opts, nil
}
