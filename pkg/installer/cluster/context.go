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

package cluster

import (
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/create"
	internalcontext "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/context"
	internalcreate "github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create"
	"google.golang.org/grpc"
)

// DefaultName is the default cluster name
const DefaultName = constants.DefaultClusterName

// Context is used to create / manipulate kubernetes-in-docker clusters
// See: NewContext()
type Context struct {
	// the internal context type, shared between implementations of more
	// advanced methods like create
	ic *internalcontext.Context
}

// NewContext returns a new cluster management context
// if name is "" the default name will be used (constants.DefaultClusterName)
func NewContext(name string, server *grpc.Server, port string) *Context {
	// wrap a new internal context
	return &Context{
		ic: internalcontext.NewContext(name, server, port),
	}
}

// Create provisions and starts a kubernetes-in-docker cluster
func (c *Context) Create(options ...create.ClusterOption) error {
	return internalcreate.Cluster(c.ic, options...)
}
