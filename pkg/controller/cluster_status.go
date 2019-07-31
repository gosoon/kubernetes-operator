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

package controller

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"

	"github.com/gosoon/glog"
)

func (c *Controller) processOperateFinished(cluster *ecsv1.KubernetesCluster) error {
	if cluster.Status.Phase != enum.Running {
		// update status
		curCluster := cluster.DeepCopy()
		curCluster.Status.Phase = enum.Running
		_, err := c.kubernetesClusterClientset.EcsV1().KubernetesClusters(cluster.Namespace).UpdateStatus(curCluster)
		if err != nil {
			glog.Errorf("update finished cluster status failed with:%v", err)
			return err
		}
	}
	return nil
}

func (c *Controller) processOperateFailed(cluster *ecsv1.KubernetesCluster) error {
	if cluster.Status.Phase != enum.Failed {
		// update status
		curCluster := cluster.DeepCopy()
		curCluster.Status.Phase = enum.Failed
		_, err := c.kubernetesClusterClientset.EcsV1().KubernetesClusters(cluster.Namespace).UpdateStatus(curCluster)
		if err != nil {
			glog.Errorf("update finished cluster status failed with:%v", err)
			return err
		}
	}
	return nil
}

func (c *Controller) processKubeCreating(cluster *ecsv1.KubernetesCluster) error {
	// if kubeCreateFailed and retry,the status is KubeCreating
	if cluster.Status.Phase != enum.Creating {
		// update status
		curCluster := cluster.DeepCopy()
		curCluster.Status.Phase = enum.Creating
		_, err := c.kubernetesClusterClientset.EcsV1().KubernetesClusters(cluster.Namespace).UpdateStatus(curCluster)
		if err != nil {
			glog.Errorf("update finished cluster status failed with:%v", err)
			return err
		}
	}
	return nil
}

func (c *Controller) processNewOperate(cluster *ecsv1.KubernetesCluster) error {
	// if the reason filed is not null,indicating that the job failed,the reason have the job create failed,
	// job timeout...
	if cluster.Status.Reason != "" {
		return nil
	}

	// if kubeCreateFailed and retry,the status is new
	if cluster.Status.Phase != enum.New {
		curCluster := cluster.DeepCopy()
		curCluster.Status.Phase = enum.New
		_, err := c.kubernetesClusterClientset.EcsV1().KubernetesClusters(cluster.Namespace).UpdateStatus(curCluster)
		if err != nil {
			glog.Errorf("update finished cluster status failed with:%v", err)
			return err
		}
	}
	return nil
}
