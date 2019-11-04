package server

import (
	"context"

	"github.com/gosoon/glog"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/installer/util/protobuf"
)

// ClusterNew is creating a new cluster
func (inst *installer) ClusterNew(cluster *ecsv1.KubernetesCluster) error {
	clusterRequest, err := protobuf.ClusterConvertToProtobuf(cluster)
	if err != nil {
		glog.Errorf("clusterRequest convert to protobuf failed with:%v", err)
		return err
	}

	_, err = inst.InstallCluster(context.Background(), clusterRequest)
	if err != nil {
		glog.Errorf("installCluster failed with %v", err)
		return err
	}

	return nil
}

// ClusterScaleUp is scale up a cluster node
func (inst *installer) ClusterScaleUp(cluster *ecsv1.KubernetesCluster, scaleUpNodeList []ecsv1.Node) error {
	// TODO
	return nil
}

// ClusterScaleDown is scale down a cluster node
func (inst *installer) ClusterScaleDown(cluster *ecsv1.KubernetesCluster, scaleDownNodeList []ecsv1.Node) error {
	// TODO
	return nil
}

// ClusterTerminating is delete a cluster
func (inst *installer) ClusterTerminating(cluster *ecsv1.KubernetesCluster) error {
	// TODO
	return nil
}
