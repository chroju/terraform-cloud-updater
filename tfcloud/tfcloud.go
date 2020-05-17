package tfcloud

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

// TfCloud represents Terraform Cloud
type TfCloud struct {
	*tfe.Client
	ctx context.Context
}

// NewTfCloud creates a new TfCloud API client
func NewTfCloud(token string) (*TfCloud, error) {
	config := &tfe.Config{
		Token: token,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tfcloud := &TfCloud{
		client,
		ctx,
	}

	return tfcloud, nil
}

// ReadWorkspace reads Terraform Cloud workspace
func (t *TfCloud) ReadWorkspace(organization, workspace string) (*tfe.Workspace, error) {
	ws, err := t.Workspaces.Read(t.ctx, organization, workspace)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// ReadWorkspaceVersion reads Terraform Cloud workspace terraform version
func (t *TfCloud) ReadWorkspaceVersion(organization, workspace string) (string, error) {
	ws, err := t.Workspaces.Read(t.ctx, organization, workspace)
	if err != nil {
		return "", err
	}
	return ws.TerraformVersion, nil
}

// UpdateWorkspaceVersion updates Terraform Cloud workspace terraform version
func (t *TfCloud) UpdateWorkspaceVersion(organization, workspace, version string) (*tfe.Workspace, error) {
	oldWs, err := t.ReadWorkspace(organization, workspace)
	if err != nil {
		return nil, err
	}

	if oldWs.TerraformVersion == version {
		return oldWs, nil
	}

	options := tfe.WorkspaceUpdateOptions{
		TerraformVersion: &version,
	}
	newWs, err := t.Client.Workspaces.Update(t.ctx, organization, workspace, options)
	if err != nil {
		return nil, err
	}
	return newWs, nil
}
