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
	"encoding/json"
	"fmt"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"

	"github.com/gosoon/glog"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Controller) processClusterNew(cluster *ecsv1.KubernetesCluster) error {
	curCluster := cluster.DeepCopy()

	namespace := curCluster.Namespace
	name := curCluster.Name
	// create kubernetes cluster job
	createClusterJob := newCreateKubernetesClusterJob(curCluster)
	_, err := c.kubeclientset.BatchV1().Jobs(namespace).Create(createClusterJob)
	if err != nil {
		glog.Errorf("create %s/%s kubernetes cluster job failed with:%v", namespace, name, err)
		c.recorder.Event(curCluster, corev1.EventTypeWarning, enum.CreateKubeJobFailed, err.Error())
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.CreateKubeJobSuccess, "")

	// configmap is record the crd operation
	configMap := newConfigMap(curCluster, createClusterJob.Name)
	_, err = c.kubeclientset.CoreV1().ConfigMaps(namespace).Create(configMap)
	if err != nil {
		glog.Errorf("create %s/%s kubeCreating configMap failed with:%v", namespace, name, err)
		return err
	}

	// update phase
	curCluster.Status.Phase = enum.Creating
	curCluster.Status.LastTransitionTime = metav1.Now()
	curCluster.Status.JobName = createClusterJob.Name
	curCluster, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update status failed with:%v", err)
		return err
	}

	// set finalizers
	curCluster = curCluster.DeepCopy()
	curCluster.Finalizers = []string{fmt.Sprintf("kubernetescluster.ecs.yun.com/%v", curCluster.Name)}
	curCluster, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).Update(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s spec failed with:%v", err, namespace, name)
		c.recorder.Event(curCluster, corev1.EventTypeWarning, enum.SetFinalizersFailed, err.Error())
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.SetFinalizersSuccess, "")

	go c.jobTTLControl(curCluster)
	return nil
}

func (c *Controller) processClusterScaleUp(cluster *ecsv1.KubernetesCluster) error {
	// if the reason filed is not null,indicating that the job failed,the reason have the job create failed,
	//job timeout...
	if cluster.Status.Phase == enum.Scaling && cluster.Status.Reason != "" {
		return nil
	}

	curCluster := cluster.DeepCopy()
	namespace := curCluster.Namespace
	name := curCluster.Name

	// diff work nodes
	var oldSpec ecsv1.KubernetesClusterSpec
	oldSpecStr := curCluster.Annotations[enum.Spec]
	err := json.Unmarshal([]byte(oldSpecStr), &oldSpec)
	if err != nil {
		glog.Errorf("get old spec failed with:%v", err)
		return err
	}
	nodeList := diffNodeList(oldSpec.NodeList, cluster.Spec.NodeList, cluster.Annotations[enum.Operation])

	// create job
	scaleUpClusterJob := newScaleUpClusterJob(curCluster, nodeList)
	_, err = c.kubeclientset.BatchV1().Jobs(namespace).Create(scaleUpClusterJob)
	if err != nil {
		glog.Errorf("create %s/%s scale up cluster job failed with:%v", namespace, name, err)
		c.recorder.Event(curCluster, corev1.EventTypeWarning, enum.CreateScaleUpJobFailed, err.Error())
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.CreateScaleUpJobSuccess, "")

	// configmap is record the crd operation
	configMap := newConfigMap(curCluster, scaleUpClusterJob.Name)
	_, err = c.kubeclientset.CoreV1().ConfigMaps(namespace).Create(configMap)
	if err != nil {
		glog.Errorf("create %s/%s scaleUp configMap failed with:%v", namespace, name, err)
		return err
	}

	// update phase to ScalingUp
	curCluster.Status.Phase = enum.Scaling
	curCluster.Status.LastTransitionTime = metav1.Now()
	curCluster.Status.Reason = ""
	curCluster.Status.JobName = scaleUpClusterJob.Name
	curCluster, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s status to ScalingUp failed with:%v", namespace, name, err)
		return err
	}

	go c.jobTTLControl(curCluster)
	return nil
}

