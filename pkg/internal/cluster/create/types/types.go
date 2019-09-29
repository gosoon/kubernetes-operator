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

import (
	"time"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/internal/apis/config"
)

// ClusterOptions holds cluster creation options
// NOTE: this is only exported for usage by the parent package and the options
// package
// See ClusterOption instead
type ClusterOptions struct {
	Name                 string
	Config               *config.Cluster
	NodeImage            string
	WaitForReady         time.Duration
	SetupKubernetes      bool
	NodeAddress          string
	Role                 ecsv1.NodeRole
	ExternalLoadBalancer string
	KubeConfigPath       string
}
