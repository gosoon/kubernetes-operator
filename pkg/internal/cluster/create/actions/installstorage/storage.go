package installstorage

import (
	"strings"

	"github.com/gosoon/kubernetes-operator/pkg/exec"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions"
	"github.com/pkg/errors"
)

type action struct{}

// NewAction returns a new action for installing storage
func NewAction() actions.Action {
	return &action{}
}

// Execute runs the action
func (a *action) Execute(ctx *actions.ActionContext) error {
	ctx.Status.Start("Installing StorageClass")
	defer ctx.Status.End(false)

	// add the default storage class
	if err := addDefaultStorageClass(); err != nil {
		return errors.Wrap(err, "failed to add default storage class")
	}

	// mark success
	ctx.Status.End(true)
	return nil
}

// a default storage class
// we need this for e2es (StatefulSet)
const defaultStorageClassManifest = `# host-path based default storage class
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  namespace: kube-system
  name: standard
  annotations:
    storageclass.beta.kubernetes.io/is-default-class: "true"
  labels:
    addonmanager.kubernetes.io/mode: EnsureExists
provisioner: kubernetes.io/host-path`

func addDefaultStorageClass() error {
	in := strings.NewReader(defaultStorageClassManifest)
	cmd := exec.Command(
		"kubectl",
		"--kubeconfig=/etc/kubernetes/admin.conf", "apply", "-f", "-",
	)
	cmd.SetStdin(in)
	return cmd.Run()
}
