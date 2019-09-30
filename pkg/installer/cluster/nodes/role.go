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

package nodes

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/internal/apis/config"

	"github.com/pkg/errors"
)

// BootstrapControlPlaneNode returns a handle to the bootstrap control plane node
func BootstrapControlPlaneNode(allNodes []config.Node) (*config.Node, error) {
	controlPlaneNode := &config.Node{}
	for _, node := range allNodes {
		if node.Role == ecsv1.ControlPlaneRole {
			controlPlaneNode.IP = node.IP
			controlPlaneNode.Role = node.Role
		}
	}

	if controlPlaneNode == nil {
		return nil, errors.Errorf(
			"expected at least one %s node",
			ecsv1.ControlPlaneRole,
		)
	}
	return controlPlaneNode, nil
}
