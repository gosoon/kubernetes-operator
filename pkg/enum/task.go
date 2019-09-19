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

package enum

import ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"

const (
	// Operation is operate annotation key
	Operation = "app.kubernetes.io/operation"

	// Spec is echo operate success and record the spec
	Spec = "app.kubernetes.io/spec"

	// state
	KubeCreating          = string("creating")
	KubeCreateFailed      = string("create-failed")
	KubeCreateFinished    = string("create-finished")
	KubeScalingUp         = string("scaling-up")
	KubeScaleUpFailed     = string("scale-up-failed")
	KubeScaleUpFinished   = string("scale-up-finished")
	KubeScalingDown       = string("scaling-down")
	KubeScaleDownFailed   = string("scale-down-failed")
	KubeScaleDownFinished = string("scale-down-finished")
	KubeTerminating       = string("terminating")
	KubeTerminateFinished = string("terminate-finished")
	KubeTerminateFailed   = string("terminate-failed")

	// TODO:to do it
	KubeUpdating         = string("updating")
	KubeUpdateFailed     = string("update-failed")
	KubeUpdateFinished   = string("update-finished")
	KubeRollbacking      = string("rollbacking")
	KubeRollbackFailed   = string("rollback-failed")
	KubeRollbackFinished = string("rollback-finished")

	// phase
	New         = ecsv1.KubernetesOperatorPhase("")
	Creating    = ecsv1.KubernetesOperatorPhase("Creating")
	Prechecking = ecsv1.KubernetesOperatorPhase("Prechecking")
	Scaling     = ecsv1.KubernetesOperatorPhase("Scaling")
	Running     = ecsv1.KubernetesOperatorPhase("Running")
	Failed      = ecsv1.KubernetesOperatorPhase("Failed")
	Terminating = ecsv1.KubernetesOperatorPhase("Terminating")

	// event message
	EcsSyncSuccess = string("ecs synced successfully")

	// event reason
	SyncedSuccess = string("Synced")

	// job events
	CreateKubeJobSuccess      = string("CreateKubeJobSuccess")
	CreateKubeJobFailed       = string("CreateKubeJobFailed")
	CreateScaleUpJobSuccess   = string("CreateScaleUpJobSuccess")
	CreateScaleUpJobFailed    = string("CreateScaleUpJobFailed")
	CreateScaleDownJobSuccess = string("CreateScaleDownJobSuccess")
	CreateScaleDownJobFailed  = string("CreateScaleDownJobFailed")
	SetFinalizersSuccess      = string("SetFinalizersSuccess")
	SetFinalizersFailed       = string("SetFinalizersFailed")
	DeleteKubeJobSuccess      = string("DeleteKubeJobSuccess")
	DeleteKubeJobFailed       = string("DeleteKubeJobFailed")
)
