package main

import (
	"fmt"
	"os"

	"github.com/chroju/terraform-cloud-updater/commands"
	"github.com/mitchellh/cli"
)

const (
	app     = "terraform-cloud-updater"
	version = "0.1.0"
)

func main() {
	c := cli.NewCLI(app, version)
	c.Args = os.Args[1:]
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c.Commands = map[string]cli.CommandFactory{
		"check": func() (cli.Command, error) {
			return &commands.CheckCommand{UI: &cli.ColoredUi{Ui: ui, WarnColor: cli.UiColorYellow, ErrorColor: cli.UiColorRed}}, nil
		},
		"update": func() (cli.Command, error) {
			return &commands.UpdateCommand{UI: &cli.ColoredUi{Ui: ui, ErrorColor: cli.UiColorRed}}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(fmt.Sprintf("Error: %s", err))
	}

	os.Exit(exitStatus)
}
