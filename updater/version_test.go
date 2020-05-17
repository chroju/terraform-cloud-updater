package updater

import (
	"testing"
)

var cases = []struct {
	src semanticVersion
	dst semanticVersion
}{
	{
		src: []int{0, 12, 0},
		dst: []int{0, 12, 0},
	},
	{
		src: []int{0, 12, 0},
		dst: []int{0, 12, 1},
	},
	{
		src: []int{0, 12, 0},
		dst: []int{0, 11, 1},
	},
	{
		src: []int{0, 12, 14},
		dst: []int{0, 12, 14},
	},
	{
		src: []int{1, 1, 14},
		dst: []int{0, 12, 14},
	},
	{
		src: []int{0, 11},
		dst: []int{0, 12, 14},
	},
	{
		src: []int{0, 12},
		dst: []int{0, 12, 14},
	},
	{
		src: []int{0, 13},
		dst: []int{0, 12, 14},
	},
	{
		src: []int{1, 0},
		dst: []int{0, 12, 14},
	},
}

func TestIsEquall(t *testing.T) {
	expected := []bool{true, false, false, true, false, false, false, false, false}
	for i, v := range cases {
		if v.src.IsEquall(v.dst) != expected[i] {
			t.Errorf("Failed: src = %v / dst = %v / expected = %v", v.src, v.dst, expected[i])
		}
	}
}

func TestIsGreaterThan(t *testing.T) {
	expected := []bool{false, false, true, false, true, true, true, true, true}
	for i, v := range cases {
		if len(v.src) == 2 {
			continue
		}
		if v.src.IsGreaterThan(v.dst) != expected[i] {
			t.Errorf("Failed: src = %v / dst = %v / expected = %v", v.src, v.dst, expected[i])
		}
	}
}

func TestIsLessThan(t *testing.T) {
	expected := []bool{false, true, false, false, false, true, true, true, true}
	for i, v := range cases {
		if len(v.src) == 2 {
			continue
		}
		if v.src.IsLessThan(v.dst) != expected[i] {
			t.Errorf("Failed: src = %v / dst = %v / expected = %v", v.src, v.dst, expected[i])
		}
	}
}

func TestIsPessimisticConstraint(t *testing.T) {
	expected := []bool{true, true, false, true, false, false, true, false, false}
	for i, v := range cases {
		if v.src.IsPessimisticConstraint(v.dst) != expected[i] {
			t.Errorf("Failed: src = %v / dst = %v / expected = %v", v.src, v.dst, expected[i])
		}
	}
}
