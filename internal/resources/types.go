// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

// Resource represents a resource as managed during the reconciliation process.
type Resource struct {
	common.ResourceCommon

	Object     client.Object
	Reconciler common.ComponentReconciler
}
