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

package context

import (
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/constants"
	createtypes "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/types"
	"google.golang.org/grpc"
)

// Context is the private shared context underlying pkg/cluster.Context
//
// NOTE: this is the internal one, it should contain reasonably trivial
// methods that are safe to share between various user facing methods
// pkg/cluster.Context is a superset of this, packages like create and delete
// consume this
type Context struct {
	Name           string
	ClusterOptions *createtypes.ClusterOptions
	Server         *grpc.Server
	Port           string
}

// NewContext returns a new internal cluster management context
// if name is "" the default name will be used
func NewContext(name string, server *grpc.Server, port string) *Context {
	if name == "" {
		name = constants.DefaultClusterName
	}
	return &Context{
		Name:   name,
		Server: server,
		Port:   port,
	}
}
