// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package mutate

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

// TanzuNamespaceMutate performs the logic to mutate resources that belong to the parent.
func TanzuNamespaceMutate(reconciler common.ComponentReconciler,
	object *metav1.Object,
) (replacedObjects []metav1.Object, skip bool, err error) {
	return []metav1.Object{*object}, false, nil
}
