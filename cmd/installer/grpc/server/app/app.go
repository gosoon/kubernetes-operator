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

package app

import (
	"context"
	"net"
	"net/http"

	"github.com/gosoon/glog"
	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type flagpole struct {
	Port string
}

// NewServerCommand returns a new cobra.Command for kube-on-kube server
func NewServerCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "cluster",
		Short: "Creates a local Kubernetes cluster",
		Long:  "Creates a local Kubernetes cluster",
		Run: func(cmd *cobra.Command, args []string) {
			run(flags)
		},
	}
	cmd.Flags().StringVar(&flags.Port, "port", "10022", "installer agent grpc server port(default 10022)")
	return cmd
}

// run is start grpc gateway
func run(flags *flagpole) {
	// start grpc server
	grpcServerEndpoint := ":" + flags.Port
	l, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	installer := NewInstaller(&Options{
		Flags:  flags,
		Server: grpcServer,
	})
	// register grpc server
	installerv1.RegisterInstallerServer(grpcServer, installer)
	reflection.Register(grpcServer)
	go func() {
		glog.Info("starting grpc server...")
		glog.Fatal(grpcServer.Serve(l))
	}()
	// start http server
	//
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = installerv1.RegisterInstallerHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		glog.Fatal(err)
	}

	glog.Info("starting http server...")
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	glog.Fatal(http.ListenAndServe(":8080", mux))

	//master := []ecsv1.Node{{IP: "127.0.0.1"}}
	//cluster := &ecsv1.KubernetesCluster{Spec: ecsv1.KubernetesClusterSpec{Cluster: ecsv1.Cluster{MasterList: master}}}
	//installer.DispatchClusterConfig(cluster)
}
