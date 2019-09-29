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

package constants

// DefaultClusterName is the default cluster Context name
const DefaultClusterName = "default"

const (
	// ExternalLoadBalancerNodeRoleValue identifies a node that hosts an
	// external load balancer for the API server in HA configurations.
	ExternalLoadBalancerNodeRoleValue string = "external-load-balancer"

	// ExternalEtcdNodeRoleValue identifies a node that hosts an external-etcd
	// instance.
	//
	// WARNING: this node type is not yet implemented!
	//
	ExternalEtcdNodeRoleValue string = "external-etcd"

	// InstallPath is write kubeadm config default path
	InstallPath string = "/tmp/install/"
)
