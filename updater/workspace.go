package updater

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	operator = "<>=~!"
)

// Workspace represents Terraform Cloud workspace
type Workspace struct {
	client           *TfCloud
	hostname         string
	organization     string
	workspace        string
	requiredVersions RequiredVersions
	currentVersion   SemanticVersion
}

// NewWorkspace creates new workspace
func NewWorkspace(token, address string, rb *tfRemoteBackend) (*Workspace, error) {
	tfc, err := NewTfCloud(address, token)
	if err != nil {
		return nil, err
	}

	cv := make([]int, 3)
	currentVersionString, err := tfc.ReadWorkspaceVersion(rb.Organization, rb.Workspace)
	if err != nil {
		return nil, err
	}
	for i, v := range strings.Split(currentVersionString, ".") {
		cv[i], _ = strconv.Atoi(v)
	}

	var rvs RequiredVersions
	if rb.RequiredVersion != "" {
		rvs = NewRequiredVersions(strings.TrimSpace(rb.RequiredVersion))
	}

	return &Workspace{
		client:           tfc,
		hostname:         rb.Hostname,
		organization:     rb.Organization,
		workspace:        rb.Workspace,
		requiredVersions: rvs,
		currentVersion:   cv,
	}, nil
}

// GetCurrentversion get terraform cloud workspace current terraform veresion
func (w *Workspace) GetCurrentVersion() SemanticVersion {
	return w.currentVersion
}

// GetCurrentversion get terraform cloud workspace current terraform veresion
func (w *Workspace) GetLatestVersion(releases []*tfRelease) (SemanticVersion, error) {
	for _, v := range releases {
		if v.Draft {
			continue
		}
		if w.requiredVersions.CheckVersionConsistency(v.SemanticVersion) {
			return v.SemanticVersion, nil
		}
	}
	return nil, fmt.Errorf("Something wrong")
}
