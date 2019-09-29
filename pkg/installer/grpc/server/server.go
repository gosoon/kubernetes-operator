package server

import (
	"context"
	"sync"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/gosoon/glog"
	"google.golang.org/grpc"
)

// Options is installer server options
type Options struct {
	Server         *grpc.Server
	ServerPort     string
	AgentPort      string
	ImagesRegistry string
}

// gserver xxx
type installer struct {
	Options *Options
}

// NewInstaller is a new installer server
func NewInstaller(opt *Options) *installer {
	return &installer{Options: opt}
}

// CopyFile is a definition of InstallerServer Interface,server do not use CopyFile
// so it is no a implementation
func (s *installer) CopyFile(
	file *installerv1.File,
	stream installerv1.Installer_CopyFileServer) error {

	return nil
}

func (s *installer) InstallCluster(
	ctx context.Context,
	cluster *installerv1.KubernetesClusterRequest) (*installerv1.InstallClusterResponse, error) {

	// images registry,retain,wait and others agent flags is inject by grpc server,
	// because it's don't have to in kubernetes controller
	cluster = s.injectClusterConfig(cluster)

	finish := make(chan bool)
	go s.DispatchClusterConfig(ctx, cluster, finish)

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
func (s *installer) DispatchClusterConfig(
	ctx context.Context,
	cluster *installerv1.KubernetesClusterRequest,
	finish chan<- bool) {

	// get all nodeList
	var clusterNodeList []installerv1.Node
	clusterNodeList = append(clusterNodeList, cluster.Spec.Cluster.NodeList...)
	clusterNodeList = append(clusterNodeList, cluster.Spec.Cluster.MasterList...)

	// dispatch config to echo node and record result in a result channel
	// set dispatch concurrent default is 100
	results := make([]chan types.DispatchConfigResult, len(clusterNodeList))
	chanLimits := make(chan bool, 100)

	var wg sync.WaitGroup
	for idx, node := range clusterNodeList {
		chanLimits <- true
		results[idx] = make(chan types.DispatchConfigResult, 1)
		wg.Add(1)
		go s.dispatchConfig(node.IP, cluster, &wg, chanLimits, results[idx])
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

func (s *installer) dispatchConfig(
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

	address := ip + ":" + s.Options.AgentPort
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

// ClusterNew is creating a new cluster
func (s *installer) ClusterNew(cluster *ecsv1.KubernetesCluster, scaleUpNodeList []ecsv1.Node) error {
	// TODO
	return nil
}

// ClusterScaleUp is scale up a cluster node
func (s *installer) ClusterScaleUp(cluster *ecsv1.KubernetesCluster, scaleUpNodeList []ecsv1.Node) error {
	// TODO
	return nil
}

// ClusterScaleDown is scale down a cluster node
func (s *installer) ClusterScaleDown(cluster *ecsv1.KubernetesCluster, scaleDonwNodeList []ecsv1.Node) error {
	// TODO
	return nil
}

// ClusterTerminating is delete a cluster
func (s *installer) ClusterTerminating(cluster *ecsv1.KubernetesCluster) error {
	// TODO
	return nil
}

//type nodeList []installerv1.Node

//func (n nodeList) Less(i, j int) bool { return n[i].IP < n[j].IP }
//func (n nodeList) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
//func (n nodeList) Len() int           { return len(n) }

//// InjectClusterConfig is set some config by server,eg:image registry and node role
//func (s *installer) InjectClusterConfig(cluster *installerv1.KubernetesClusterRequest) *installerv1.KubernetesClusterRequest {
//cluster.Spec.Cluster.ImagesRegistry = s.Options.ImagesRegistry

//// set master role
//// grpc server select some node and set to ControlPlaneRole,SecondaryControlPlaneRole,WorkerRole
//if len(cluster.Spec.Cluster.MasterList) > 0 {
//masterList := cluster.Spec.Cluster.MasterList
//sort.Sort(nodeList(masterList))
//for idx, master := range cluster.Spec.Cluster.MasterList {
//cluster.Spec.Cluster.MasterList[idx].Role = string(ecsv1.SecondaryControlPlaneRole)
//if master.IP == masterList[0].IP {
//cluster.Spec.Cluster.MasterList[idx].Role = string(ecsv1.ControlPlaneRole)
//}
//}
//}

//// set node role,default all node is worker
//if len(cluster.Spec.Cluster.NodeList) > 0 {
//for idx := range cluster.Spec.Cluster.NodeList {
//cluster.Spec.Cluster.NodeList[idx].Role = string(ecsv1.WorkerRole)
//}
//}

//return cluster
//}
