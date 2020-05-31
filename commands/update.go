package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/chroju/terraform-cloud-updater/updater"
	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
)

type UpdateCommand struct {
	UI cli.Ui
}

func (c *UpdateCommand) Run(args []string) int {
	var root, token string
	var updateVer *updater.SemanticVersion

	currentDir, _ := os.Getwd()
	f := flag.NewFlagSet("check", flag.ExitOnError)
	f.StringVar(&token, "token", "", "Terraform Cloud token")
	f.StringVar(&root, "root-path", currentDir, "Terraform config root path (default: current directory)")
	if err := f.Parse(args[1:]); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	ws, err := InitCLI(root, token)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if args[0] == "latest" {
		updateVer, err = ws.GetLatestVersion()
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
	} else {
		updateVer, err = updater.NewSemanticVersion(args[0])
		if err != nil {
			c.UI.Error(fmt.Sprintf("%s is not valid version", args[0]))
			c.UI.Output(helpMessageUpdate)
			return 1
		}
	}

	currentVer, err := ws.GetCurrentVersion()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if currentVer.String() == updateVer.String() {
		c.UI.Warn(fmt.Sprintf("Already latest version %s", updateVer.String()))
		return 0
	}

	if !ws.IsCompatibleVersion(updateVer) {
		c.UI.Error("This version is not compatible with required version.")
		if args[0] == "latest" {
			c.UI.Info(fmt.Sprintf("New version %s is available, but it is not compatible with required version %s", updateVer.String(), ws.GetRequiredVersions().String()))
		} else {
			c.UI.Error(fmt.Sprintf("Version %s is not compatible with required version %s", updateVer.String(), ws.GetRequiredVersions().String()))
		}
		c.UI.Info(fmt.Sprintf("\nLink to: %s", ws.GetSettingsLink()))
		return 3
	}

	if err = ws.UpdateVersion(updateVer); err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	c.UI.Info(fmt.Sprintf("Updated: %s -> %s", currentVer, updateVer))
	c.UI.Info(fmt.Sprintf("\nLink to: %s", ws.GetSettingsLink()))
	return 0
}

func (c *UpdateCommand) Help() string {
	return strings.TrimSpace(helpMessageUpdate)
}

func (c *UpdateCommand) Synopsis() string {
	return "Update Terraform cloud workspace terraform version"
}

const helpMessageUpdate = `
Usage: terraform-cloud-updater update <version> [OPTION]

Notes:
  version is must be in the correct semantic version format like 0.12.1, v0.12.2 .
  Or you can specify "latest" to automatically update to the latest version.

Options:
  --token        Terraform Cloud token        (default: TFE_TOKEN env var or parse from your .terraformrc)
  --root-path    Terraform config root path   (default: current directory)

`
