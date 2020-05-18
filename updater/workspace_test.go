package updater

import (
	"reflect"
	"testing"
)

var releases = []*tfRelease{
	{
		Draft:           true,
		SemanticVersion: []int{0, 13, 0},
	},
	{
		Draft:           false,
		SemanticVersion: []int{0, 12, 26},
	},
	{
		Draft:           false,
		SemanticVersion: []int{0, 12, 25},
	},
	{
		Draft:           false,
		SemanticVersion: []int{0, 12, 24},
	},
	{
		Draft:           false,
		SemanticVersion: []int{0, 12, 23},
	},
}

func TestGetDesiredVersion(t *testing.T) {
	cases := []struct {
		requiredVersions []*RequiredVersion
		expected         SemanticVersion
	}{
		{
			requiredVersions: []*RequiredVersion{
				{
					Operator:        ">",
					SemanticVersion: []int{0, 12, 24},
				},
			},
			expected: []int{0, 12, 26},
		},
		{
			requiredVersions: []*RequiredVersion{
				{
					Operator:        "~>",
					SemanticVersion: []int{0, 12},
				},
			},
			expected: []int{0, 12, 26},
		},
		{
			requiredVersions: []*RequiredVersion{
				{
					Operator:        "~>",
					SemanticVersion: []int{0, 12, 0},
				},
			},
			expected: []int{0, 12, 26},
		},
		{
			requiredVersions: []*RequiredVersion{
				{
					Operator:        "<",
					SemanticVersion: []int{0, 12, 26},
				},
				{
					Operator:        ">=",
					SemanticVersion: []int{0, 12, 22},
				},
			},
			expected: []int{0, 12, 25},
		},
		{
			requiredVersions: nil,
			expected:         []int{0, 13, 0},
		},
	}

	w := &Workspace{}
	for _, v := range cases {
		w.requiredVersions = v.requiredVersions
		result, err := w.GetLatestVersion(releases)
		if err != nil {
			t.Errorf("Failed: requiredVersions = %v / err = %s", v.requiredVersions, err)
		} else if reflect.DeepEqual(result, &(v.expected)) {
			t.Errorf("Failed: requiredVersions = %v / want = %v / get = %v", v.requiredVersions, v.expected, result)
		}
	}
}
