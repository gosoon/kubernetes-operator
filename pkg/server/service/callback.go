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
	"encoding/json"
	"fmt"

	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/gosoon/kubernetes-operator/pkg/kuberesource"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/gosoon/glog"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *service) CreateClusterCallback(ctx context.Context, region string, namespace string, name string, result *types.Callback) error {
	clientset := s.opt.KubernetesClusterClientset
	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get %s/%s cluster failed with:%v", namespace, name, err)
		return err
	}

	// if already callback and return
	if kubernetesCluster.Annotations[enum.Operation] != enum.KubeCreating {
		return errors.New("callback is already done.Current operation not is kubeCreating")
	}

	// if job failed,get the detail log from job's pod
	if !result.Success {
		// update operation annotations
		if _, existed := kubernetesCluster.Annotations[enum.Operation]; existed {
			kubernetesCluster.Annotations[enum.Operation] = enum.KubeCreateFailed
		}
		_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
		if err != nil {
			glog.Errorf("update %s/%s operation to KubeCreateFailed failed with:%v", namespace, name, err)
			return err
		}
		return nil
	}

	// save kubeconfig to configmap
	configmap := kuberesource.NewConfigMap(kubernetesCluster)
	configmap.Name = fmt.Sprintf("%v-%v-kubeconfig", namespace, name)
	configmap.Data = map[string]string{"kubeconfig": result.KubeConfig}
	_, err = s.opt.KubeClientset.CoreV1().ConfigMaps(namespace).Create(configmap)
	if err != nil {
		glog.Errorf("create %s/%s kubeconfig configMap failed with:%v", namespace, name, err)
		return err
	}

	// update operation annotations
	if _, existed := kubernetesCluster.Annotations[enum.Operation]; existed {
		kubernetesCluster.Annotations[enum.Operation] = enum.KubeCreateFinished
	}

	// update spec annotations
	if _, existed := kubernetesCluster.Annotations[enum.Spec]; existed {
		specBytes, err := json.Marshal(kubernetesCluster.Spec)
		if err != nil {
			glog.Errorf("marshal cluster spec field failed with :%v", err)
			return err
		}
		kubernetesCluster.Annotations[enum.Spec] = string(specBytes)
	}

	_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
	if err != nil {
		glog.Errorf("update %s/%s operation to KubeCreateFinished failed with:%v", namespace, name, err)
		return err
	}

	return nil
}

func (s *service) ScaleUpCallback(ctx context.Context, region string, namespace string, name string, result *types.Callback) error {
	clientset := s.opt.KubernetesClusterClientset

	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get %s/%s cluster failed with:%v", namespace, name, err)
		return err
	}

	// if already callback and return
	if kubernetesCluster.Annotations[enum.Operation] != enum.KubeScalingUp {
		return errors.New("callback is already done.current operation not is kubeScalingUp")
	}

	// if job failed,get the detail log from job's pod log
	if !result.Success {
		// update operation annotations
		kubernetesCluster.Annotations[enum.Operation] = enum.KubeScaleUpFailed

		_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
		if err != nil {
			glog.Errorf("update %s/%s operation to KubeScaleUpFailed failed with:%v", namespace, name, err)
			return err
		}
		return nil
	}

	// update operation annotations
	kubernetesCluster.Annotations[enum.Operation] = enum.KubeScaleUpFinished

	// update spec annotations
	if _, existed := kubernetesCluster.Annotations[enum.Spec]; existed {
		specBytes, err := json.Marshal(kubernetesCluster.Spec)
		if err != nil {
			glog.Errorf("marshal cluster spec field failed with :%v", err)
			return err
		}
		kubernetesCluster.Annotations[enum.Spec] = string(specBytes)
	}

	_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
	if err != nil {
		glog.Errorf("update %s/%s operation to KubeScaleUpFinished failed with:%v", namespace, name, err)
		return err
	}

	return nil
}

func (s *service) ScaleDownCallback(ctx context.Context, region string, namespace string, name string, result *types.Callback) error {
	clientset := s.opt.KubernetesClusterClientset

	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get %s/%s cluster failed with:%v", namespace, name, err)
		return err
	}

	// if already callback and return
	if kubernetesCluster.Annotations[enum.Operation] != enum.KubeScalingDown {
		return errors.New("callback is already done.current operation not is kubeScalingDown")
	}

	// if job failed,get the detail log from job's pod
	if !result.Success {
		// update operation annotations
		kubernetesCluster.Annotations[enum.Operation] = enum.KubeScaleDownFailed

		_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
		if err != nil {
			glog.Errorf("update %s/%s operation to KubeScaleDownFailed failed with:%v", namespace, name, err)
			return err
		}
		return nil
	}

	// update operation annotations
	kubernetesCluster.Annotations[enum.Operation] = enum.KubeScaleDownFinished

	// update spec annotations
	if _, existed := kubernetesCluster.Annotations[enum.Spec]; existed {
		specBytes, err := json.Marshal(kubernetesCluster.Spec)
		if err != nil {
			glog.Errorf("marshal cluster spec field failed with :%v", err)
			return err
		}
		kubernetesCluster.Annotations[enum.Spec] = string(specBytes)
	}

	_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
	if err != nil {
		glog.Errorf("update %s/%s operation to KubeScaleDownFinished failed with:%v", namespace, name, err)
		return err
	}

	return nil
}

func (s *service) DeleteClusterCallback(ctx context.Context, region string, namespace string, name string, result *types.Callback) error {
	clientset := s.opt.KubernetesClusterClientset
	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get %s/%s cluster failed with:%v", namespace, name, err)
		return err
	}

	// if already callback and return
	if kubernetesCluster.Annotations[enum.Operation] != enum.KubeTerminating {
		return errors.New("callback is already done.current operation not is kubeTerminating")
	}

	// if job failed,get the detail log from job's pod
	if !result.Success {
		// update operation annotations
		kubernetesCluster.Annotations[enum.Operation] = enum.KubeTerminateFailed

		_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
		if err != nil {
			glog.Errorf("update %s/%s operation to KubeTerminateFailed failed with:%v", namespace, name, err)
			return err
		}
		return nil
	}

	kubernetesCluster.Finalizers = []string{}
	// update finalizers to null
	_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
	if err != nil {
		glog.Errorf("update %s/%s finalizers to nil failed with:%v", namespace, name, err)
		return err
	}

	return nil
}
