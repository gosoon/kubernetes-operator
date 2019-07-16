package service

import (
	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/gosoon/kubernetes-operator/pkg/types"
	"golang.org/x/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *service) ScaleUp(namespace string, name string, clusterInfo *types.EcsClient) error {
	// TODO: operation failed callback and check receive error
	clientset := s.opt.KubernetesClusterClientset

	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	if !clusterInfo.Retry {
		// update node list
		kubernetesCluster.Spec.NodeList = clusterInfo.NodeList
	}

	// update operation annotations
	if _, existed := kubernetesCluster.Annotations[enum.Operation]; existed {
		kubernetesCluster.Annotations[enum.Operation] = enum.KubeScalingUp
	}

	_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
	if err != nil {
		glog.Errorf("create callback update kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	return nil
}

func (s *service) ScaleDown(namespace string, name string, clusterInfo *types.EcsClient) error {
	// TODO: operation failed callback and check receive error
	clientset := s.opt.KubernetesClusterClientset

	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	if !clusterInfo.Retry {
		// update node list
		kubernetesCluster.Spec.NodeList = clusterInfo.NodeList
	}

	// update operation annotations
	if _, existed := kubernetesCluster.Annotations[enum.Operation]; existed {
		kubernetesCluster.Annotations[enum.Operation] = enum.KubeScalingDown
	}

	_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
	if err != nil {
		glog.Errorf("create callback update kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	return nil
}
