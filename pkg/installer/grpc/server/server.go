package server

import (
	"context"
	"net"
	"net/http"
	"sync"

	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"
	"github.com/gosoon/kubernetes-operator/pkg/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/gosoon/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	ServerPort     = "10022"
	AgentPort      = "10023"
	ImagesRegistry = "registry.cn-hangzhou.aliyuncs.com/aliyun_kube_system"
)

// Options is installer server options
type Options struct {
	ServerPort     string
	AgentPort      string
	ImagesRegistry string
}

// gserver xxx
type installer struct {
	opt *Options
}

// NewInstaller is a new installer server
func NewInstaller(opt *Options) *installer {
	grpcInstaller := &installer{opt: opt}

	//start grpc installer's grpc server and http server
	go grpcInstaller.Run()

	return grpcInstaller
}

// run is start grpc gateway
// TODO: remove http server,only use grcp server
func (inst *installer) Run() {
	grpcServerEndpoint := ":" + inst.opt.ServerPort
	l, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// register grpc server
	installerv1.RegisterInstallerServer(grpcServer, inst)
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

// CopyFile is a definition of InstallerServer Interface,server do not use CopyFile
// so it is no a implementation
func (inst *installer) CopyFile(
	file *installerv1.File,
	stream installerv1.Installer_CopyFileServer) error {

	return nil
}

func (inst *installer) InstallCluster(
	ctx context.Context,
	cluster *installerv1.KubernetesClusterRequest) (*installerv1.InstallClusterResponse, error) {

	// images registry,retain,wait and others agent flags is inject by grpc server,
	// because it's don't have to in kubernetes controller
	cluster = inst.injectClusterConfig(cluster)

	finish := make(chan bool)
	go inst.DispatchClusterConfig(ctx, cluster, finish)

	select {
	case <-ctx.Done():
		// callback install failed
	case success := <-finish:
		if !success {
			// callback install failed
		}
	}
	return nil, nil
}

// InstallCluster is send KubernetesCluster config to all installer agent
func (inst *installer) DispatchClusterConfig(
	ctx context.Context,
	cluster *installerv1.KubernetesClusterRequest,
	finish chan<- bool) {

	// get all nodeList
	var clusterNodeList []installerv1.Node
	clusterNodeList = append(clusterNodeList, cluster.Spec.Cluster.NodeList...)
	clusterNodeList = append(clusterNodeList, cluster.Spec.Cluster.MasterList...)

	// dispatch config to echo node and record result in a result channel
	// set dispatch concurrent default is 16
	results := make([]chan types.DispatchConfigResult, len(clusterNodeList))
	chanLimits := make(chan bool, 16)

	var wg sync.WaitGroup
	for idx, node := range clusterNodeList {
		chanLimits <- true
		results[idx] = make(chan types.DispatchConfigResult, 1)
		wg.Add(1)
		go inst.dispatchConfig(node.IP, cluster, &wg, chanLimits, results[idx])
	}
	wg.Wait()

	// get all result by node
	success := true
	for _, result := range results {
		res := <-result
		// if dispatchConfig failed and record log
		if !res.Success {
			success = false
			glog.Errorf("dispatchConfig res.Host %v failed with:%v", res.Host, res.Message)
		}
	}
	finish <- success
}

func (inst *installer) dispatchConfig(
	ip string,
	cluster *installerv1.KubernetesClusterRequest,
	wg *sync.WaitGroup,
	chanLimits <-chan bool,
	result chan<- types.DispatchConfigResult) {

	defer wg.Done()
	defer func() { <-chanLimits }()

	failedResult := func(err error) {
		result <- types.DispatchConfigResult{
			Host:    ip,
			Success: false,
			Message: err.Error(),
		}
	}

	address := ip + ":" + inst.opt.AgentPort
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		failedResult(err)
		return
	}
	defer conn.Close()

	client := installerv1.NewInstallerClient(conn)
	//send cluster config to installer agent
	_, err = client.InstallCluster(context.Background(), cluster)
	if err != nil {
		failedResult(err)
		return
	}
}
