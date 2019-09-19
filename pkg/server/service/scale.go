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

package service

import (
	"context"

	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/gosoon/kubernetes-operator/pkg/types"
	"github.com/pkg/errors"

	"github.com/gosoon/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *service) ScaleUp(ctx context.Context, region string, namespace string, name string, clusterInfo *types.EcsClient) error {
	clientset := s.opt.KubernetesClusterClientset

	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get %s/%s cluster failed with:%v", namespace, name, err)
		return err
	}

	// if latest task must be finished and start next task
	admit := validPhase(kubernetesCluster)
	if !admit {
		return errors.Errorf("current operation is [%v],only in running,failed can be start next task",
			kubernetesCluster.Status.Phase)
	}

	if !clusterInfo.Retry {
		// update node list
		kubernetesCluster.Spec.Cluster.NodeList = clusterInfo.NodeList
	}

	// update operation annotations
	kubernetesCluster.Annotations[enum.Operation] = enum.KubeScalingUp

	_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
	if err != nil {
		glog.Errorf("update %s/%s operation to KubeScalingUp failed with:%v", namespace, name, err)
		return err
	}

	return nil
}

func (s *service) ScaleDown(ctx context.Context, region string, namespace string, name string, clusterInfo *types.EcsClient) error {
	clientset := s.opt.KubernetesClusterClientset

	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get %s/%s cluster failed with:%v", namespace, name, err)
		return err
	}

	// if latest task must be finished and start next task
	admit := validPhase(kubernetesCluster)
	if !admit {
		return errors.Errorf("current operation is [%v],only in running,failed can be start next task",
			kubernetesCluster.Status.Phase)
	}

	if !clusterInfo.Retry {
		// update node list
		kubernetesCluster.Spec.Cluster.NodeList = clusterInfo.NodeList
	}

	// update operation annotations
	kubernetesCluster.Annotations[enum.Operation] = enum.KubeScalingDown

	_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
	if err != nil {
		glog.Errorf("update %s/%s operation to KubeScalingDown failed with:%v", namespace, name, err)
		return err
	}

	return nil
}
