package controller

import (
	"github.com/gosoon/glog"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
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
	// if kubeCreateFailed and retry,the status is new
	curCluster := cluster.DeepCopy()
	curCluster.Status.Phase = enum.New
	_, err := c.kubernetesClusterClientset.EcsV1().KubernetesClusters(cluster.Namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update finished cluster status failed with:%v", err)
		return err
	}
	return nil
}
