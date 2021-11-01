// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"github.com/spf13/cobra"
)

// TanzuNsCtlCommand represents the base command when called without any subcommands.
type TanzuNsCtlCommand struct {
	*cobra.Command
}

// NewTanzuNsCtlCommand returns an instance of the TanzuNsCtlCommand.
func NewTanzuNsCtlCommand() *TanzuNsCtlCommand {
	c := &TanzuNsCtlCommand{
		Command: &cobra.Command{
			Use:   "TanzuNsCtl",
			Short: "Manage webstore stuff like a boss",
			Long:  "Manage webstore stuff like a boss",
		},
	}

	c.addSubCommands()

	return c
}

// Run represents the main entry point into the command
// This is called by main.main() to execute the root command.
func (c *TanzuNsCtlCommand) Run() {
	cobra.CheckErr(c.Execute())
}

// addSubCommands adds any additional subCommands to the root command.
func (c *TanzuNsCtlCommand) addSubCommands() {
	c.newInitCommand()
	c.newGenerateCommand()
	//+kubebuilder:scaffold:operator-builder:subcommands
}
