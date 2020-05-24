package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chroju/terraform-cloud-updater/updater"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type CLIConfig struct {
	Token           string
	Hostname        string
	Organization    string
	Workspace       string
	RequiredVersion string
}

type tfRc struct {
	Credentials []credential `hcl:"credentials,block"`
}

type credential struct {
	Name  string `hcl:"name,label"`
	Token string `hcl:"token"`
}

// InitCLI initialize CLI config and creates a new workspace
func InitCLI(root, token string) (*updater.Workspace, error) {
	config, err := ParseTfFiles(root)
	if err != nil {
		return nil, err
	}

	if token != "" {
		config.Token = token
	}

	tfc, err := updater.NewTfCloud(config.Hostname, config.Token)
	if err != nil {
		return nil, err
	}

	ws, err := updater.NewWorkspace(tfc, &updater.Config{Organization: config.Organization, Workspace: config.Workspace, RequiredVersion: config.RequiredVersion})
	if err != nil {
		return nil, err
	}

	return ws, nil
}

// ParseTfFiles parses local Terraform files and creates a new CLIConfig
func ParseTfFiles(root string) (*CLIConfig, error) {
	config, err := parseTfRemoteBackend(root)
	if err != nil {
		return nil, err
	}

	token := os.Getenv("TFE_TOKEN")
	if token == "" {
		home := os.Getenv("HOME")
		if configPath := os.Getenv("TF_CLI_CONFIG_FILE"); configPath != "" {
			home = configPath
		}
		token, _ = parseTerraformrc(home + "/.terraformrc")
	}
	config.Token = token

	return config, nil
}

func parseTfRemoteBackend(root string) (*CLIConfig, error) {
	var config *CLIConfig
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
							config = &CLIConfig{
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

	if config == nil {
		return nil, fmt.Errorf("Remote backend config is not found")
	}

	return config, nil
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
