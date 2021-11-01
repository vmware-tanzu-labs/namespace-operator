// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package dependencies

import (
	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

// TanzuNamespaceCheckReady performs the logic to determine if a TanzuNamespace object is ready.
func TanzuNamespaceCheckReady(reconciler common.ComponentReconciler) (bool, error) {
	return true, nil
}
