package controller

import (
	"fmt"

	"github.com/gosoon/glog"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
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
		c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.CreateKubeJobFailed, err.Error())
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.CreateKubeJobSuccess, "")

	// set finalizers
	curCluster.Finalizers = []string{fmt.Sprintf("kubernetescluster.ecs.yun.com/%v", curCluster.Name)}
	k, err := c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).Update(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s spec failed with:%v", err, namespace, name)
		c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.SetFinalizersFailed, err.Error())
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.SetFinalizersSuccess, "")

	// update phase
	curCluster = k.DeepCopy()
	curCluster.Status.Phase = enum.Creating
	k, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update status failed with:%v", err)
		return err
	}
	return nil
}

func (c *Controller) processClusterScaleUp(cluster *ecsv1.KubernetesCluster) error {
	curCluster := cluster.DeepCopy()

	namespace := curCluster.Namespace
	name := curCluster.Name
	scaleUpClusterJob := newScaleUpClusterJob(namespace, name)
	_, err := c.kubeclientset.BatchV1().Jobs(namespace).Create(scaleUpClusterJob)
	if err != nil {
		glog.Errorf("create %s/%s scale up cluster job failed with:%v", namespace, name, err)
		c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.CreateScaleUpJobFailed, err.Error())
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.CreateScaleUpJobSuccess, "")

	// update phase to ScalingUp
	curCluster.Status.Phase = enum.Scaling
	_, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s status to ScalingUp failed with:%v", namespace, name, err)
		return err
	}
	return nil
}

func (c *Controller) processClusterScaleDown(cluster *ecsv1.KubernetesCluster) error {
	curCluster := cluster.DeepCopy()

	namespace := curCluster.Namespace
	name := curCluster.Name
	scaleDownClusterJob := newScaleDownClusterJob(namespace, name)
	_, err := c.kubeclientset.BatchV1().Jobs(namespace).Create(scaleDownClusterJob)
	if err != nil {
		glog.Errorf("create %s/%s scale up cluster job failed with:%v", namespace, name, err)
		c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.CreateScaleDownJobFailed, err.Error())
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.CreateScaleDownJobSuccess, "")

	// update phase to ScalingDown
	curCluster.Status.Phase = enum.Scaling
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
	deleteClusterJob := newDeleteKubernetesClusterJob(namespace, name)
	_, err := c.kubeclientset.BatchV1().Jobs(namespace).Create(deleteClusterJob)
	if err != nil {
		glog.Errorf("create delete %s/%s kubernetes cluster job failed with:%v", namespace, name, err)
		c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.DeleteKubeJobFailed, "")
		return err
	}
	c.recorder.Event(curCluster, corev1.EventTypeNormal, enum.DeleteKubeJobSuccess, "")

	// update label if delete job created
	deleteLabel := deleteClusterJob.Name
	if curCluster.Labels == nil {
		curCluster.Labels = map[string]string{}
	}
	if _, existed := curCluster.Labels[deleteLabel]; !existed {
		curCluster.Labels[deleteLabel] = DeleteJobLabelCreated
	}
	k, err := c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).Update(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s spec failed with:%v", namespace, name, err)
		return err
	}

	// update status to Terminating
	curCluster = k.DeepCopy()
	curCluster.Status.Phase = enum.Terminating
	k, err = c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update %s/%s status failed with:%v", namespace, name, err)
		return err
	}
	return nil
}

func (c *Controller) processClusterNotExistInCache(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	glog.Warningf("KubernetesCluster: %s/%s does not exist in local cache, will delete it from Neutron ...",
		namespace, name)

	glog.Infof("[Neutron] Deleting kubernetesCluster: %s/%s ...", namespace, name)

	// neutron.Delete(namespace, name)
	// delete job and kubernetes cluster
	// update DeletionTimestamp
	deleteClusterJob := newDeleteKubernetesClusterJob(namespace, name)
	_, err = c.kubeclientset.BatchV1().Jobs(namespace).Create(deleteClusterJob)
	if err != nil {
		glog.Errorf("create delete %s/%s kubernetes cluster job failed with:%v", namespace, name, err)
		return err
	}

	// call back and delete deleteJob and createJob
	return nil
}
