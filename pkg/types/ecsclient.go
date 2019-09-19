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

package types

import ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"

// EcsClient xxx
type EcsClient struct {
	Name          string       `json:"name"`
	Namespace     string       `json:"namespace"`
	Region        string       `json:"region"`
	TimeoutMins   string       `json:"timeoutMins"`
	ClusterType   string       `json:"clusterType"`
	ContainerCIDR string       `json:"containerCIDR"`
	ServiceCIDR   string       `json:"serviceCIDR"`
	MasterList    []ecsv1.Node `json:"masterList"`
	MasterVIP     string       `json:"masterVIP"`
	NodeList      []ecsv1.Node `json:"nodeList"`
	EtcdList      []ecsv1.Node `json:"etcdList"`

	// PrivateSSHKey,because of ssh private key has multiple special characters, use base64 encode in it
	PrivateSSHKey string `json:"privateSSHKey"`
	Retry         bool   `json:"retry"`
}

// Callback xxx
type Callback struct {
	Name       string       `json:"name"`
	Namespace  string       `json:"namespace"`
	Region     string       `json:"region"`
	MasterList []ecsv1.Node `json:"masterList"`
	NodeList   []ecsv1.Node `json:"nodeList"`
	EtcdList   []ecsv1.Node `json:"etcdList"`
	KubeConfig string       `json:"kubeconfig"`
	Success    bool         `json:"success"`
	Message    string       `json:"message"`
}
