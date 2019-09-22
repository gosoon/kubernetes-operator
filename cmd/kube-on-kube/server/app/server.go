package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/gosoon/glog"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"
	"github.com/gosoon/kubernetes-operator/pkg/installer/util/protobuf"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"google.golang.org/grpc"
)

type Options struct {
	Flags  *flagpole
	Server *grpc.Server
}

// install xxx
type installer struct {
	Options *Options
}

func NewInstaller(opt *Options) *installer {
	return &installer{Options: opt}
}

func (s *installer) CopyFile(stream installerv1.Installer_CopyFileServer) error {
	return nil
}

func (s *installer) InstallCluster(installerv1.Installer_InstallClusterServer) error {
	return nil
}

// InstallCluster is send KubernetesCluster config to all installer agent
func (s *installer) DispatchClusterConfig(cluster *ecsv1.KubernetesCluster) {
	var clusterNodeList []ecsv1.Node
	clusterNodeList = append(clusterNodeList, cluster.Spec.Cluster.NodeList...)
	clusterNodeList = append(clusterNodeList, cluster.Spec.Cluster.MasterList...)

	results := make([]chan types.DispatchConfigResult, len(clusterNodeList))
	chanLimits := make(chan bool, 100)

	// convet clusterProto
	clusterProto := protobuf.ClusterConvertToProtobuf(cluster)

	var wg sync.WaitGroup
	for idx, node := range clusterNodeList {
		chanLimits <- true
		results[idx] = make(chan types.DispatchConfigResult, 1)
		wg.Add(1)
		go s.dispatchConfig(node, clusterProto, &wg, chanLimits, results[idx])
	}

	wg.Wait()
	// finish <- true
}

func (s *installer) dispatchConfig(
	node ecsv1.Node,
	clusterProto *installerv1.KubernetesClusterRequest,
	wg *sync.WaitGroup,
	chanLimits <-chan bool,
	result chan<- types.DispatchConfigResult) {

	defer wg.Done()
	defer func() { <-chanLimits }()

	// node.IP : port
	address := "127.0.0.1:10023"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		glog.Error(err)
		return
	}
	defer conn.Close()

	failedResult := func(err error) {
		result <- types.DispatchConfigResult{
			Host:    node.IP,
			Success: false,
			Message: err.Error(),
		}
	}

	client := installerv1.NewInstallerClient(conn)

	// connect installer agent
	stream, err := client.InstallCluster(context.Background())
	if err != nil {
		glog.Error(err)
		failedResult(err)
		return
	}

	// send clusterProto to installer agent
	err = stream.Send(clusterProto)
	if err != nil {
		glog.Error(err)
		failedResult(err)
		return
	}

	// grpc client send EOF
	_, err = stream.CloseAndRecv()
	if err != nil {
		glog.Error(err)
		return
	}
}

func (s *installer) CopyFrom(fileName string) error {
	address := "127.0.0.1:10023"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect: %v", err)
	}
	defer conn.Close()
	client := installerv1.NewInstallerClient(conn)

	stream, err := client.CopyFile(context.Background(), &installerv1.File{Name: fileName})
	if err != nil {
		glog.Error(err)
		return err
	}
	defer stream.CloseSend()

	destFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			file, err := stream.Recv()
			if err == io.EOF {
				// read done.
				fmt.Println("read done ")
				close(waitc)
				return
			}
			if err != nil {
				fmt.Println("Failed to receive a note : %v", err)
			}
			_, err = destFile.Write(file.Content)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	<-waitc
	return nil
}
