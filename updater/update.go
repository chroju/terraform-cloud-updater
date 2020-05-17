package updater

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	operator = "<>=~!"
)

type Updater struct {
	Backend          *tfRemoteBackend
	Tfc              *tfcloud.TfCloud
	CurrentVersion   semanticVersion
	RequiredVersions RequiredVersions
	ReleaseVersions  []*tfRelease
}

func NewUpdater(root, token string) (*Updater, error) {
	rb, err := parseTfRemoteBackend(root)
	if err != nil {
		return nil, err
	}

	tfc, err := initializeTfCloud(token)
	if err != nil {
		return nil, err
	}

	releaseVersions, err := getTfReleases()
	if err != nil {
		return nil, err
	}

	currentVersion := make([]int, 3)
	currentVersionString, err := tfc.ReadWorkspaceVersion(rb.Organization, rb.Workspace)
	if err != nil {
		return nil, err
	}
	for i, v := range strings.Split(currentVersionString, ".") {
		currentVersion[i], _ = strconv.Atoi(v)
	}

	if rb.RequiredVersion == "" {
		return &Updater{
			Backend:          rb,
			Tfc:              tfc,
			CurrentVersion:   currentVersion,
			RequiredVersions: nil,
			ReleaseVersions:  releaseVersions,
		}, nil
	}

	trimedRv := strings.TrimSpace(rb.RequiredVersion)
	rvs := parseRequiredVersion(trimedRv)

	return &Updater{
		Backend:          rb,
		Tfc:              tfc,
		CurrentVersion:   currentVersion,
		RequiredVersions: rvs,
		ReleaseVersions:  releaseVersions,
	}, nil
}

func (u *Updater) GetDesiredVersion() (*tfRelease, error) {
	for _, v := range u.ReleaseVersions {
		if v.Draft {
			continue
		}
		if u.RequiredVersions.CheckVersionConsistency(v.SemanticVersion) {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Something wrong")
}

func initializeTfCloud(token string) (*tfcloud.TfCloud, error) {
	if token == "" {
		token = os.Getenv("TFE_TOKEN")
	}
	if token == "" {
		home := os.Getenv("HOME")
		token, _ = parseTerraformrc(home + "/.terraformrc")
	}
	if token == "" {
		return nil, fmt.Errorf("Terraform cloud token is not found")
	}

	tfc, err := tfcloud.NewTfCloud(token)
	if err != nil {
		return nil, err
	}
	return tfc, nil
}

func parseRequiredVersion(versionString string) []*RequiredVersion {
	if strings.Contains(versionString, ",") {
		split := strings.Split(versionString, ",")
		rvs := make([]*RequiredVersion, len(split))
		for i, v := range split {
			rvs[i] = parseRequiredVersion(v)[0]
		}
		return rvs
	}

	var rv *RequiredVersion
	split := strings.Split(versionString, ".")
	var sv []int
	if len(split) != 3 && !strings.Contains(split[0], "~>") {
		sv = []int{0, 0, 0}
	} else {
		sv = make([]int, len(split))
	}
	for i, v := range split {
		if i == 0 && strings.ContainsAny(v, operator) {
			index := strings.LastIndex("", operator)
			rv.Operator = v[:index]
			sv[0], _ = strconv.Atoi(strings.TrimSpace(v[index+1:]))
		}
		sv[i], _ = strconv.Atoi(v)
	}
	rv.SemanticVersion = sv

	return []*RequiredVersion{rv}
}
