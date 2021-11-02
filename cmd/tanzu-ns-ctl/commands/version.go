// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"
var apiVersions = []string{
	"v1alpha1",
	"v1alpha2",
	//+kubebuilder:scaffold:operator-builder:apiversions
}

// newVersionCommand creates a new instance of the version subcommand.
func (c *TanzuNsCtlCommand) newVersionCommand() {
	initCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  "Print the version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			versionInfo := struct {
				CLIVersion  string   `json:"cliVersion"`
				APIVersions []string `json:"apiVersions"`
			}{
				CLIVersion:  version,
				APIVersions: apiVersions,
			}

			output, err := json.Marshal(versionInfo)
			if err != nil {
				return fmt.Errorf("failed to determine versionInfo, %w", err)
			}

			outputStream := os.Stdout

			if _, err := outputStream.WriteString(fmt.Sprintln(string(output))); err != nil {
				return fmt.Errorf("failed to write to stdout, %w", err)
			}

			return nil
		},
	}

	c.AddCommand(initCmd)
}
