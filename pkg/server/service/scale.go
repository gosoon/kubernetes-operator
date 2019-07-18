package service

import (
	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/gosoon/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *service) ScaleUp(region string, namespace string, name string, clusterInfo *types.EcsClient) error {
	clientset := s.opt.KubernetesClusterClientset

	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get %s/%s cluster failed with:%v", namespace, name, err)
		return err
	}

	// if latest task must be finished and start next task
	admit, err := validOperate(kubernetesCluster)
	if !admit {
		return err
	}

	if !clusterInfo.Retry {
		// update node list
		kubernetesCluster.Spec.NodeList = clusterInfo.NodeList
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

func (s *service) ScaleDown(region string, namespace string, name string, clusterInfo *types.EcsClient) error {
	clientset := s.opt.KubernetesClusterClientset

	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get %s/%s cluster failed with:%v", namespace, name, err)
		return err
	}

	// if latest task must be finished and start next task
	admit, err := validOperate(kubernetesCluster)
	if !admit {
		return err
	}

	if !clusterInfo.Retry {
		// update node list
		kubernetesCluster.Spec.NodeList = clusterInfo.NodeList
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
