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

package server

import installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"

// injectClusterConfig is set some config by server,eg:image registry and node role
// current only inject images registry
func (inst *installer) injectClusterConfig(cluster *installerv1.KubernetesClusterRequest) *installerv1.KubernetesClusterRequest {
	cluster.Spec.Cluster.ImagesRegistry = inst.opt.ImagesRegistry
	return cluster
}
