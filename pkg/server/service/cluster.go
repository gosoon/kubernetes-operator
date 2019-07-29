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
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/gosoon/glog"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *service) CreateCluster(region string, namespace string, name string, clusterInfo *types.EcsClient) error {
	clientset := s.opt.KubernetesClusterClientset
	if !clusterInfo.Retry {
		kubernetesCluster := &ecsv1.KubernetesCluster{
			TypeMeta: metav1.TypeMeta{
				Kind:       "KubernetesCluster",
				APIVersion: "ecs.yun.com/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterInfo.Name,
				Namespace: clusterInfo.Namespace,
			},
			Spec: ecsv1.KubernetesClusterSpec{
				TimeoutMins:   clusterInfo.TimeoutMins,
				ClusterType:   clusterInfo.ClusterType,
				ContainerCIDR: clusterInfo.ContainerCIDR,
				ServiceCIDR:   clusterInfo.ServiceCIDR,
				MasterList:    clusterInfo.MasterList,
				NodeList:      clusterInfo.NodeList,
				EtcdList:      clusterInfo.EtcdList,
				Region:        region,
				AuthConfig: ecsv1.AuthConfig{
					PrivateSSHKey: clusterInfo.PrivateSSHKey,
				},
			},
		}

		// set annotations
		kubernetesCluster.Annotations = map[string]string{
			enum.Operation: enum.KubeCreating,
			enum.Spec:      "",
		}

		_, err := clientset.EcsV1().KubernetesClusters(namespace).Create(kubernetesCluster)
		if err != nil {
			glog.Errorf("create %s/%s cluster failed with:%v", namespace, name, err)
			return err
		}
	} else {
		// create retry
		kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			glog.Errorf("get %s/%s cluster failed with:%v", namespace, name, err)
			return err
		}

		if _, existed := kubernetesCluster.Annotations[enum.Operation]; existed {
			kubernetesCluster.Annotations[enum.Operation] = enum.KubeCreating
		}

		_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
		if err != nil {
			glog.Errorf("update %s/%s operation to KubeCreating failed with:%v", namespace, name, err)
			return err
		}
	}
	return nil
}

func (s *service) DeleteCluster(region string, namespace string, name string, clusterInfo *types.EcsClient) error {
	clientset := s.opt.KubernetesClusterClientset
	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get %s/%s cluster failed with:%v", namespace, name, err)
		return err
	}

	// only in running,failed,finished can delete cluster
	admit, _ := validOperate(kubernetesCluster)
	if !admit {
		return errors.Errorf("current operation is [%v],only in running,failed,finished can delete cluster",
			kubernetesCluster.Annotations[enum.Operation])
	}

	// update operation annotations
	kubernetesCluster.Annotations[enum.Operation] = enum.KubeTerminating

	_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
	if err != nil {
		glog.Errorf("update %s/%s operation to terminating failed with:%v", namespace, name, err)
		return err
	}

	if !clusterInfo.Retry {
		// set DeletePropagation to Foreground,apiserver first set crd DeletionTimestamp field
		// ref: https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/
		p := metav1.DeletePropagationForeground
		err = clientset.EcsV1().KubernetesClusters(namespace).Delete(name, &metav1.DeleteOptions{PropagationPolicy: &p})
		if err != nil {
			glog.Errorf("update %s/%s DeletePropagation failed with:%v", namespace, name, err)
			return err
		}
	}
	return nil
}
