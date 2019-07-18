package service

import (
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"

	"github.com/pkg/errors"
)

// valid operate
func validOperate(kubernetesCluster *ecsv1.KubernetesCluster) (bool, error) {
	operation := kubernetesCluster.Annotations[enum.Operation]
	if operation == enum.KubeCreating || operation == enum.KubeScalingUp || operation == enum.KubeScalingDown ||
		operation == enum.KubeTerminating {
		return false, errors.Errorf("the latest operation is [%v] and not finished,please wait for", operation)
	}
	return true, nil
}
