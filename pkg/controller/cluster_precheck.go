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
	"context"
	"fmt"
	"time"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/gosoon/kubernetes-operator/pkg/precheck"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/gosoon/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// precheckTimeout must be greater than sshServerTimeout
	precheckTimeout  time.Duration = 35 * time.Second
	sshServerTimeout time.Duration = 30 * time.Second
)

func (c *Controller) processClusterPrecheck(cluster *ecsv1.KubernetesCluster) error {
	// start precheck
	success, msg := runPrecheck(cluster)

	// if precheck is success and process cluster operate
	if success {
		switch cluster.Annotations[enum.Operation] {
		case enum.KubeCreating:
			return c.processClusterNew(cluster)
		case enum.KubeScalingUp:
			return c.processClusterScaleUp(cluster)
		case enum.KubeScalingDown:
			return c.processClusterScaleDown(cluster)
		}
	}

	curCluster := cluster.DeepCopy()
	namespace := curCluster.Namespace
	// handle precheck failed
	curCluster.Status.Phase = enum.Failed
	curCluster.Status.Reason = fmt.Sprintf("[precheck] failed with:\n %v", msg)
	curCluster.Status.LastTransitionTime = metav1.Now()
	curCluster, err := c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
	if err != nil {
		glog.Errorf("update status to [Failed] failed with:%v", err)
	}
	return err
}

// runPrecheck is range all node list and exec check script
func runPrecheck(cluster *ecsv1.KubernetesCluster) (bool, string) {
	results := make([]chan types.PrecheckResult, len(cluster.Spec.Cluster.NodeList))
	finished := make(chan bool)

	ctx, cacnel := context.WithTimeout(context.Background(), precheckTimeout)
	defer cacnel()

	// new precheck
	precheckCli := precheck.New(&precheck.Options{
		Timeout: sshServerTimeout,
	})

	// run precheck
	go precheckCli.HostEnv(cluster, results, finished)

	checkResult := true
	checkMsg := ""
	select {
	case <-ctx.Done():
		msg := fmt.Sprintf("[precheck] %s timeout", cluster.Name)
		checkResult = false
		checkMsg = msg

	case <-finished:
		for _, result := range results {
			res := <-result
			if !res.Success {
				checkResult = false
				checkMsg += fmt.Sprintf("%v : %v\n", res.Host, res.Result)
			}
		}
	}
	return checkResult, checkMsg
}
