package app

import (
	"fmt"
	"io"
	"os"

	"github.com/gosoon/glog"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/create"
	"github.com/gosoon/kubernetes-operator/pkg/installer/util/protobuf"
	"github.com/pkg/errors"

	"google.golang.org/grpc"
)

// Options xxx
type Options struct {
	Flags   *Flagpole
	Context *cluster.Context
	Server  *grpc.Server
}

// install xxx
type installer struct {
	Options *Options
}

func NewInstaller(opt *Options) *installer {
	return &installer{Options: opt}
}

// CopyFile is installer agent server send file to client
func (s *installer) CopyFile(f *installerv1.File, stream installerv1.Installer_CopyFileServer) error {
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

	return nil
}

// InstallCluster is a grpc server and receiver installer server send clusterProto data
func (s *installer) InstallCluster(stream installerv1.Installer_InstallClusterServer) error {
	cluster := &ecsv1.KubernetesCluster{}
	for {
		clusterProto, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&installerv1.InstallClusterResponse{})
			break
		}
		if err != nil {
			return err
		}

		// receive proto
		fmt.Printf("proto %v \n", clusterProto)
		// receiver operator ecs config
		cluster = protobuf.ClusterConvertToTypes(clusterProto)

	}
	err := s.DoInstallCluster(cluster)
	if err != nil {
		glog.Error(err)
	}
	return nil
}

// TODO : uninstall

// doInstallCluster xxx
func (s *installer) DoInstallCluster(cluster *ecsv1.KubernetesCluster) error {
	// 1.pull image and copy kubeadm,kubectl,kubelet to local
	// 2.init kubeadm config and install
	// 3.callback server
	if err := s.Options.Context.Create(
		create.WithNodeImage(s.Options.Flags.ImageName, s.Options.Flags.Registry, cluster),
		create.WithNodeAddressAndRole(cluster),
		create.WithConfig(cluster),
		create.WithExternalLoadBalancer(cluster),
		create.WithKubeConfigPath(),
		create.WaitForReady(s.Options.Flags.Wait),
	); err != nil {
		//if utilErrors, ok := err.(util.Errors); ok {
		//for _, problem := range utilErrors.Errors() {
		//glog.Error(problem)
		//}
		//return errors.New("aborting due to invalid configuration")
		//}
		return errors.Wrap(err, "failed to create cluster")
	}
	return nil
}
