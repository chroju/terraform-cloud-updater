package updater

import (
	"testing"
)

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
