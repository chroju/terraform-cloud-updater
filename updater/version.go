package updater

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type semanticVersion []int

func (s *semanticVersion) String() string {
	var result string
	for _, v := range *s {
		result = strings.Join([]string{result, strconv.Itoa(v)}, ".")
	}
	return result
}

func (s *semanticVersion) IsEquall(target semanticVersion) bool {
	return reflect.DeepEqual(s, &target)
}

func (s *semanticVersion) IsNotEquall(target semanticVersion) bool {
	return !reflect.DeepEqual(s, &target)
}

func (s *semanticVersion) IsGreaterThan(target semanticVersion) bool {
	for i, v := range *s {
		if v > target[i] {
			return true
		}
		if v < target[i] {
			return false
		}
	}
	return false
}

func (s *semanticVersion) IsGreaterThanOrEqual(target semanticVersion) bool {
	return !s.IsLessThan(target)
}

func (s *semanticVersion) IsLessThan(target semanticVersion) bool {
	for i, v := range *s {
		if v > target[i] {
			return false
		}
		if v < target[i] {
			return true
		}
	}
	return false
}

func (s *semanticVersion) IsLessThanOrEqual(target semanticVersion) bool {
	return !s.IsGreaterThan(target)
}

func (s *semanticVersion) IsPessimisticConstraint(target semanticVersion) bool {
	for i, v := range *s {
		if v != target[i] {
			return false
		}
		if i == 1 {
			break
		}
	}
	return true
}

type RequiredVersion struct {
	Operator        string
	SemanticVersion semanticVersion
}

func (r *RequiredVersion) String() string {
	return fmt.Sprintf("%s %v", r.Operator, r.SemanticVersion)
}

type RequiredVersions []*RequiredVersion

func (r *RequiredVersions) CheckVersionConsistency(s semanticVersion) bool {
	for _, v := range *r {
		requiredVersion := v.SemanticVersion
		switch v.Operator {
		case "":
			if !s.IsEquall(requiredVersion) {
				return false
			}
		case "=":
			if !s.IsEquall(requiredVersion) {
				return false
			}
		case "!=":
			if !s.IsNotEquall(requiredVersion) {
				return false
			}
		case ">":
			if !s.IsGreaterThan(requiredVersion) {
				return false
			}
		case ">=":
			if !s.IsGreaterThanOrEqual(requiredVersion) {
				return false
			}
		case "<":
			if !s.IsLessThan(requiredVersion) {
				return false
			}
		case "<=":
			if !s.IsLessThanOrEqual(requiredVersion) {
				return false
			}
		case "~>":
			if !s.IsPessimisticConstraint(requiredVersion) {
				return false
			}
		}
	}
	return true
}
