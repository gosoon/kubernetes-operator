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

package agent

import (
	"context"
	"io"
	"os"
	"time"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"
	installcluster "github.com/gosoon/kubernetes-operator/pkg/installer/cluster"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/create"
	"github.com/gosoon/kubernetes-operator/pkg/installer/util/protobuf"

	"github.com/gosoon/glog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Options xxx
type Options struct {
	Server *grpc.Server
	Port   string
	Wait   time.Duration
}

// agent xxx
type agent struct {
	Options *Options
}

// NewAgent xxx
func NewAgent(opt *Options) *agent {
	return &agent{Options: opt}
}

// CopyFile is installer agent server send file to client
func (s *agent) CopyFile(f *installerv1.File, stream installerv1.Installer_CopyFileServer) error {
	file, err := os.Open(f.Name)
	if err != nil {
		glog.Error(err)
		return err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	tmpFile := &installerv1.File{Name: f.Name}
	for {
		_, err := file.Read(buf)
		if err == io.EOF {
			return nil
		}

		tmpFile.Content = buf
		if err := stream.Send(tmpFile); err != nil {
			glog.Error(err)
			return err
		}
	}
}

// InstallCluster is a grpc server and receiver installer server send clusterProto data
func (s *agent) InstallCluster(
	ctx context.Context,
	clusterRequest *installerv1.KubernetesClusterRequest) (*installerv1.InstallClusterResponse, error) {

	cluster := &ecsv1.KubernetesCluster{}
	// receiver operator ecs config
	cluster, _ = protobuf.ClusterConvertToTypes(clusterRequest)

	// set role for all nodes
	cluster = injectClusterNodeRole(cluster)

	err := s.ClusterNew(cluster)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return nil, nil
}

// ClusterNew xxx
// 1.pull image and copy kubeadm,kubectl,kubelet to local
// 2.init kubeadm config and install
// 3.callback server
func (s *agent) ClusterNew(cluster *ecsv1.KubernetesCluster) error {
	// create a cluster context and create the cluster
	ctx := installcluster.NewContext(constants.DefaultClusterName, s.Options.Server, s.Options.Port)

	if err := ctx.Create(
		create.WithNodeImage(cluster),
		create.WithNodeAddressAndRole(cluster),
		create.WithConfig(cluster),
		create.WithExternalLoadBalancer(cluster),
		create.WithKubeConfigPath(),
		create.WaitForReady(s.Options.Wait),
	); err != nil {
		return errors.Wrap(err, "failed to create cluster")
	}

	return nil
}
