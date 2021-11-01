// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package helpers

import (
	"fmt"

	common "github.com/vmware-tanzu-labs/namespace-operator/apis/common"
	tenancyv1alpha2 "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha2"
)

// TanzuNamespaceUnique returns only one TanzuNamespace and returns an error if more than one are found.
func TanzuNamespaceUnique(
	reconciler common.ComponentReconciler,
) (
	*tenancyv1alpha2.TanzuNamespace,
	error,
) {
	components, err := TanzuNamespaceList(reconciler)
	if err != nil {
		return nil, err
	}

	if len(components.Items) != 1 {
		return nil, fmt.Errorf("expected only 1 TanzuNamespace; found %v\n", len(components.Items))
	}

	component := components.Items[0]

	return &component, nil
}

// TanzuNamespaceList gets a TanzuNamespaceList from the cluster.
func TanzuNamespaceList(
	reconciler common.ComponentReconciler,
) (
	*tenancyv1alpha2.TanzuNamespaceList,
	error,
) {
	components := &tenancyv1alpha2.TanzuNamespaceList{}
	if err := reconciler.List(reconciler.GetContext(), components); err != nil {
		reconciler.GetLogger().V(0).Info("unable to retrieve TanzuNamespaceList from cluster")

		return nil, err
	}

	return components, nil
}