func (c *Controller) processClusterScaleDown(cluster *ecsv1.KubernetesCluster) error {
	//if the reason filed is not null,indicating that the job failed,the reason have the job create failed,
	//job timeout...
	if cluster.Status.Phase == enum.Scaling && cluster.Status.Reason != "" {
		return nil
	}

	curCluster := cluster.DeepCopy()
	namespace := curCluster.Namespace
	name := curCluster.Name

	// diff work nodes
	var oldSpec ecsv1.KubernetesClusterSpec
	oldSpecStr := curCluster.Annotations[enum.Spec]
	err := json.Unmarshal([]byte(oldSpecStr), &oldSpec)
	if err != nil {
		glog.Errorf("get old spec failed with:%v", err)
		return err
	}
	nodeList := diffNodeList(oldSpec.NodeList, cluster.Spec.NodeList, cluster.Annotations[enum.Operation])

	// create job
	scaleDownClusterJob := newScaleDownClusterJob(curCluster, nodeList)
	_, err = c.kubeclientset.BatchV1().Jobs(namespace).Create(scaleDownClusterJob)
	if err != nil {
		glog.Errorf("create %s/%s scale up cluster job failed with:%v", namespace, name, err)
		c.recorder.Event(curCluster, corev1.EventTypeWarning, enum.CreateScaleDownJobFailed, err.Error())
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.CreateScaleDownJobSuccess, "")

	// configmap is record the crd operation
	configMap := newConfigMap(curCluster, scaleDownClusterJob.Name)
	_, err = c.kubeclientset.CoreV1().ConfigMaps(namespace).Create(configMap)
	if err != nil {
		glog.Errorf("create %s/%s scaleDown configMap failed with:%v", namespace, name, err)
		return err
	}

	// update phase to ScalingDown
	curCluster.Status.Phase = enum.Scaling
	curCluster.Status.LastTransitionTime = metav1.Now()
	curCluster.Status.Reason = ""
	curCluster.Status.JobName = scaleDownClusterJob.Name
	curCluster, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s status to ScalingUp failed with:%v", namespace, name, err)
		return err
	}

	go c.jobTTLControl(curCluster)
	return nil
}

func (c *Controller) processClusterTerminating(cluster *ecsv1.KubernetesCluster) error {
	// if the reason filed is not null,indicating that the last terminating job failed,the reason have the job create failed,
	// job timeout...
	if cluster.Status.Phase == enum.Terminating && cluster.Status.Reason != "" {
		return nil
	}

	curCluster := cluster.DeepCopy()
	namespace := curCluster.Namespace
	name := curCluster.Name

	// create job
	deleteClusterJob := newDeleteKubernetesClusterJob(curCluster)
	_, err := c.kubeclientset.BatchV1().Jobs(namespace).Create(deleteClusterJob)
	if err != nil {
		glog.Errorf("create delete-%s/%s-kubernetes cluster job failed with:%v", namespace, name, err)
		c.recorder.Event(curCluster, corev1.EventTypeWarning, enum.DeleteKubeJobFailed, "")
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.DeleteKubeJobSuccess, "")

	// configmap is record the crd operation
	configMap := newConfigMap(curCluster, deleteClusterJob.Name)
	_, err = c.kubeclientset.CoreV1().ConfigMaps(namespace).Create(configMap)
	if err != nil {
		glog.Errorf("create %s/%s delete configMap failed with:%v", namespace, name, err)
		return err
	}

	// update status to Terminating
	curCluster.Status.Phase = enum.Terminating
	curCluster.Status.LastTransitionTime = metav1.Now()
	curCluster.Status.Reason = ""
	curCluster.Status.JobName = deleteClusterJob.Name
	curCluster, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s status failed with:%v", namespace, name, err)
		return err
	}

	go c.jobTTLControl(curCluster)
	return nil
}
