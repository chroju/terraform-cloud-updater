package updater

import (
	"fmt"
	"strings"
)

// Workspace represents Terraform Cloud workspace
type Workspace struct {
	client           TfCloud
	tfRelease        TfReleases
	hostname         string
	organization     string
	workspace        string
	requiredVersions RequiredVersions
}

// Config is Terraform Cloud workspace config
type Config struct {
	Organization    string
	Workspace       string
	RequiredVersion string
}

// NewWorkspace creates new workspace
func NewWorkspace(tfcloud TfCloud, config *Config) (*Workspace, error) {
	ws := &Workspace{
		client:           tfcloud,
		organization:     config.Organization,
		workspace:        config.Workspace,
		requiredVersions: nil,
	}
	ws.tfRelease = NewTfReleases()

	if config.RequiredVersion != "" {
		rvs, err := NewRequiredVersions(strings.TrimSpace(config.RequiredVersion))
		if err != nil {
			return nil, err
		}
		ws.requiredVersions = rvs
	}

	return ws, nil
}

// GetCurrentVersion get terraform cloud workspace current terraform veresion
func (w *Workspace) GetCurrentVersion() (*SemanticVersion, error) {
	cv, err := w.client.ReadWorkspaceVersion(w.organization, w.workspace)
	if err != nil {
		return nil, err
	}
	return cv, nil
}

// GetLatestVersion get terraform cloud workspace current terraform veresion
func (w *Workspace) GetLatestVersion() (*SemanticVersion, error) {
	releases, err := w.tfRelease.List()
	if err != nil {
		return nil, err
	}

	for _, v := range releases {
		if v.Draft {
			continue
		}
		if w.requiredVersions.CheckVersionConsistency(v.SemanticVersion) {
			return v.SemanticVersion, nil
		}
	}

	return nil, fmt.Errorf("No version is consistent with required versions %v", w.requiredVersions)
}

// UpdateVersion update terraform cloud workspace terraform version
func (w *Workspace) UpdateVersion(s *SemanticVersion) error {
	if !w.requiredVersions.CheckVersionConsistency(s) {
		return fmt.Errorf("Version %v is not consistent with required version '%v'", s, w.requiredVersions)
	}
	if err := w.client.UpdateWorkspaceVersion(w.organization, w.workspace, s); err != nil {
		return err
	}
	return nil
}
