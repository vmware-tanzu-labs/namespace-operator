// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0
/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"errors"
	"io"
	"log"

	"github.com/spf13/cobra"
)

const defaultManifest = `apiVersion: tenancy.platform.cnr.vmware.com/v1alpha1
kind: TanzuNamespace
metadata:
  name: tanzunamespace-sample
spec:
  tanzuNamespaceName: "tanzu-namespace"
  tanzuLimitRangeDefaultCpuLimit: "250m"
  tanzuLimitRangeDefaultMemoryLimit: "64Mi"
  tanzuLimitRangeDefaultCpuRequest: "250m"
  tanzuLimitRangeDefaultMemoryRequest: "64Mi"
  tanzuLimitRangeMaxCpuLimit: "1000m"
  tanzuLimitRangeMaxMemoryLimit: "512Mi"
  tanzuResourceQuotaCpuRequests: "2"
  tanzuResourceQuotaMemoryRequests: "4Gi"
  tanzuResourceQuotaCpuLimits: "2"
  tanzuResourceQuotaMemoryLimits: "4Gi"
  networkPolicies:
    - targetPodLabels: {}
      ingressTCPPorts: []
      ingressUDPPorts: []
      ingressPodLabels: {}
      ingressNamespaceLabels: {}
      egressTCPPorts: []
      egressUDPPorts: []
      egressPodLabels: {}
      egressNamespaceLabels: {}
`

func NewInitCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "init <file>",
		Short: "Initialize a new manifest file",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a file argument")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			out, err := openOutput(args[0])
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := out.Close(); err != nil {
					log.Fatal(err)
				}
			}()

			if err := runInit(out); err != nil {
				log.Fatal(err)
			}
		},
	}

	return &cmd
}

func runInit(w io.Writer) error {
	if _, err := w.Write([]byte(defaultManifest)); err != nil {
		return err
	}
	return nil
}
