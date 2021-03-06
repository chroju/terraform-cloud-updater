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
	Hostname        string
}

// NewWorkspace creates new workspace
func NewWorkspace(tfcloud TfCloud, config *Config) (*Workspace, error) {
	hostname := "app.terraform.io"
	if config.Hostname != "" {
		hostname = config.Hostname
	}

	ws := &Workspace{
		client:           tfcloud,
		hostname:         hostname,
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

// GetRequiredVersions get required versions
func (w *Workspace) GetRequiredVersions() *RequiredVersions {
	return &w.requiredVersions
}

// GetSettingsLink get workspace settings link
func (w *Workspace) GetSettingsLink() string {
	return fmt.Sprintf("https://%s/app/%s/workspaces/%s/settings/general", w.hostname, w.organization, w.workspace)
}

// GetCurrentVersion get terraform cloud workspace current terraform veresion
func (w *Workspace) GetCurrentVersion() (*SemanticVersion, error) {
	cv, err := w.client.ReadWorkspaceVersion(w.organization, w.workspace)
	if err != nil {
		return nil, err
	}
	return cv, nil
}

// GetLatestVersion get latest terraform version
func (w *Workspace) GetLatestVersion() (*SemanticVersion, error) {
	releases, err := w.tfRelease.List()
	if err != nil {
		return nil, err
	}

	for _, v := range releases {
		if v.Draft {
			continue
		}
		return v.SemanticVersion, nil
	}

	return nil, fmt.Errorf("Something is wrong to get latest terraform version")
}

// GetCompatibleLatestVersion get latest terraform version compatible with required versions
func (w *Workspace) GetCompatibleLatestVersion() (*SemanticVersion, error) {
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

	return nil, fmt.Errorf("No version is compatbile with required versions '%v'", w.requiredVersions)
}

// UpdateVersion update terraform cloud workspace terraform version
func (w *Workspace) UpdateVersion(s *SemanticVersion) error {
	if !w.requiredVersions.CheckVersionConsistency(s) {
		return fmt.Errorf("Version %v is not compatbile with required version '%v'", s, w.requiredVersions)
	}
	if err := w.client.UpdateWorkspaceVersion(w.organization, w.workspace, s); err != nil {
		return err
	}
	return nil
}

// UpdateLatestVersion update terraform cloud workspace terraform to the latest version
func (w *Workspace) UpdateCompatibleLatestVersion() (*SemanticVersion, error) {
	newVersion, err := w.GetCompatibleLatestVersion()
	if err != nil {
		return nil, err
	}

	if err = w.UpdateVersion(newVersion); err != nil {
		return nil, err
	}
	return newVersion, nil
}

// IsCompatibleVersion returns whether a given veresion is compatible with workspace required version
func (w *Workspace) IsCompatibleVersion(s *SemanticVersion) bool {
	return w.requiredVersions.CheckVersionConsistency(s)
}
