// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const defaultManifest = `apiVersion: tenancy.platform.cnr.vmware.com/v1alpha2
kind: TanzuNamespace
metadata:
  name: tanzunamespace-sample
spec:
  namespace: "tanzu-namespace"
  resources:
    limits:
      cpu: "250m"
      memory: "64Mi"
    requests:
      cpu: "250m"
      memory: "64Mi"
    max:
      cpu: "500m"
      memory: "500m"
    quota:
      requests:
        cpu: "2000m"
        memory: "4Gi"
      limits:
        cpu: "2000m"
        memory: "4Gi"
`

// newInitCommand creates a new instance of the init subcommand.
func (c *TanzuNsCtlCommand) newInitCommand() {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Write a sample custom resource manifest for a workload to standard out",
		Long:  "Write a sample custom resource manifest for a workload to standard out",
		RunE: func(cmd *cobra.Command, args []string) error {
			outputStream := os.Stdout

			if _, err := outputStream.WriteString(defaultManifest); err != nil {
				return fmt.Errorf("failed to write to stdout, %w", err)
			}

			return nil
		},
	}

	c.AddCommand(initCmd)
}
