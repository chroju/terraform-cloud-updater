package updater

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	tfReleaseURL = "https://api.github.com/repos/hashicorp/terraform/releases"
)

type tfVersion struct {
	Draft           bool   `json:"draft"`
	Version         string `json:"tag_name"`
	SemanticVersion SemanticVersion
}

type SemanticVersion []int

func (s *SemanticVersion) String() string {
	var result string
	for _, v := range *s {
		result = strings.Join([]string{result, strconv.Itoa(v)}, ".")
	}
	return result
}

func getTerraformVersions() ([]*tfVersion, error) {
	resp, err := http.Get(tfReleaseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tfVersions []*tfVersion
	if err = json.Unmarshal(body, &tfVersions); err != nil {
		return nil, err
	}
	for _, v := range tfVersions {
		v.Version = v.Version[1:]
		split := strings.Split(v.Version, ".")
		var sv SemanticVersion
		sv = make([]int, len(split))
		for i, v2 := range split {
			sv[i], _ = strconv.Atoi(v2)
		}
		v.SemanticVersion = sv
	}
	return tfVersions, nil
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
