package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
)

type CheckCommand struct {
	UI cli.Ui
}

func (c *CheckCommand) Run(args []string) int {
	currentDir, _ := os.Getwd()
	token := flag.String("token", "", "Terraform Cloud token")
	root := flag.String("root-path", currentDir, "Terraform config root path (default: current directory)")

	ws, err := InitCLI(*root, *token)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	currentVer, err := ws.GetCurrentVersion()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	latestVer, err := ws.GetLatestVersion()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if currentVer.String() != latestVer.String() {
		c.UI.Warn("New version is available.")
		c.UI.Info(fmt.Sprintf("%s -> %s", currentVer.String(), latestVer.String()))
	} else {
		c.UI.Info("No updates available.")
	}

	return 0
}

func (c *CheckCommand) Help() string {
	return strings.TrimSpace(helpMessageCheck)
}

func (c *CheckCommand) Synopsis() string {
	return "Check if new Terraform version is available"
}

const helpMessageCheck = `
Usage: terraform-cloud-updater check [OPTION]

--token        Terraform Cloud token        (default: TFE_TOKEN env var or parse from your .terraformrc)
--root-path    Terraform config root path   (default: current directory)
`
