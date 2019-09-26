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

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
)

// CreateCluster is dispatch kubernetes cluster config to all installer agent
func (s *service) CreateCluster(ctx context.Context, region string, name string, cluster *ecsv1.KubernetesCluster) error {

}
