package updater

import (
	"context"
	"strconv"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// TfCloud represents Terraform Cloud API wrapper
type TfCloud interface {
	ReadWorkspaceVersion(org, workspace string) (*SemanticVersion, error)
	UpdateWorkspaceVersion(org, workspace string, sv *SemanticVersion) error
	// *tfe.Client
	// ctx context.Context
}

type tfcloudImpl struct {
	*tfe.Client
	ctx context.Context
}

// NewTfCloud creates a new TfCloud interface
func NewTfCloud(address, token string) (TfCloud, error) {
	config := &tfe.Config{
		Token: token,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return &tfcloudImpl{
		client,
		ctx,
	}, nil
}

func (t *tfcloudImpl) readWorkspace(organization, workspace string) (*tfe.Workspace, error) {
	ws, err := t.Workspaces.Read(t.ctx, organization, workspace)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// ReadWorkspaceVersion reads Terraform Cloud workspace terraform version
func (t *tfcloudImpl) ReadWorkspaceVersion(org, workspace string) (*SemanticVersion, error) {
	ws, err := t.Workspaces.Read(t.ctx, org, workspace)
	if err != nil {
		return nil, err
	}

	sv := make([]int, 3)
	for i, v := range strings.Split(ws.TerraformVersion, ".") {
		sv[i], _ = strconv.Atoi(v)
	}

	return &SemanticVersion{Versions: sv}, nil
}

// UpdateWorkspaceVersion updates Terraform Cloud workspace terraform version
func (t *tfcloudImpl) UpdateWorkspaceVersion(org, workspace string, sv *SemanticVersion) error {
	oldWs, err := t.readWorkspace(org, workspace)
	if err != nil {
		return err
	}

	if oldWs.TerraformVersion == sv.String() {
		return nil
	}

	options := tfe.WorkspaceUpdateOptions{
		TerraformVersion: tfe.String(sv.String()),
	}
	if _, err = t.Client.Workspaces.Update(t.ctx, org, workspace, options); err != nil {
		return err
	}
	return nil
}
