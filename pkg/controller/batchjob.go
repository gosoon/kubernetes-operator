/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controller

import (
	"fmt"
	"time"

	"github.com/gosoon/kubernetes-operator/pkg/apis/ecs"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/gosoon/kubernetes-operator/pkg/utils/pointer"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Image                 = string("busybox:latest")
	RestartPolicy         = string("Never")
	ActiveDeadlineSeconds = int32(10 * 60)
	Kind                  = string("KubernetesCluster")
)

func newCreateKubernetesClusterJob(cluster *ecsv1.KubernetesCluster) *batchv1.Job {
	jobName := fmt.Sprintf("create-%v-%v-job", cluster.Namespace, cluster.Name)
	completions := pointer.Int32Ptr(1)
	parallelism := pointer.Int32Ptr(1)
	backoffLimit := pointer.Int32Ptr(0)
	// 10 minutes
	ActiveDeadlineSeconds := pointer.Int64Ptr(10 * 60)

	// pack envs
	envs := compressEnvs(cluster, enum.KubeCreating)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: cluster.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         fmt.Sprintf("%v/v1", ecs.GroupName), // not define and occur invalid error
					Kind:               Kind,                                // not define and occur invalid error
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
							Env:   envs,
						},
					},
				},
			},
		},
	}

	job.Spec.Template.Spec.Containers[0].Env = append(job.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{
			Name:  "NODE_HOSTS",
			Value: convertNodesToString(cluster.Spec.NodeList),
		})

	return job
}

func newDeleteKubernetesClusterJob(cluster *ecsv1.KubernetesCluster) *batchv1.Job {
	jobName := fmt.Sprintf("delete-%v-%v-job", cluster.Namespace, cluster.Name)
	completions := pointer.Int32Ptr(1)
	parallelism := pointer.Int32Ptr(1)
	backoffLimit := pointer.Int32Ptr(0)
	// 60 minutes
	ttlSecondsAfterFinished := pointer.Int32Ptr(60 * 60)
	// 10 minutes
	ActiveDeadlineSeconds := pointer.Int64Ptr(10 * 60)

	// pack envs
	envs := compressEnvs(cluster, enum.KubeTerminating)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:            jobName,
			Namespace:       cluster.Namespace,
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
							Env:   envs,
						},
					},
				},
			},
		},
	}

	job.Spec.Template.Spec.Containers[0].Env = append(job.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{
			Name:  "NODE_HOSTS",
			Value: convertNodesToString(cluster.Spec.NodeList),
		})

	return job
}

func newScaleUpClusterJob(cluster *ecsv1.KubernetesCluster, diffNodeList []ecsv1.Node) *batchv1.Job {
	// diff work node
	namespace := cluster.Namespace
	name := cluster.Name
	jobName := fmt.Sprintf("scale-up-%v-%v-job-%v", namespace, name, time.Now().Unix())
	completions := pointer.Int32Ptr(1)
	parallelism := pointer.Int32Ptr(1)
	backoffLimit := pointer.Int32Ptr(0)
	// 10 minutes
	ActiveDeadlineSeconds := pointer.Int64Ptr(10 * 60)

	// pack envs
	envs := compressEnvs(cluster, enum.KubeScalingUp)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         fmt.Sprintf("%v/v1", ecs.GroupName),
					Kind:               Kind,
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
							Env:   envs,
						},
					},
				},
			},
		},
	}

	job.Spec.Template.Spec.Containers[0].Env = append(job.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{
			Name:  "NODE_HOSTS",
			Value: convertNodesToString(diffNodeList),
		})
	return job
}

func newScaleDownClusterJob(cluster *ecsv1.KubernetesCluster, diffNodeList []ecsv1.Node) *batchv1.Job {
	// diff work node
	namespace := cluster.Namespace
	name := cluster.Name
	jobName := fmt.Sprintf("scale-down-%v-%v-job-%v", namespace, name, time.Now().Unix())
	completions := pointer.Int32Ptr(1)
	parallelism := pointer.Int32Ptr(1)
	backoffLimit := pointer.Int32Ptr(0)
	// 10 minutes
	ActiveDeadlineSeconds := pointer.Int64Ptr(10 * 60)

	// pack envs
	envs := compressEnvs(cluster, enum.KubeScalingDown)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         fmt.Sprintf("%v/v1", ecs.GroupName), // not define and occur invalid error
					Kind:               Kind,
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
							Env:   envs,
						},
					},
				},
			},
		},
	}

	job.Spec.Template.Spec.Containers[0].Env = append(job.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{
			Name:  "NODE_HOSTS",
			Value: convertNodesToString(diffNodeList),
		})
	return job
}

func compressEnvs(cluster *ecsv1.KubernetesCluster, operation string) []corev1.EnvVar {
	// pack hostsYAML
	hostsYAML := compressHostsYAML(cluster)

	envs := []corev1.EnvVar{
		{
			Name:  "MASTER_HOSTS",
			Value: convertNodesToString(cluster.Spec.MasterList),
		},
		{
			Name:  "MASTER_VIP",
			Value: cluster.Spec.MasterVIP,
		},
		{
			Name:  "ETCD_HOSTS",
			Value: convertNodesToString(cluster.Spec.EtcdList),
		},
		{
			Name:  "OPERATION",
			Value: operation,
		},
		{
			Name:  "HOSTS_YAML",
			Value: hostsYAML,
		},
		{
			Name:  "PRIVATE_KEY",
			Value: cluster.Spec.AuthConfig.PrivateSSHKey,
		},
		{
			Name: "CLUSTER_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name: "CLUSTER_NAMESPACE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				},
			},
		},
	}
	return envs
}
