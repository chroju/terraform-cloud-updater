package updater

import (
	"fmt"
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
	{
		Draft:           false,
		Tag:             "v0.12.22",
		SemanticVersion: &SemanticVersion{Versions: []int{0, 12, 22}},
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
			expected: &SemanticVersion{Versions: []int{0, 12, 25}},
		},
		{
			requiredVersions: []*RequiredVersion{
				{
					Operator:        "~>",
					SemanticVersion: &SemanticVersion{Versions: []int{0, 12, 2}},
				},
			},
			expected: &SemanticVersion{Versions: []int{0, 12, 25}},
		},
		{
			requiredVersions: []*RequiredVersion{
				{
					Operator:        "~>",
					SemanticVersion: &SemanticVersion{Versions: []int{0, 12, 0}},
				},
			},
			expected: &SemanticVersion{Versions: []int{0, 12, 25}},
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
					Operator:        ">=",
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

	w, err := newWorkspaceForTest()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for _, v := range cases {
		w.requiredVersions = v.requiredVersions
		err := w.UpdateVersion(v.updateVersion)
		if (err != nil) != v.expectError {
			t.Errorf("Failed: requiredVersions = %v / updateVersion = %v / want = %v", v.requiredVersions, v.updateVersion, v.expectError)
		}
	}
}

func TestUpdateLatestVersion(t *testing.T) {
	cases := []struct {
		requiredVersions RequiredVersions
		expect           *SemanticVersion
	}{
		{
			requiredVersions: []*RequiredVersion{
				{
					Operator:        ">=",
					SemanticVersion: &SemanticVersion{Versions: []int{0, 12, 22}},
				},
			},
			expect: &SemanticVersion{Versions: []int{0, 12, 25}},
		},
		{
			requiredVersions: []*RequiredVersion{
				{
					Operator:        "<",
					SemanticVersion: &SemanticVersion{Versions: []int{0, 12, 24}},
				},
			},
			expect: &SemanticVersion{Versions: []int{0, 12, 23}},
		},
	}

	w, err := newWorkspaceForTest()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	initialVersion := &SemanticVersion{Versions: []int{0, 12, 20}}

	for _, v := range cases {
		_ = w.UpdateVersion(initialVersion)

		w.requiredVersions = v.requiredVersions
		newVersion, err := w.UpdateLatestVersion()
		if err != nil {
			t.Errorf("Failed: %s / requiredVersions = %v /  want = %v", err.Error(), v.requiredVersions, v.expect)
		} else if !reflect.DeepEqual(newVersion, v.expect) {
			t.Errorf("Failed: requiredVersions = %v /  want = %v", v.requiredVersions, v.expect)
		}
	}
}

func newWorkspaceForTest() (*Workspace, error) {
	token := os.Getenv("TFE_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("Please set your Terraform Cloud token to env var TFE_TOKEN")
	}
	org := os.Getenv("TFE_ORG")
	if org == "" {
		return nil, fmt.Errorf("Please set your Terraform Cloud organization to env var TFE_ORG")
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
	return w, nil
}
