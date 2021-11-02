// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package tanzunamespace

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	tenancyv1alpha2 "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha2"
)

// CreateNetworkPolicyTanzuNetworkPolicy creates the tanzu-network-policy NetworkPolicy resource.
func CreateNetworkPolicyTanzuNetworkPolicy(
	parent *tenancyv1alpha2.TanzuNamespace) (metav1.Object, error) {
	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			// NOTE: code markers were added before functionality for complex data types existed.  these are not functional
			//       or even syntactically correct but for documentation/visualization purposes only.
			"apiVersion": "networking.k8s.io/v1",
			"kind":       "NetworkPolicy",
			"metadata": map[string]interface{}{
				"name":      "tanzu-network-policy",
				"namespace": parent.Spec.Namespace,
			},
			"spec": map[string]interface{}{
				"podSelector": map[string]interface{}{},
				"policyTypes": []interface{}{
					"Ingress",
					"Egress",
				},
				"egress": []interface{}{
					map[string]interface{}{
						"ports": []interface{}{
							map[string]interface{}{
								"protocol": "UDP",
								"port":     53,
							},
						},
					},
				},
			},
		},
	}

	return resourceObj, nil
}
