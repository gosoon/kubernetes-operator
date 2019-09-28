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

package installer

// Installer is a interface define all installer implementation operation
type Installer interface {
	// ClusterNew is creating a new cluster
	ClusterNew(cluster *ecsv1.KubernetesCluster) error

	// ClusterScaleUp is scale up a cluster node
	ClusterScaleUp(cluster *ecsv1.KubernetesCluster, scaleUpNodeList []ecsv1.Node) error

	// ClusterScaleDown is scale down a cluster node
	ClusterScaleDown(cluster *ecsv1.KubernetesCluster, scaleDonwNodeList []ecsv1.Node) error

	// ClusterTerminating is delete a cluster
	ClusterTerminating(cluster *ecsv1.KubernetesCluster) error
}
