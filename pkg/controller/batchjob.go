package controller

import (
	"fmt"

	"github.com/gosoon/kubernetes-operator/pkg/apis/ecs"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/utils/pointer"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Image                 string = "busybox:latest"
	RestartPolicy         string = "Never"
	ActiveDeadlineSeconds int32  = 10 * 60
)

func newCreateKubernetesClusterBatchJob(cluster *ecsv1.KubernetesCluster) *batchv1.Job {
	jobName := fmt.Sprintf("create-%v-%v-cluster", cluster.Namespace, cluster.Name)
	completions := pointer.Int32Ptr(1)
	parallelism := pointer.Int32Ptr(1)
	backoffLimit := pointer.Int32Ptr(0)
	// 10 minutes
	ActiveDeadlineSeconds := pointer.Int64Ptr(10 * 60)

	fmt.Println("%+v", cluster)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: cluster.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         fmt.Sprintf("%v/v1", ecs.GroupName), // occur invalid error
					Kind:               "KubernetesCluster",                 // occur invalid error
					Name:               cluster.Name,
					UID:                cluster.UID,
					Controller:         pointer.BoolPtr(true),
					BlockOwnerDeletion: pointer.BoolPtr(true),
				},
			},
		},
		Spec: batchv1.JobSpec{
			Parallelism:           parallelism,
			Completions:           completions,
			BackoffLimit:          backoffLimit,
			ActiveDeadlineSeconds: ActiveDeadlineSeconds,
			//  TTLSecondsAfterFinished :  ,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  jobName,
							Image: Image,
						},
					},
				},
			},
		},
	}
	return job
}

func newDeleteKubernetesClusterBatchJob(name string, namespace string) *batchv1.Job {
	jobName := fmt.Sprintf("delete-%v-%v-cluster", namespace, name)
	completions := pointer.Int32Ptr(1)
	parallelism := pointer.Int32Ptr(1)
	backoffLimit := pointer.Int32Ptr(0)
	// 10 minutes
	ActiveDeadlineSeconds := pointer.Int64Ptr(10 * 60)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:            jobName,
			Namespace:       namespace,
			OwnerReferences: []metav1.OwnerReference{},
		},
		Spec: batchv1.JobSpec{
			Parallelism:           parallelism,
			Completions:           completions,
			BackoffLimit:          backoffLimit,
			ActiveDeadlineSeconds: ActiveDeadlineSeconds,
			//  TTLSecondsAfterFinished :  ,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  jobName,
							Image: Image,
						},
					},
				},
			},
		},
	}
	return job
}
