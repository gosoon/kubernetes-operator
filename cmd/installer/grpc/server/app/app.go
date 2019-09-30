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
	grpcserver "github.com/gosoon/kubernetes-operator/pkg/installer/grpc/server"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// flagpole contains all required flags by grpc agent and server
// because of agent deploy in all nodes,so agent required flags are passed through the server.
type flagpole struct {
	ServerPort     string
	AgentPort      string
	ImagesRegistry string
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
	cmd.Flags().StringVar(&flags.ServerPort, "serverPort", "10022", "installer grpc server port(default 10022)")
	cmd.Flags().StringVar(&flags.AgentPort, "agentPort", "10023", "installer grpc agent port(default 10023)")
	cmd.Flags().StringVar(&flags.ImagesRegistry, "registry", "registry.cn-hangzhou.aliyuncs.com/aliyun_kube_system", "kubernetes image registry")
	return cmd
}

// run is start grpc gateway
func run(flags *flagpole) {
	grpcServerEndpoint := ":" + flags.ServerPort
	l, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	installer := grpcserver.NewInstaller(&grpcserver.Options{
		ServerPort:     flags.ServerPort,
		AgentPort:      flags.AgentPort,
		ImagesRegistry: flags.ImagesRegistry,
		Server:         grpcServer,
	})

	// register grpc server
	installerv1.RegisterInstallerServer(grpcServer, installer)
	reflection.Register(grpcServer)
	go func() {
		glog.Info("starting grpc server...")
		glog.Fatal(grpcServer.Serve(l))
	}()

	// start http server
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
}
