// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package tanzunamespace

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	tenancyv1alpha2 "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha2"
)

// CreateNamespaceParentSpecNamespace creates the parent.Spec.Namespace Namespace resource.
func CreateNamespaceParentSpecNamespace(
	parent *tenancyv1alpha2.TanzuNamespace) (metav1.Object, error) {
	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Namespace",
			"metadata": map[string]interface{}{
				// Namespace name which is created and then enforced by related policy objects
				// such as LimitRange, ResourceQuota, and NetworkPolicy., controlled by namespace
				"name": parent.Spec.Namespace,
			},
		},
	}

	return resourceObj, nil
}
