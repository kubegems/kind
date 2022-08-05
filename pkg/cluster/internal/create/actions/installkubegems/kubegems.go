package installkubegems

import (
	"sigs.k8s.io/kind/pkg/errors"

	"sigs.k8s.io/kind/pkg/cluster/internal/create/actions"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"
)

type action struct{}

// NewAction returns a new action for installing KubeGems
func NewAction() actions.Action {
	return &action{}
}

// Execute runs the action
func (a *action) Execute(ctx *actions.ActionContext) error {
	ctx.Status.Start("Installing KubeGems 🎁")
	defer ctx.Status.End(false)

	allNodes, err := ctx.Nodes()
	if err != nil {
		return err
	}

	// get the target node for this task
	controlPlanes, err := nodeutils.ControlPlaneNodes(allNodes)
	if err != nil {
		return err
	}
	node := controlPlanes[0] // kind expects at least one always

	// install the manifest

	if err := node.Command(
		"kubectl", "create", "--kubeconfig=/etc/kubernetes/admin.conf",
		"namespace", "kubegems-installer",
	).Run(); err != nil {
		return errors.Wrap(err, "failed to apply kubegems installer namespace")
	}
	if err := node.Command(
		"kubectl", "create", "--kubeconfig=/etc/kubernetes/admin.conf",
		"-f", "https://github.com/kubegems/kubegems/raw/main/deploy/installer.yaml",
	).Run(); err != nil {
		return errors.Wrap(err, "failed to apply kubegems installer")
	}
	if err := node.Command(
		"kubectl", "create", "--kubeconfig=/etc/kubernetes/admin.conf",
		"namespace", "kubegems",
	).Run(); err != nil {
		return errors.Wrap(err, "failed to apply kubegems namespace")
	}
	if err := node.Command(
		"kubectl", "create", "--kubeconfig=/etc/kubernetes/admin.conf",
		"-f", "https://raw.githubusercontent.com/kubegems/kind/main/kubegems-kind.yaml",
	).Run(); err != nil {
		return errors.Wrap(err, "failed to apply kubegems")
	}

	// mark success
	ctx.Status.End(true)
	return nil
}
