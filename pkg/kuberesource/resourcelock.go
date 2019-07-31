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

package kuberesource

import (
	"os"

	"github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned/scheme"
	"github.com/gosoon/kubernetes-operator/pkg/controller"

	"github.com/gosoon/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

func NewResourceLock(kubeClient *kubernetes.Clientset) (resourcelock.Interface, error) {
	// init eventRecorder
	eventBroadcaster := record.NewBroadcaster()
	eventRecorder := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: controller.ComponentName})
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})

	// init host identity
	id, err := os.Hostname()
	if err != nil {
		glog.Errorf("get hostname error: %v", err)
		return nil, err
	}
	id = id + "_" + string(uuid.NewUUID())

	rl, err := resourcelock.New("endpoints",
		"kube-system",
		controller.ComponentName,
		kubeClient.CoreV1(),
		kubeClient.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: eventRecorder,
		})

	if err != nil {
		glog.Errorf("error creating lock: %v", err)
		return nil, err
	}
	return rl, nil
}
