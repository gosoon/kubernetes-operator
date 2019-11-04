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

package controller

import (
	"encoding/json"

	"github.com/gosoon/glog"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
)

// processClusterNew first determinated use which installer and install cluster.
func (c *Controller) processClusterNew(ecs *ecsv1.KubernetesCluster) error {
	deployMode := ecs.Spec.Cluster.DeployMode

	var err error
	switch deployMode {
	case ecsv1.BinaryDeployMode:
		err = c.sshInstaller.ClusterNew(ecs)
	case ecsv1.ContainerDeployMode:
		err = c.grpcInstaller.ClusterNew(ecs)
	}

	if err != nil {
		glog.Errorf("install cluster %s failed with %v", ecs.Name, err)
		return err
	}

	return nil
}

// processClusterScaleUp is add some nodes to cluster.
func (c *Controller) processClusterScaleUp(ecs *ecsv1.KubernetesCluster) error {
	// diff work nodes
	var oldSpec ecsv1.KubernetesClusterSpec
	oldSpecStr := ecs.Annotations[enum.Spec]
	err := json.Unmarshal([]byte(oldSpecStr), &oldSpec)
	if err != nil {
		glog.Errorf("get old spec failed with:%v", err)
		return err
	}
	nodeList := diffNodeList(oldSpec.Cluster.NodeList, ecs.Spec.Cluster.NodeList, ecs.Annotations[enum.Operation])

	deployMode := ecs.Spec.Cluster.DeployMode
	switch deployMode {
	case ecsv1.BinaryDeployMode:
		err = c.sshInstaller.ClusterScaleUp(ecs, nodeList)
	case ecsv1.ContainerDeployMode:
		err = c.grpcInstaller.ClusterScaleUp(ecs, nodeList)
	}

	if err != nil {
		glog.Errorf("scaleup cluster %s failed with %v", ecs.Name, err)
		return err
	}
	return nil
}

// processClusterScaleDown is delete some nodes to cluster.
func (c *Controller) processClusterScaleDown(ecs *ecsv1.KubernetesCluster) error {
	// diff work nodes
	var oldSpec ecsv1.KubernetesClusterSpec
	oldSpecStr := ecs.Annotations[enum.Spec]
	err := json.Unmarshal([]byte(oldSpecStr), &oldSpec)
	if err != nil {
		glog.Errorf("get old spec failed with:%v", err)
		return err
	}
	nodeList := diffNodeList(oldSpec.Cluster.NodeList, ecs.Spec.Cluster.NodeList, ecs.Annotations[enum.Operation])

	deployMode := ecs.Spec.Cluster.DeployMode
	switch deployMode {
	case ecsv1.BinaryDeployMode:
		err = c.sshInstaller.ClusterScaleDown(ecs, nodeList)
	case ecsv1.ContainerDeployMode:
		err = c.grpcInstaller.ClusterScaleDown(ecs, nodeList)
	}

	if err != nil {
		glog.Errorf("scaledown cluster %s failed with %v", ecs.Name, err)
		return err
	}

	return nil
}

// processClusterTerminating is delete a cluster.
func (c *Controller) processClusterTerminating(ecs *ecsv1.KubernetesCluster) error {
	deployMode := ecs.Spec.Cluster.DeployMode

	var err error
	switch deployMode {
	case ecsv1.BinaryDeployMode:
		err = c.sshInstaller.ClusterTerminating(ecs)
	case ecsv1.ContainerDeployMode:
		err = c.grpcInstaller.ClusterTerminating(ecs)
	}

	if err != nil {
		glog.Errorf("terminating cluster %s failed with %v", ecs.Name, err)
		return err
	}

	return nil
}
