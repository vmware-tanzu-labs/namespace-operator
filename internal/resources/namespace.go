// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	v1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

const (
	NamespaceKind = "Namespace"
)

// NamespaceIsReady defines the criteria for a namespace to be condsidered ready.
func NamespaceIsReady(resource common.ComponentResource) (bool, error) {
	var namespace v1.Namespace
	if err := getObject(resource, &namespace, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if namespace.Name == "" {
		return false, nil
	}

	// if the namespace is terminating, it is not considered ready
	if namespace.Status.Phase == v1.NamespaceTerminating {
		return false, nil
	}

	// finally, rely on the active field to determine if this namespace is ready
	return namespace.Status.Phase == v1.NamespaceActive, nil
}

// NamespaceForResourceIsReady checks to see if the namespace of a resource is ready.
func NamespaceForResourceIsReady(resource common.ComponentResource) (bool, error) {
	// create a stub namespace resource to pass to the NamespaceIsReady method
	namespace := &Resource{
		Reconciler: resource.GetReconciler(),
	}

	// insert the inherited fields
	namespace.Name = resource.GetNamespace()
	namespace.Group = ""
	namespace.Version = "v1"
	namespace.Kind = NamespaceKind

	return NamespaceIsReady(namespace)
}
