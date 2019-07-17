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

	KubeUpdating          = string("updating")
	KubeUpdateFailed      = string("update-failed")
	KubeUpdateFinished    = string("update-finished")
	KubeRollbacking       = string("rollbacking")
	KubeRollbackFailed    = string("rollback-failed")
	KubeRollbackFinished  = string("rollback-finished")
	KubeTerminating       = string("terminating")
	KubeTerminateFinished = string("terminate-finished")
	KubeTerminateFailed   = string("terminate-failed")

	// phase
	New         = ecsv1.KubernetesOperatorPhase("")
	Scaling     = ecsv1.KubernetesOperatorPhase("Scaling")
	Creating    = ecsv1.KubernetesOperatorPhase("Creating")
	Running     = ecsv1.KubernetesOperatorPhase("Running")
	Failed      = ecsv1.KubernetesOperatorPhase("Failed")
	Terminating = ecsv1.KubernetesOperatorPhase("Terminating")

	// event message
	EcsSyncSuccess = string("ecs synced successfully")

	// event reason
	SyncedSuccess = string("Synced")

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
