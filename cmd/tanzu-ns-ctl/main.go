// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/vmware-tanzu-labs/namespace-operator/cmd/tanzu-ns-ctl/commands"
)

func main() {
	tanzunsctl := commands.NewTanzuNsCtlCommand()
	tanzunsctl.Run()
}
