package installcni

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/gosoon/kubernetes-operator/pkg/exec"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions"
	"github.com/pkg/errors"
)

type action struct{}

// NewAction returns a new action for installing default CNI
func NewAction() actions.Action {
	return &action{}
}

// Execute runs the action
func (a *action) Execute(ctx *actions.ActionContext) error {
	ctx.Status.Start("Installing CNI")
	defer ctx.Status.End(false)

	//allNodes, err := ctx.Nodes()
	//if err != nil {
	//return err
	//}

	// get the target node for this task
	//node, err := nodes.BootstrapControlPlaneNode(allNodes)
	//if err != nil {
	//return err
	//}

	// read the manifest from the node
	var raw bytes.Buffer
	if err := exec.Command("cat", "/tmp/manifests/default-cni.yaml").SetStdout(&raw).Run(); err != nil {
		return errors.Wrap(err, "failed to read CNI manifest")
	}
	manifest := raw.String()

	// TODO: remove this check?
	// backwards compatibility for mounting your own manifest file to the default
	// location
	// NOTE: this is intentionally undocumented, as an internal implementation
	// detail. Going forward users should disable the default CNI and install
	// their own, or use the default. The internal templating mechanism is
	// not intended for external usage and is unstable.
	if strings.Contains(manifest, "would you kindly template this file") {
		t, err := template.New("cni-manifest").Parse(manifest)
		if err != nil {
			return errors.Wrap(err, "failed to parse CNI manifest template")
		}
		var out bytes.Buffer
		err = t.Execute(&out, &struct {
			PodSubnet string
		}{
			PodSubnet: ctx.Cluster.Config.Networking.PodSubnet,
		})
		if err != nil {
			return errors.Wrap(err, "failed to execute CNI manifest template")
		}
		manifest = out.String()
	}

	// install the manifest
	if err := exec.Command(
		"kubectl", "create", "--kubeconfig=/etc/kubernetes/admin.conf",
		"-f", "-",
	).SetStdin(strings.NewReader(manifest)).Run(); err != nil {
		return errors.Wrap(err, "failed to apply overlay network")
	}

	// mark success
	ctx.Status.End(true)
	return nil
}
