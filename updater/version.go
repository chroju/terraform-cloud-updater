package updater

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type SemanticVersion []int

func (s *SemanticVersion) String() string {
	var result string
	for _, v := range *s {
		result = strings.Join([]string{result, strconv.Itoa(v)}, ".")
	}
	return result
}

func (s *SemanticVersion) IsEquall(target SemanticVersion) bool {
	return reflect.DeepEqual(s, &target)
}

func (s *SemanticVersion) IsNotEquall(target SemanticVersion) bool {
	return !reflect.DeepEqual(s, &target)
}

func (s *SemanticVersion) IsGreaterThan(target SemanticVersion) bool {
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

func (s *SemanticVersion) IsGreaterThanOrEqual(target SemanticVersion) bool {
	return !s.IsLessThan(target)
}

func (s *SemanticVersion) IsLessThan(target SemanticVersion) bool {
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

func (s *SemanticVersion) IsLessThanOrEqual(target SemanticVersion) bool {
	return !s.IsGreaterThan(target)
}

func (s *SemanticVersion) IsPessimisticConstraint(target SemanticVersion) bool {
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
	SemanticVersion SemanticVersion
}

func (r *RequiredVersion) String() string {
	return fmt.Sprintf("%s %v", r.Operator, r.SemanticVersion)
}

type RequiredVersions []*RequiredVersion

func (r *RequiredVersions) CheckVersionConsistency(s SemanticVersion) bool {
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

func NewRequiredVersions(versionString string) RequiredVersions {
	if strings.Contains(versionString, ",") {
		split := strings.Split(versionString, ",")
		rvs := make([]*RequiredVersion, len(split))
		for i, v := range split {
			rvs[i] = NewRequiredVersions(v)[0]
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
