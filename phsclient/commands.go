package main

import (
	"git.aemian.com/phs/client/phsclient/command"
	"github.com/mitchellh/cli"
	"os"
)

// Commands is the mapping of all the available Consul commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}
	Commands = map[string]cli.CommandFactory{
		"test": func() (cli.Command, error) {
			return &command.TestCommand{
				Revision:          GitCommit,
				Version:           Version,
				VersionPrerelease: VersionPrerelease,
				Ui:                ui,
			}, nil
		},
		"request": func() (cli.Command, error) {
			return &command.RequestCommand{
				Revision:          GitCommit,
				Version:           Version,
				VersionPrerelease: VersionPrerelease,
				Ui:                ui,
			}, nil

		},
	}
}
