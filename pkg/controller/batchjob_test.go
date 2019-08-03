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
	"testing"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewCreateKubernetesClusterJob(t *testing.T) {
	testCases := []*ecsv1.KubernetesCluster{
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase: enum.Running,
			},
		},
	}
	for _, test := range testCases {
		_ = newCreateKubernetesClusterJob(test)
	}
}
