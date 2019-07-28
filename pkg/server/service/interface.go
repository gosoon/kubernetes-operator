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

package service

import (
	clientset "github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"k8s.io/client-go/kubernetes"
)

type Options struct {
	KubernetesClusterClientset clientset.Interface
	KubeClientset              kubernetes.Interface
}

type service struct {
	opt *Options
}

func New(opt *Options) Interface {
	return &service{opt: opt}
}

type Interface interface {
	// cluster
	CreateCluster(region string, namespace string, name string, clusterInfo *types.EcsClient) error
	DeleteCluster(region string, namespace string, name string, clusterInfo *types.EcsClient) error

	// scale
	ScaleUp(region string, namespace string, name string, clusterInfo *types.EcsClient) error
	ScaleDown(region string, namespace string, name string, clusterInfo *types.EcsClient) error

	// callback
	CreateClusterCallback(region string, namespace string, name string, result *types.Callback) error
	ScaleUpCallback(region string, namespace string, name string, result *types.Callback) error
	ScaleDownCallback(region string, namespace string, name string, result *types.Callback) error
	DeleteClusterCallback(region string, namespace string, name string, result *types.Callback) error

	// logs
	GetClusterOperationLogs(region string, namespace string, name string) (string, error)
}
