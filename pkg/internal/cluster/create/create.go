package create

import (
	"github.com/gosoon/glog"
	"github.com/gosoon/kubernetes-operator/pkg/cluster/create"
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
)

// Cluster creates a cluster
func Cluster(ctx *context.Context, options ...create.ClusterOption) error {
	// apply options, do defaulting etc.
	opts, err := collectOptions(options...)
	if err != nil {
		return err
	}

	// then validate 检查 docker 暴露端口等信息
	//if err := opts.Config.Validate(); err != nil {
	//return err
	//}

	// 拉取 master node image and docker cp bin to local
	// attempt to explicitly pull the required node images if they doesn't exist locally
	// we don't care if this errors, we'll still try to run which also pulls
	ensureNodeImages(status, opts.Config)

	// 4.prepare node, copy docker image bin to local and  为 node 和 master 生成配置文件
	// Create node containers implementing defined config Nodes
	if err := provisionNodes(status, opts.Config, ctx.Name(), ctx.ClusterLabel()); err != nil {
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
		//loadbalancer.NewAction(), // setup external loadbalancer  // 检查是否有 loadbalancer
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

	// run all actions
	// TODO: compare node list and get local ip
	localIP := getLocalIP(opts.Config)
	if localIP == "" {
		return nil
	}
	role = ""

	actionsContext := actions.NewActionContext(opts.Config, ctx, status, localIP, role)
	for _, action := range actionsToRun {
		if err := action.Execute(actionsContext); err != nil {
			//if !opts.Retain {
			_ = delete.Cluster(ctx)
			//}
			return err
		}
	}

	// print how to set KUBECONFIG to point to the cluster etc.
	printUsage(ctx.Name())

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

	// do post processing for options
	// first ensure we at least have a default cluster config
	//if opts.Config == nil {
	//cfg, err := encoding.Load("")
	//if err != nil {
	//return nil, err
	//}
	//opts.Config = cfg
	//}

	// if NodeImage was set, override the image on all nodes
	if opts.NodeImage != "" {
		// Apply image override to all the Nodes defined in Config
		// TODO(fabrizio pandini): this should be reconsidered when implementing
		//     https://github.com/kubernetes-sigs/kind/issues/133
		for i := range opts.Config.Nodes {
			opts.Config.Nodes[i].Image = opts.NodeImage
		}
	}

	// default config fields (important for usage as a library, where the config
	// may be constructed in memory rather than from disk)
	//encoding.Scheme.Default(opts.Config)

	return opts, nil
}
