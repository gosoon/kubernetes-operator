package service

import (
	"encoding/json"

	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/gosoon/kubernetes-operator/pkg/types"
	"golang.org/x/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *service) CreateClusterCallback(region string, namespace string, name string, result *types.CallBack) error {
	// TODO: operation failed callback and check receive error
	clientset := s.opt.KubernetesClusterClientset
	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	// if job failed,get the detail log from job's pod
	if !result.Success {
		// update operation annotations
		if _, existed := kubernetesCluster.Annotations[enum.Operation]; existed {
			kubernetesCluster.Annotations[enum.Operation] = enum.KubeCreateFailed
		}
		_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
		if err != nil {
			glog.Errorf("create kubernetesCluster %v/%v failed with:%v", namespace, name, err)
			return err
		}
		return nil
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
		glog.Errorf("create callback update kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	return nil
}

func (s *service) ScaleUpCallback(region string, namespace string, name string, result *types.CallBack) error {
	// TODO: operation failed callback and check receive error
	clientset := s.opt.KubernetesClusterClientset

	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	// if job failed,get the detail log from job's pod
	if !result.Success {
		// update operation annotations
		if _, existed := kubernetesCluster.Annotations[enum.Operation]; existed {
			kubernetesCluster.Annotations[enum.Operation] = enum.KubeScaleUpFinished
		}
		_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
		if err != nil {
			glog.Errorf("scale up kubernetesCluster %v/%v failed with:%v", namespace, name, err)
			return err
		}
		return nil
	}

	// update operation annotations
	if _, existed := kubernetesCluster.Annotations[enum.Operation]; existed {
		kubernetesCluster.Annotations[enum.Operation] = enum.KubeScaleUpFinished
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
		glog.Errorf("scale up and update kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	return nil
}

func (s *service) ScaleDownCallback(region string, namespace string, name string, result *types.CallBack) error {
	clientset := s.opt.KubernetesClusterClientset

	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	// if job failed,get the detail log from job's pod
	if !result.Success {
		// update operation annotations
		if _, existed := kubernetesCluster.Annotations[enum.Operation]; existed {
			kubernetesCluster.Annotations[enum.Operation] = enum.KubeScaleDownFinished
		}
		_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
		if err != nil {
			glog.Errorf("scale down kubernetesCluster %v/%v failed with:%v", namespace, name, err)
			return err
		}
		return nil
	}

	// update operation annotations
	if _, existed := kubernetesCluster.Annotations[enum.Operation]; existed {
		kubernetesCluster.Annotations[enum.Operation] = enum.KubeScaleDownFinished
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
		glog.Errorf("scale down and update kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	return nil
}

func (s *service) DeleteClusterCallback(region string, namespace string, name string, result *types.CallBack) error {
	// TODO: check have running task
	clientset := s.opt.KubernetesClusterClientset
	kubernetesCluster, err := clientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("get kubernetesCluster %v/%v failed with:%v", namespace, name, err)
		return err
	}

	// if job failed,get the detail log from job's pod
	if !result.Success {
		// update operation annotations
		if _, existed := kubernetesCluster.Annotations[enum.Operation]; existed {
			kubernetesCluster.Annotations[enum.Operation] = enum.KubeTerminateFailed
		}
		_, err = clientset.EcsV1().KubernetesClusters(namespace).Update(kubernetesCluster)
		if err != nil {
			glog.Errorf("create kubernetesCluster %v/%v failed with:%v", namespace, name, err)
			return err
		}
		return nil
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
