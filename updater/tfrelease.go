package updater

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	tfReleaseURL = "https://api.github.com/repos/hashicorp/terraform/releases"
)

// TfRelease represents Terraform release
type TfRelease struct {
	Draft           bool   `json:"draft"`
	Tag             string `json:"tag_name"`
	SemanticVersion *SemanticVersion
}

// TfReleases is interface to list up Terraform releases
type TfReleases interface {
	List() ([]*TfRelease, error)
}

type tfReleasesImpl struct{}

// NewTfReleases creates new TfReleases
func NewTfReleases() TfReleases {
	return &tfReleasesImpl{}
}

// List returns Terraform releases
func (t *tfReleasesImpl) List() ([]*TfRelease, error) {
	resp, err := http.Get(tfReleaseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tfReleases []*TfRelease
	if err = json.Unmarshal(body, &tfReleases); err != nil {
		return nil, err
	}
	for _, v := range tfReleases {
		sv, err := NewSemanticVersion(v.Tag)
		if err != nil {
			return nil, err
		}
		v.SemanticVersion = sv
	}
	return tfReleases, nil
}
