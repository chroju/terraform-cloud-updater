package updater

import (
	"os"
	"reflect"
	"testing"
)

var releases = []*TfRelease{
	{
		Draft:           true,
		Tag:             "v0.13.0",
		SemanticVersion: &SemanticVersion{Versions: []int{0, 13, 0}},
	},
	{
		Draft:           false,
		Tag:             "v0.12.26",
		SemanticVersion: &SemanticVersion{Versions: []int{0, 12, 26}},
	},
	{
		Draft:           false,
		Tag:             "v0.12.25",
		SemanticVersion: &SemanticVersion{Versions: []int{0, 12, 25}},
	},
	{
		Draft:           false,
		Tag:             "v0.12.24",
		SemanticVersion: &SemanticVersion{Versions: []int{0, 12, 24}},
	},
	{
		Draft:           false,
		Tag:             "v0.12.23",
		SemanticVersion: &SemanticVersion{Versions: []int{0, 12, 23}},
	},
}

type TfReleasesMock struct{}

func (t *TfReleasesMock) List() ([]*TfRelease, error) {
	return releases, nil
}

func TestGetLatestVersion(t *testing.T) {
	cases := []struct {
		requiredVersions RequiredVersions
		expected         *SemanticVersion
	}{
		{
			requiredVersions: []*RequiredVersion{
				{
					Operator:        "~>",
					SemanticVersion: &SemanticVersion{Versions: []int{0, 12}},
				},
			},
			expected: &SemanticVersion{Versions: []int{0, 12, 26}},
		},
	}

	w := &Workspace{
		tfRelease: &TfReleasesMock{},
	}
	for _, v := range cases {
		w.requiredVersions = v.requiredVersions
		result, err := w.GetLatestVersion()
		if err != nil {
			t.Errorf("Failed: requiredVersions = %v / err = %s", v.requiredVersions, err)
		} else if reflect.DeepEqual(result, &(v.expected)) {
			t.Errorf("Failed: requiredVersions = %v / want = %v / get = %v", v.requiredVersions, v.expected, result)
		}
	}
}

func TestUpdateVersion(t *testing.T) {
	cases := []struct {
		requiredVersions RequiredVersions
		updateVersion    *SemanticVersion
		expectError      bool
	}{
		{
			requiredVersions: []*RequiredVersion{
				{
					Operator:        ">",
					SemanticVersion: &SemanticVersion{Versions: []int{0, 12, 24}},
				},
			},
			updateVersion: &SemanticVersion{Versions: []int{0, 12, 25}},
			expectError:   false,
		},
		{
			requiredVersions: []*RequiredVersion{
				{
					Operator:        "<",
					SemanticVersion: &SemanticVersion{Versions: []int{0, 12, 24}},
				},
			},
			updateVersion: &SemanticVersion{Versions: []int{0, 12, 25}},
			expectError:   true,
		},
	}

	token := os.Getenv("TFE_TOKEN")
	if token == "" {
		t.Error("Please set your Terraform Cloud token to env var TFE_TOKEN")
	}
	org := os.Getenv("TFE_ORG")
	if org == "" {
		t.Error("Please set your Terraform Cloud organization to env var TFE_ORG")
	}
	workspace := os.Getenv("TFE_WORKSPACE")
	if workspace == "" {
		workspace = "sample"
	}

	client, _ := NewTfCloud("app.terraform.io", token)
	config := &Config{
		Organization: org,
		Workspace:    workspace,
	}
	w, _ := NewWorkspace(client, config)
	for _, v := range cases {
		w.requiredVersions = v.requiredVersions
		err := w.UpdateVersion(v.updateVersion)
		if (err != nil) != v.expectError {
			t.Errorf("Failed: requiredVersions = %v / updateVersion = %v / want = %v", v.requiredVersions, v.updateVersion, v.expectError)
		}
	}
}
