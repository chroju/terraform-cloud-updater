package updater

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// SemanticVersion represents semantic version
type SemanticVersion struct {
	Versions []int
}

func (s *SemanticVersion) String() string {
	result := make([]string, len(s.Versions))
	for i, v := range s.Versions {
		result[i] = strconv.Itoa(v)
	}
	return strings.Join(result, ".")
}

// NewSemanticVersion creates a new SemanticVersion from the string represents semantic version
func NewSemanticVersion(versionString string) (*SemanticVersion, error) {
	split := strings.Split(versionString, ".")
	sv := make([]int, len(split))
	for i, v := range split {
		converted, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		sv[i] = converted
	}
	return &SemanticVersion{Versions: sv}, nil
}

// RequiredVersion represents Terraform required version
type RequiredVersion struct {
	Operator        Operator
	SemanticVersion *SemanticVersion
}

// Operator is required version operator
type Operator string

const (
	blank                 Operator = ""
	equal                 Operator = "="
	notEqual              Operator = "!="
	greaterThan           Operator = ">"
	greaterThanOrEqual    Operator = ">="
	lessThan              Operator = "<"
	lessThanEqual         Operator = "<="
	pessimisticConstraint Operator = "~>"
	operators             string   = "=!<>~"
)

func (r *RequiredVersion) String() string {
	return fmt.Sprintf("%s %v", r.Operator, r.SemanticVersion)
}

// RequiredVersions represents Terraform required versions
type RequiredVersions []*RequiredVersion

// NewRequiredVersions returns new RequiredVersions from given constraints
func NewRequiredVersions(versionString string) (RequiredVersions, error) {
	if strings.Contains(versionString, ",") {
		split := strings.Split(versionString, ",")
		rvs := make([]*RequiredVersion, len(split))
		for i, v := range split {
			rv, err := NewRequiredVersions(v)
			if err != nil {
				return nil, err
			}
			rvs[i] = rv[0]
		}
		return rvs, nil
	}

	rv := &RequiredVersion{}
	if strings.ContainsAny(versionString, operators) {
		index := strings.LastIndexAny(versionString, operators)
		rv.Operator = Operator(strings.TrimSpace(versionString[:index+1]))
		versionString = strings.TrimSpace(versionString[index+1:])
	} else {
		rv.Operator = blank
	}

	sv, err := NewSemanticVersion(versionString)
	if err != nil {
		return nil, err
	}
	rv.SemanticVersion = sv

	return []*RequiredVersion{rv}, nil
}

func (r *RequiredVersions) String() string {
	result := make([]string, len(*r))
	for i, v := range *r {
		result[i] = v.String()
	}
	return strings.Join(result, ", ")
}

// CheckVersionConsistency checks given semantic version is consistent with requreid versions
func (r *RequiredVersions) CheckVersionConsistency(s *SemanticVersion) bool {
	for _, v := range *r {
		switch v.Operator {
		case blank, equal:
			if !v.IsEquall(s) {
				return false
			}
		case notEqual:
			if !v.IsNotEquall(s) {
				return false
			}
		case greaterThan:
			if !v.IsGreaterThan(s) {
				return false
			}
		case greaterThanOrEqual:
			if !v.IsGreaterThanOrEqual(s) {
				return false
			}
		case lessThan:
			if !v.IsLessThan(s) {
				return false
			}
		case lessThanEqual:
			if !v.IsLessThanOrEqual(s) {
				return false
			}
		case pessimisticConstraint:
			if !v.IsPessimisticConstraint(s) {
				return false
			}
		}
	}
	return true
}

func (r *RequiredVersion) IsEquall(target *SemanticVersion) bool {
	return reflect.DeepEqual(r.SemanticVersion, target)
}

func (r *RequiredVersion) IsNotEquall(target *SemanticVersion) bool {
	return !reflect.DeepEqual(r.SemanticVersion, target)
}

func (r *RequiredVersion) IsGreaterThan(target *SemanticVersion) bool {
	for i, v := range r.SemanticVersion.Versions {
		if target.Versions[i] > v {
			return true
		}
		if target.Versions[i] < v {
			return false
		}
	}
	if len(r.SemanticVersion.Versions) == 2 {
		return true
	}
	return false
}

func (r *RequiredVersion) IsGreaterThanOrEqual(target *SemanticVersion) bool {
	return !r.IsLessThan(target)
}

func (r *RequiredVersion) IsLessThan(target *SemanticVersion) bool {
	for i, v := range r.SemanticVersion.Versions {
		if target.Versions[i] < v {
			return true
		}
		if target.Versions[i] > v {
			return false
		}
	}
	return false
}

func (r *RequiredVersion) IsLessThanOrEqual(target *SemanticVersion) bool {
	return !r.IsGreaterThan(target)
}

func (r *RequiredVersion) IsPessimisticConstraint(target *SemanticVersion) bool {
	// `~> 0.9` is equivalent to `>= 0.9` and `< 1.0`
	// First, check `>=`
	if r.IsLessThan(target) {
		return false
	}
	// Second, check `<`
	nextVersion := make([]int, 2)
	if r.SemanticVersion.Versions[1] == 9 {
		nextVersion[1] = 0
	} else {
		nextVersion[1] = r.SemanticVersion.Versions[1] + 1
	}
	r.SemanticVersion.Versions = nextVersion

	if r.IsGreaterThanOrEqual(target) {
		return false
	}

	return true
}
