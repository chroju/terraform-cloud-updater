package updater

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type tfRemoteBackend struct {
	Organization    string
	Workspace       string
	RequiredVersion string
	Hostname        string
}

type tfRc struct {
	Credentials []credential `hcl:"credentials,block"`
}

type credential struct {
	Name  string `hcl:"name,label"`
	Token string `hcl:"token"`
}

func parseTfRemoteBackend(root string) (*tfRemoteBackend, error) {
	var rb *tfRemoteBackend
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || !strings.HasSuffix(info.Name(), ".tf") {
				return nil
			}

			src, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			file, diags := hclwrite.ParseConfig(src, path, hcl.InitialPos)
			if diags.HasErrors() {
				return diags
			}

			for _, block := range file.Body().Blocks() {
				if block.Type() == "terraform" {
					for _, subBlock := range block.Body().Blocks() {
						if subBlock.Type() == "backend" && subBlock.Labels()[0] == "remote" {
							subBlockBody := subBlock.Body()
							rb = &tfRemoteBackend{
								Organization:    parseAttribute(subBlockBody.GetAttribute("organization")),
								Hostname:        parseAttribute(subBlockBody.GetAttribute("hostname")),
								Workspace:       parseAttribute(subBlockBody.Blocks()[0].Body().GetAttribute("name")),
								RequiredVersion: parseAttribute(block.Body().GetAttribute("required_version")),
							}
						}
					}
				}
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	if rb == nil {
		return nil, fmt.Errorf("Remote backend config is not found")
	}

	return rb, nil
}

func parseTerraformrc(path string) (string, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return "", fmt.Errorf("Parse %s failed", path)
	}
	var tfrc tfRc
	diags = gohcl.DecodeBody(f.Body, nil, &tfrc)
	if diags.HasErrors() {
		return "", fmt.Errorf("Decode %s failed", path)
	}
	return tfrc.Credentials[0].Token, nil
}

func parseAttribute(a *hclwrite.Attribute) string {
	if a == nil {
		return ""
	}
	return string(a.Expr().BuildTokens(nil)[1].Bytes)
}
