package updater

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	tfReleaseURL = "https://api.github.com/repos/hashicorp/terraform/releases"
)

type tfRelease struct {
	Draft           bool   `json:"draft"`
	Tag             string `json:"tag_name"`
	SemanticVersion SemanticVersion
}

func getTfReleases() ([]*tfRelease, error) {
	resp, err := http.Get(tfReleaseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tfReleases []*tfRelease
	if err = json.Unmarshal(body, &tfReleases); err != nil {
		return nil, err
	}
	for _, v := range tfReleases {
		v.Tag = v.Tag[1:]
		split := strings.Split(v.Tag, ".")
		var sv SemanticVersion
		sv = make([]int, len(split))
		for i, v2 := range split {
			sv[i], _ = strconv.Atoi(v2)
		}
		v.SemanticVersion = sv
	}
	return tfReleases, nil
}
