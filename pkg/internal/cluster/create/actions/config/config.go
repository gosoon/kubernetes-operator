package config

import (
	"fmt"

	"github.com/gosoon/glog"
	"github.com/gosoon/kubernetes-operator/pkg/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/cluster/nodes"
	"github.com/gosoon/kubernetes-operator/pkg/internal/apis/config"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/kubeadm"
	"github.com/pkg/errors"
)

// Action implements action for creating the kubeadm config
// and deployng it on the bootrap control-plane node.
type Action struct{}

// NewAction returns a new action for creating the kubadm config
func NewAction() actions.Action {
	return &Action{}
}

// Execute runs the action
func (a *Action) Execute(ctx *actions.ActionContext) error {
	ctx.Status.Start("Creating kubeadm config")
	defer ctx.Status.End(false)

	kubeVersion := ctx.Config.KubeVersion

	if ctx.Config.ExternalLoadBalancer == "" {

	}

	// get the control plane endpoint, in case the cluster has an external load balancer in
	// front of the control-plane nodes
	//controlPlaneEndpoint, controlPlaneEndpointIPv6, err := nodes.GetControlPlaneEndpoint(allNodes)
	//if err != nil {
	//return err
	//}

	// create kubeadm init config
	fns := []func() error{}

	configData := kubeadm.ConfigData{
		//ClusterName:          ctx.ClusterContext.Name(),
		KubernetesVersion:    kubeVersion,
		ControlPlaneEndpoint: ctx.LocalIP,
		APIBindPort:          kubeadm.APIServerPort,
		APIServerAddress:     ctx.Config.Networking.APIServerAddress,
		Token:                kubeadm.Token,
		PodSubnet:            ctx.Config.Networking.PodSubnet,
		ServiceSubnet:        ctx.Config.Networking.ServiceSubnet,
		ControlPlane:         true,
		IPv6:                 ctx.Config.Networking.IPFamily == "ipv6",
		NodeAddress:          ctx.LocalIP,
	}

	if ctx.Role == constants.WorkerNodeRoleValue {
		configData.ControlPlane = false
	}

	fns = append(fns, func() error {
		return writeKubeadmConfig(ctx.Config, configData)
	})

	return nil
}

// getKubeadmConfig generates the kubeadm config contents for the cluster
// by running data through the template.
func getKubeadmConfig(cfg *config.Cluster, data kubeadm.ConfigData) (path string, err error) {
	// generate the config contents
	config, err := kubeadm.Config(data)
	if err != nil {
		return "", err
	}
	// fix all the patches to have name metadata matching the generated config
	//patches, jsonPatches := setPatchNames(
	//allPatchesFromConfig(cfg),
	//)
	// apply patches
	// TODO(bentheelder): this does not respect per node patches at all
	// either make patches cluster wide, or change this
	//patched, err := kustomize.Build([]string{config}, patches, jsonPatches)
	//if err != nil {
	//return "", err
	//}
	return config, nil
}

func allPatchesFromConfig(cfg *config.Cluster) (patches []string, jsonPatches []config.PatchJSON6902) {
	return cfg.KubeadmConfigPatches, cfg.KubeadmConfigPatchesJSON6902
}

// setPatchNames sets the targeted object name on every patch to be the fixed
// name we use when generating config objects (we have one of each type, all of
// which have the same fixed name)
func setPatchNames(patches []string, jsonPatches []config.PatchJSON6902) ([]string, []config.PatchJSON6902) {
	fixedPatches := make([]string, len(patches))
	fixedJSONPatches := make([]config.PatchJSON6902, len(jsonPatches))
	for i, patch := range patches {
		// insert the generated name metadata
		fixedPatches[i] = fmt.Sprintf("metadata:\nname: %s\n%s", kubeadm.ObjectName, patch)
	}
	for i, patch := range jsonPatches {
		// insert the generated name metadata
		patch.Name = kubeadm.ObjectName
		fixedJSONPatches[i] = patch
	}
	return fixedPatches, fixedJSONPatches
}

// writeKubeadmConfig writes the kubeadm configuration in the specified node
func writeKubeadmConfig(cfg *config.Cluster, data kubeadm.ConfigData) error {
	// get the node ip address
	//nodeAddress, nodeAddressIPv6, err := node.IP()
	//if err != nil {
	//return errors.Wrap(err, "failed to get IP for node")
	//}

	// configure the right protocol addresses
	//if cfg.Networking.IPFamily == "ipv6" {
	//data.NodeAddress = nodeAddressIPv6
	//}

	kubeadmConfig, err := getKubeadmConfig(cfg, data)
	if err != nil {
		return errors.Wrap(err, "failed to generate kubeadm config content")
	}

	glog.Infof("Using kubeadm config:\n" + kubeadmConfig)

	// copy the config to the node
	if err := nodes.WriteFile("/tmp/kubeadm.conf", kubeadmConfig); err != nil {
		// TODO(bentheelder): logging here
		return errors.Wrap(err, "failed to copy kubeadm config to node")
	}

	return nil
}
