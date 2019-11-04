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

import (
	"context"

	"github.com/gosoon/glog"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/installer/util/protobuf"
)

// ClusterNew is creating a new cluster
func (inst *installer) ClusterNew(cluster *ecsv1.KubernetesCluster) error {
	clusterRequest, err := protobuf.ClusterConvertToProtobuf(cluster)
	if err != nil {
		glog.Errorf("clusterRequest convert to protobuf failed with:%v", err)
		return err
	}

	_, err = inst.InstallCluster(context.Background(), clusterRequest)
	if err != nil {
		glog.Errorf("installCluster failed with %v", err)
		return err
	}

	return nil
}

// ClusterScaleUp is scale up a cluster node
func (inst *installer) ClusterScaleUp(cluster *ecsv1.KubernetesCluster, scaleUpNodeList []ecsv1.Node) error {
	// TODO
	return nil
}

// ClusterScaleDown is scale down a cluster node
func (inst *installer) ClusterScaleDown(cluster *ecsv1.KubernetesCluster, scaleDownNodeList []ecsv1.Node) error {
	// TODO
	return nil
}

// ClusterTerminating is delete a cluster
func (inst *installer) ClusterTerminating(cluster *ecsv1.KubernetesCluster) error {
	// TODO
	return nil
}
