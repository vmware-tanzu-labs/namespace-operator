// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package tanzunamespace

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	tenancyv1alpha2 "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha2"
)

// CreateResourceQuotaTanzuResourceQuota creates the tanzu-resource-quota ResourceQuota resource.
func CreateResourceQuotaTanzuResourceQuota(
	parent *tenancyv1alpha2.TanzuNamespace) (metav1.Object, error) {
	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ResourceQuota",
			"metadata": map[string]interface{}{
				"name":      "tanzu-resource-quota",
				"namespace": parent.Spec.Namespace,
			},
			"spec": map[string]interface{}{
				"hard": map[string]interface{}{
					// Default CPU requests quota to be enforced on the sum of all applications which get deployed into this namespace., controlled by resources.quota.requests.cpu
					"requests.cpu": parent.Spec.Resources.Quota.Requests.Cpu,
					// Default Memory requests quota to be enforced on the sum of all applications which get deployed into this namespace., controlled by resources.quota.requests.memory
					"requests.memory": parent.Spec.Resources.Quota.Requests.Memory,
					// Default CPU limits quota to be enforced on the sum of all applications which get deployed into this namespace., controlled by resources.quota.limits.cpu
					"limits.cpu": parent.Spec.Resources.Quota.Limits.Cpu,
					// Default Memory limits quota to be enforced on the sum of all applications which get deployed into this namespace., controlled by resources.quota.limits.memory
					"limits.memory": parent.Spec.Resources.Quota.Limits.Memory,
				},
			},
		},
	}

	return resourceObj, nil
}
