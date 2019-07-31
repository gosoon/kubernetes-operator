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
	"fmt"
	"time"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"

	"github.com/gosoon/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	// retry
	period = 10 * time.Second

	// ten minutes timeout
	timeout = int64(10 * 60)
)

// jobTTLControl is handle job create failed or job running timeout.
func (c *Controller) jobTTLControl(cluster *ecsv1.KubernetesCluster) {
	name := cluster.Name
	namespace := cluster.Namespace

	stopCh := make(chan struct{})
	defer close(stopCh)

	wait.Until(func() {
		curCluster, err := c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			glog.Errorf("get %s/%s kubernetesCluster object failed with:%v", namespace, name, err)
			stopCh <- struct{}{}
			return
		}

		// if job finished and return
		oldOperation := cluster.Annotations[enum.Operation]
		newOperation := curCluster.Annotations[enum.Operation]
		if oldOperation != newOperation {
			stopCh <- struct{}{}
			return
		}

		// if job running timeout and set operation status to failed
		createTime := cluster.CreationTimestamp.Unix()
		nowTime := time.Now().Unix()
		if nowTime-createTime > timeout {
			curCluster = curCluster.DeepCopy()
			// update kubernetesCluster annotation operation and status
			curCluster.Status.Reason = fmt.Sprintf("the [%v] job running timeout(>%vs)", oldOperation, timeout)
			curCluster.Status.Phase = enum.Failed
			_, err := c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(curCluster)
			if err != nil {
				glog.Errorf("jobTTLTimeout update %s/%s cluster status failed with:%v", namespace, name, err)
			}
			stopCh <- struct{}{}
			return
		}
	}, period, stopCh)
}
