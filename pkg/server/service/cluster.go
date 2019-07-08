package service

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"golang.org/x/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *service) CreateCluster(namespace string, name string, kubernetesCluster *ecsv1.KubernetesCluster) error {
	clientset := s.opt.KubernetesClusterClientset

	_, err := clientset.EcsV1().KubernetesClusters(namespace).Create(kubernetesCluster)
	if err != nil {
		glog.Errorf("create kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}
	return nil
}

func (s *service) DeleteCluster(namespace string, name string) error {
	clientset := s.opt.KubernetesClusterClientset

	// set DeletePropagation to Foreground,apiserver first set crd DeletionTimestamp field
	// ref: https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/
	p := metav1.DeletePropagationForeground
	err := clientset.EcsV1().KubernetesClusters(namespace).Delete(name, &metav1.DeleteOptions{PropagationPolicy: &p})
	if err != nil {
		glog.Errorf("delete kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}
	return nil
}

func (s *service) DeleteClusterCallback(namespace string, name string) error {
	// TODO: check have running task
	clientset := s.opt.KubernetesClusterClientset

	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	kubernetesCluster.Finalizers = []string{}
	// update finalizers to null
	_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
	if err != nil {
		glog.Errorf("update kubernetesCluster %v/%v finalizers failed with:%v", namespace, name, err)
		return err
	}

	return nil
}
