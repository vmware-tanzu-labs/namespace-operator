// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package wait

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

// TanzuNamespaceWait performs the logic to wait for resources that belong to the parent.
func TanzuNamespaceWait(reconciler common.ComponentReconciler,
	object *metav1.Object,
) (ready bool, err error) {
	return true, nil
}
