// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package tanzunamespace

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	tenancyv1alpha2 "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha2"
)

// CreateLimitRangeTanzuLimitRange creates the tanzu-limit-range LimitRange resource.
func CreateLimitRangeTanzuLimitRange(
	parent *tenancyv1alpha2.TanzuNamespace) (metav1.Object, error) {
	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "LimitRange",
			"metadata": map[string]interface{}{
				"name":      "tanzu-limit-range",
				"namespace": parent.Spec.Namespace,
			},
			"spec": map[string]interface{}{
				"limits": []interface{}{
					map[string]interface{}{
						"default": map[string]interface{}{
							// Default CPU limits to be applied to applications which get deployed into this namespace,
							// but are missing a resources declaration., controlled by resources.limits.cpu
							"cpu": parent.Spec.Resources.Limits.Cpu,
							// Default Memory limits to be applied to applications which get deployed into this namespace,
							// but are missing a resources declaration., controlled by resources.limits.memory
							"memory": parent.Spec.Resources.Limits.Memory,
						},
						"defaultRequest": map[string]interface{}{
							// Default CPU requests to be applied to applications which get deployed into this namespace,
							// but are missing a resources declaration., controlled by resources.requests.cpu
							"cpu": parent.Spec.Resources.Requests.Cpu,
							// Default Memory requests to be applied to applications which get deployed into this namespace,
							// but are missing a resources declaration., controlled by resources.requests.memory
							"memory": parent.Spec.Resources.Requests.Memory,
						},
						"max": map[string]interface{}{
							// Default maximum CPU limits for an individual application which get deployed into this namespace., controlled by resources.max.cpu
							"cpu": parent.Spec.Resources.Max.Cpu,
							// Default maximum Memory limits for an individual application which get deployed into this namespace., controlled by resources.max.memory
							"memory": parent.Spec.Resources.Max.Memory,
						},
						"type": "Container",
					},
				},
			},
		},
	}

	return resourceObj, nil
}
