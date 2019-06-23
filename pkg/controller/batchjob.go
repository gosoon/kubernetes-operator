package controller

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/utils/pointer"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Image                 string = "ecs-create-k8s"
	RestartPolicy         string = "Never"
	ActiveDeadlineSeconds int32  = 10 * 60
)

func newCreateKubernetesClusterBatchJob(cluster *ecsv1.KubernetesCluster) *batchv1.Job {
	completions := pointer.Int32Ptr(1)
	parallelism := pointer.Int32Ptr(1)
	backoffLimit := pointer.Int32Ptr(0)
	// 10 minutes
	ActiveDeadlineSeconds := pointer.Int64Ptr(10 * 60)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
		},
		Spec: batchv1.JobSpec{
			Parallelism:           parallelism,
			Completions:           completions,
			BackoffLimit:          backoffLimit,
			ActiveDeadlineSeconds: ActiveDeadlineSeconds,
			//  TTLSecondsAfterFinished :  ,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "demo",
							Image: Image,
						},
					},
				},
			},
		},
	}
	return job
}
