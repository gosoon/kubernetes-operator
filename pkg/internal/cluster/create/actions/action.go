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

package actions

import (
	createtypes "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/types"
	"github.com/gosoon/kubernetes-operator/pkg/internal/util/cli"
	"google.golang.org/grpc"
)

type Action interface {
	Execute(ctx *ActionContext) error
}

// ActionContext is data supplied to all actions
type ActionContext struct {
	Status  *cli.Status
	Cluster *createtypes.ClusterOptions
	Server  *grpc.Server
	Port    string
}

func NewActionContext(
	cluster *createtypes.ClusterOptions,
	server *grpc.Server,
	port string,
	status *cli.Status,
) *ActionContext {
	return &ActionContext{
		Status:  status,
		Cluster: cluster,
		Server:  server,
		Port:    port,
	}
}
