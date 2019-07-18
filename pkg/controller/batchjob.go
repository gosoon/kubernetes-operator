package controller

import (
	"fmt"
	"time"

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

func newCreateKubernetesClusterJob(cluster *ecsv1.KubernetesCluster) *batchv1.Job {
	jobName := fmt.Sprintf("create-%v-%v-job", cluster.Namespace, cluster.Name)
	completions := pointer.Int32Ptr(1)
	parallelism := pointer.Int32Ptr(1)
	backoffLimit := pointer.Int32Ptr(0)
	// 10 minutes
	ActiveDeadlineSeconds := pointer.Int64Ptr(10 * 60)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: cluster.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         fmt.Sprintf("%v/v1", ecs.GroupName), // not define and occur invalid error
					Kind:               "KubernetesCluster",                 // not define and occur invalid error
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
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  jobName,
							Image: Image,
							Env: []corev1.EnvVar{
								{
									Name:  "CLUSTER_MASTER_LIST",
									Value: convertNodesToString(cluster.Spec.MasterList),
								},
								{
									Name:  "CLUSTER_NODE_LIST",
									Value: convertNodesToString(cluster.Spec.NodeList),
								},
								{
									Name:  "CLUSTER_ETCD_LIST",
									Value: convertNodesToString(cluster.Spec.EtcdList),
								},
							},
						},
					},
				},
			},
		},
	}
	return job
}

func newDeleteKubernetesClusterJob(namespace string, name string) *batchv1.Job {
	jobName := fmt.Sprintf("delete-%v-%v-job", namespace, name)
	completions := pointer.Int32Ptr(1)
	parallelism := pointer.Int32Ptr(1)
	backoffLimit := pointer.Int32Ptr(0)
	// 60 minutes
	ttlSecondsAfterFinished := pointer.Int32Ptr(60 * 60)
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
			// if you want to clean up finished jobs automatically,
			// enabled with feature gate TTLAfterFinished.
			TTLSecondsAfterFinished: ttlSecondsAfterFinished,
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

func newScaleUpClusterJob(cluster *ecsv1.KubernetesCluster) *batchv1.Job {
	// diff work node
	namespace := cluster.Namespace
	name := cluster.Name
	jobName := fmt.Sprintf("scale-up-%v-%v-job-%v", namespace, name, time.Now().Unix())
	completions := pointer.Int32Ptr(1)
	parallelism := pointer.Int32Ptr(1)
	backoffLimit := pointer.Int32Ptr(0)
	// 10 minutes
	ActiveDeadlineSeconds := pointer.Int64Ptr(10 * 60)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         fmt.Sprintf("%v/v1", ecs.GroupName), // not define and occur invalid error
					Kind:               "KubernetesCluster",                 // not define and occur invalid error
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

func newScaleDownClusterJob(cluster *ecsv1.KubernetesCluster) *batchv1.Job {
	// diff work node
	namespace := cluster.Namespace
	name := cluster.Name
	jobName := fmt.Sprintf("scale-down-%v-%v-job-%v", namespace, name, time.Now().Unix())
	completions := pointer.Int32Ptr(1)
	parallelism := pointer.Int32Ptr(1)
	backoffLimit := pointer.Int32Ptr(0)
	// 10 minutes
	ActiveDeadlineSeconds := pointer.Int64Ptr(10 * 60)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         fmt.Sprintf("%v/v1", ecs.GroupName), // not define and occur invalid error
					Kind:               "KubernetesCluster",                 // not define and occur invalid error
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
