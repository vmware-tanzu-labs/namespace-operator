// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

// validateWorkload validates the unmarshaled version of the workload resource
// manifest.
func validateWorkload(
	workload common.Component,
) error {
	defaultWorkloadGVK := workload.GetComponentGVK()
	component := workload.(runtime.Object)

	if defaultWorkloadGVK != component.GetObjectKind().GroupVersionKind() {
		return fmt.Errorf(
			"expected resource of kind: '%s', with group '%s' and version '%s'; "+
				"found resource of kind '%s', with group '%s' and version '%s'",
			defaultWorkloadGVK.Kind,
			defaultWorkloadGVK.Group,
			defaultWorkloadGVK.Version,
			component.GetObjectKind().GroupVersionKind().Kind,
			component.GetObjectKind().GroupVersionKind().Group,
			component.GetObjectKind().GroupVersionKind().Version,
		)
	}

	return nil
}
