// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT
/*

 */
package main

import (
	"os"

	"github.com/vmware-tanzu-labs/namespace-operator/cmd/tanzu-ns-ctl/commands"
)

func main() {
	if err := commands.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
