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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Cluster contains local cluster configuration
type Cluster struct {
	// TypeMeta representing the type of the object and its API schema version.
	metav1.TypeMeta

	// Nodes contains the list of nodes defined in the Cluster
	Nodes []Node

	ExternalLoadBalancer string

	/* Advanced fields */

	// Networking contains cluster wide network settings
	Networking Networking

	// KubeadmConfigPatches are applied to the generated kubeadm config as
	// strategic merge patches to `kustomize build` internally
	// https://github.com/kubernetes/community/blob/master/contributors/devel/strategic-merge-patch.md
	// This should be an inline yaml blob-string
	KubeadmConfigPatches []string

	// KubeadmConfigPatchesJSON6902 are applied to the generated kubeadm config
	// as patchesJson6902 to `kustomize build`
	KubeadmConfigPatchesJSON6902 []PatchJSON6902

	// kubernetes version
	KubeVersion string

	// ImagesRegistry contains all images during install cluster
	ImagesRegistry string
}

// Node contains settings for a node in the Cluster.
type Node struct {
	IP string

	// Role defines the role of the node in the in the Kubernetes cluster
	Role ecsv1.NodeRole
}

// Networking contains cluster wide network settings
type Networking struct {
	// IPFamily is the network cluster model, currently it can be ipv4 or ipv6
	IPFamily ClusterIPFamily
	// APIServerPort is the listen port on the host for the Kubernetes API Server
	// Defaults to a random port on the host
	APIServerPort int32
	// APIServerAddress is the listen address on the host for the Kubernetes
	// API Server. This should be an IP address.
	//
	// Defaults to 0.0.0.0
	APIServerAddress string
	// PodSubnet is the CIDR used for pod IPs
	PodSubnet string
	// ServiceSubnet is the CIDR used for services VIPs
	ServiceSubnet string
	// If DisableDefaultCNI is true  will not install the default CNI setup.
	// Instead the user should install their own CNI after creating the cluster.
	DisableDefaultCNI bool
}

// ClusterIPFamily defines cluster network IP family
type ClusterIPFamily string

const (
	// IPv4Family sets ClusterIPFamily to ipv4
	IPv4Family ClusterIPFamily = "ipv4"
	// IPv6Family sets ClusterIPFamily to ipv6
	IPv6Family ClusterIPFamily = "ipv6"
)

// PatchJSON6902 represents an inline kustomize json 6902 patch
// https://tools.ietf.org/html/rfc6902
type PatchJSON6902 struct {
	// these fields specify the patch target resource
	Group   string
	Version string
	Kind    string
	// Name and Namespace are optional
	// NOTE: technically name is required now, but we default it elsewhere
	// Third party users of this type / library would need to set it.
	Name      string
	Namespace string
	// Patch should contain the contents of the json patch as a string
	Patch string
}
