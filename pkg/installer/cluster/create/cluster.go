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

package create

import (
	"errors"
	"time"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/nodes"
	internalencoding "github.com/gosoon/kubernetes-operator/pkg/internal/apis/config/encoding"
	internaltypes "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/types"
)

// ClusterOption is a cluster creation option
type ClusterOption func(*internaltypes.ClusterOptions) (*internaltypes.ClusterOptions, error)

// WithNodeImages is kubernetes version image,image contain kubelet,kubectl,kubeadm binary
func WithNodeImage(imageName string, registry string, cluster *ecsv1.KubernetesCluster) ClusterOption {
	return func(o *internaltypes.ClusterOptions) (*internaltypes.ClusterOptions, error) {
		o.NodeImage = nodes.ConstructImage(registry, cluster)
		if imageName != "" {
			o.NodeImage = imageName
		}
		return o, nil
	}
}

// WithNodeAddressAndRole configures the cluster local IP and role from kubernetescluster config
func WithNodeAddressAndRole(cluster *ecsv1.KubernetesCluster) ClusterOption {
	return func(o *internaltypes.ClusterOptions) (*internaltypes.ClusterOptions, error) {
		nodeAddress, role := nodes.ConfigNodeAddressAndRole(cluster)
		if nodeAddress == "" || role == "" {
			return o, errors.New("get nodeAddress and role failed")
		}
		o.NodeAddress, o.Role = nodeAddress, role
		return o, nil
	}
}

// WithConfig configures creating the cluster using the config file at path
func WithConfig(cluster *ecsv1.KubernetesCluster) ClusterOption {
	return func(o *internaltypes.ClusterOptions) (*internaltypes.ClusterOptions, error) {
		o.Config = internalencoding.Ecsv1ToInternalCluster(cluster, o.NodeAddress)
		return o, nil
	}
}

// WithExternalLoadBalancer configures creating the cluster externalLoadBalancer
// if exist externalLoadBalancer that is a vip or select a controlPlane nodeAddress by install server
func WithExternalLoadBalancer(cluster *ecsv1.KubernetesCluster) ClusterOption {
	return func(o *internaltypes.ClusterOptions) (*internaltypes.ClusterOptions, error) {
		o.ExternalLoadBalancer = cluster.Spec.Cluster.ExternalLoadBalancer
		return o, nil
	}
}

// WithKubeConfigPath configures creating the cluster kube config path,default ~/.kube/config-default
func WithKubeConfigPath() ClusterOption {
	return func(o *internaltypes.ClusterOptions) (*internaltypes.ClusterOptions, error) {
		o.KubeConfigPath = nodes.KubeConfigPath(constants.DefaultClusterName)
		return o, nil
	}
}

// WaitForReady configures create to use interval as maximum wait time for the control plane node to be ready
func WaitForReady(interval time.Duration) ClusterOption {
	return func(o *internaltypes.ClusterOptions) (*internaltypes.ClusterOptions, error) {
		o.WaitForReady = interval
		return o, nil
	}
}
