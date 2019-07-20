package controller

import (
	"encoding/json"
	"fmt"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"

	"github.com/gosoon/glog"
	corev1 "k8s.io/api/core/v1"
)

func (c *Controller) processClusterNew(cluster *ecsv1.KubernetesCluster) error {
	curCluster := cluster.DeepCopy()

	namespace := curCluster.Namespace
	name := curCluster.Name
	// create kubernetes cluster
	createClusterJob := newCreateKubernetesClusterJob(curCluster)
	_, err := c.kubeclientset.BatchV1().Jobs(namespace).Create(createClusterJob)
	if err != nil {
		glog.Errorf("create %s/%s kubernetes cluster job failed with:%v", namespace, name, err)
		c.recorder.Event(curCluster, corev1.EventTypeWarning, enum.CreateKubeJobFailed, err.Error())
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.CreateKubeJobSuccess, "")

	// set finalizers
	curCluster.Finalizers = []string{fmt.Sprintf("kubernetescluster.ecs.yun.com/%v", curCluster.Name)}
	curCluster, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).Update(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s spec failed with:%v", err, namespace, name)
		c.recorder.Event(curCluster, corev1.EventTypeWarning, enum.SetFinalizersFailed, err.Error())
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.SetFinalizersSuccess, "")

	// configmap is record the crd operation
	configMap := newConfigMap(curCluster, createClusterJob.Name)
	_, err = c.kubeclientset.CoreV1().ConfigMaps(namespace).Create(configMap)
	if err != nil {
		glog.Errorf("create %s/%s kubeCreating configMap failed with:%v", namespace, name, err)
		return err
	}

	// update phase
	curCluster = curCluster.DeepCopy()
	curCluster.Status.Phase = enum.Creating
	curCluster.Status.JobName = createClusterJob.Name
	_, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update status failed with:%v", err)
		return err
	}
	return nil
}

func (c *Controller) processClusterScaleUp(cluster *ecsv1.KubernetesCluster) error {
	// when status scaling-up + scaling, direct return
	if cluster.Status.Phase == enum.Scaling {
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
	curCluster.Status.JobName = scaleUpClusterJob.Name
	_, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s status to ScalingUp failed with:%v", namespace, name, err)
		return err
	}
	return nil
}

func (c *Controller) processClusterScaleDown(cluster *ecsv1.KubernetesCluster) error {
	// when status scaling-down + scaling, direct return
	if cluster.Status.Phase == enum.Scaling {
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
	curCluster.Status.JobName = scaleDownClusterJob.Name
	_, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s status to ScalingUp failed with:%v", namespace, name, err)
		return err
	}
	return nil
}

func (c *Controller) processClusterTerminating(cluster *ecsv1.KubernetesCluster) error {
	curCluster := cluster.DeepCopy()

	namespace := curCluster.Namespace
	name := curCluster.Name

	// create job
	deleteClusterJob := newDeleteKubernetesClusterJob(curCluster)
	_, err := c.kubeclientset.BatchV1().Jobs(namespace).Create(deleteClusterJob)
	if err != nil {
		glog.Errorf("create delete %s/%s kubernetes cluster job failed with:%v", namespace, name, err)
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
	curCluster.Status.JobName = deleteClusterJob.Name
	_, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s status failed with:%v", namespace, name, err)
		return err
	}
	return nil
}
