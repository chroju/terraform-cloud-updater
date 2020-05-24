package updater

import (
	"reflect"
	"testing"
)

func TestNewSemanticVersion(t *testing.T) {
	var cases = []struct {
		src      string
		expected *SemanticVersion
	}{
		{
			src: "v0.12.0",
			expected: &SemanticVersion{
				Versions: []int{0, 12, 0},
				Status:   "",
			},
		},
		{
			src: "0.12.25",
			expected: &SemanticVersion{
				Versions: []int{0, 12, 25},
				Status:   "",
			},
		},
		{
			src: "v0.13.0-beta",
			expected: &SemanticVersion{
				Versions: []int{0, 13, 0},
				Status:   "beta",
			},
		},
		{
			src: "1.1.0-rc2",
			expected: &SemanticVersion{
				Versions: []int{1, 1, 0},
				Status:   "rc2",
			},
		},
	}

	for _, v := range cases {
		got, err := NewSemanticVersion(v.src)
		if err != nil {
			t.Errorf("Failed of error: %s / src = %s / want = %s", err.Error(), v.src, v.expected)
		} else if !reflect.DeepEqual(got, v.expected) {
			t.Errorf("Failed: / src = %s / want = %s / got = %s", v.src, v.expected, got)
		}
	}
}

func TestCheckVersionConsistency(t *testing.T) {
	var cases = []struct {
		src      string
		dst      string
		expected bool
	}{
		{
			src:      "0.12.0",
			dst:      "0.12.0",
			expected: true,
		},
		{
			src:      "= 0.12.0",
			dst:      "0.12.1",
			expected: false,
		},
		{
			src:      "> 0.12.0",
			dst:      "0.13.1",
			expected: true,
		},
		{
			src:      "> 0.12",
			dst:      "0.12.24",
			expected: true,
		},
		{
			src:      ">= 0.12.0",
			dst:      "0.12.0",
			expected: true,
		},
		{
			src:      "!= 0.12.2",
			dst:      "0.12.2",
			expected: false,
		},
		{
			src:      ">= 0.12.2, < 0.12.20",
			dst:      "0.12.19",
			expected: true,
		},
		{
			src:      ">= 0.12.2, < 0.12.20",
			dst:      "0.12.1",
			expected: false,
		},
		{
			src:      "~> 0.12",
			dst:      "0.12.1",
			expected: true,
		},
		{
			src:      "~> 0.12",
			dst:      "0.13.0",
			expected: false,
		},
		{
			src:      "~> 0.12.2",
			dst:      "0.12.5",
			expected: true,
		},
	}
	for _, v := range cases {
		src, _ := NewRequiredVersions(v.src)
		dst, _ := NewSemanticVersion(v.dst)
		if src.CheckVersionConsistency(dst) != v.expected {
			t.Errorf("Failed: src = %v / dst = %v / want = %v", v.src, v.dst, v.expected)
		}
	}
}
