/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/nodes"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/kubeadm"

	"github.com/gosoon/glog"
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

	kubeVersion := ctx.Cluster.Config.KubeVersion

	// create kubeadm init config
	configData := kubeadm.ConfigData{
		ClusterName:          ctx.Cluster.Name, // to do
		KubernetesVersion:    kubeVersion,
		ControlPlaneEndpoint: ctx.Cluster.ExternalLoadBalancer, // is external load balancer ???
		APIBindPort:          kubeadm.APIServerPort,
		APIServerAddress:     ctx.Cluster.Config.Networking.APIServerAddress,
		Token:                kubeadm.Token,
		PodSubnet:            ctx.Cluster.Config.Networking.PodSubnet,
		ServiceSubnet:        ctx.Cluster.Config.Networking.ServiceSubnet,
		ControlPlane:         true,
		//IPv6:                 ctx.Config.Networking.IPFamily == "ipv6",
		NodeAddress: ctx.Cluster.NodeAddress,
	}

	if ctx.Cluster.Role == ecsv1.WorkerRole {
		configData.ControlPlane = false
	}

	if err := writeKubeadmConfig(configData); err != nil {
		return err
	}

	// mark success
	ctx.Status.End(true)

	return nil
}

// generateKubeadmConfig generates the kubeadm config contents for the cluster
// by running data through the template.
func generateKubeadmConfig(data kubeadm.ConfigData) (path string, err error) {
	// generate the config contents
	config, err := kubeadm.Config(data)
	if err != nil {
		return "", err
	}
	return config, nil
}

// writeKubeadmConfig writes the kubeadm configuration in the specified node
func writeKubeadmConfig(data kubeadm.ConfigData) error {
	kubeadmConfig, err := generateKubeadmConfig(data)
	if err != nil {
		return errors.Wrap(err, "failed to generate kubeadm config content")
	}

	glog.Infof("Using kubeadm config:\n" + kubeadmConfig)

	// copy the config to the node
	if err := nodes.WriteFile(constants.InstallPath+"kubeadm.conf", kubeadmConfig); err != nil {
		return errors.Wrap(err, "failed to copy kubeadm config to node")
	}

	return nil
}
