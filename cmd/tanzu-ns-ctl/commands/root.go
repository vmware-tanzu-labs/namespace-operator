// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT
/*

 */

package commands

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "tanzu-ns-ctl <subcommand>",
		Short: "workload manager",
	}

	cmd.AddCommand(NewInitCommand())
	cmd.AddCommand(NewGenerateCommand())

	return &cmd
}
