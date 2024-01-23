package main

import (
	"github.com/fatih/color"
	"github.com/hashicorp/cli"
	trellisCmd "github.com/roots/trellis-cli/cmd"
	"github.com/roots/trellis-cli/trellis"
	"os"

	"github.com/ItinerisLtd/trellis-cli-kinsta/cmd"
)

var version = "canary"

func main() {
	c := cli.NewCLI("trellis-kinsta", version)
	c.Args = os.Args[1:]

	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColor{Code: int(color.FgYellow), Bold: false},
		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	trellis := trellis.NewTrellis()
	if err := trellis.LoadGlobalCliConfig(); err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	c.Commands = map[string]cli.CommandFactory{
		"site": func() (cli.Command, error) {
			return &trellisCmd.NamespaceCommand{
				HelpText:     "Usage: trellis kinsta site <subcommand> [<args>]",
				SynopsisText: "Manage your Kinsta sites",
			}, nil
		},
		"site list": func() (cli.Command, error) {
			return cmd.NewSiteListCommand(ui, trellis), nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}
