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
	"context"

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
	CreateCluster(ctx context.Context, region string, namespace string, name string, clusterInfo *types.EcsClient) error
	DeleteCluster(ctx context.Context, region string, namespace string, name string, clusterInfo *types.EcsClient) error

	// scale
	ScaleUp(ctx context.Context, region string, namespace string, name string, clusterInfo *types.EcsClient) error
	ScaleDown(ctx context.Context, region string, namespace string, name string, clusterInfo *types.EcsClient) error

	// callback
	CreateClusterCallback(ctx context.Context, region string, namespace string, name string, result *types.Callback) error
	ScaleUpCallback(ctx context.Context, region string, namespace string, name string, result *types.Callback) error
	ScaleDownCallback(ctx context.Context, region string, namespace string, name string, result *types.Callback) error
	DeleteClusterCallback(ctx context.Context, region string, namespace string, name string, result *types.Callback) error

	// logs
	GetClusterOperationLogs(ctx context.Context, region string, namespace string, name string) (string, error)
}
