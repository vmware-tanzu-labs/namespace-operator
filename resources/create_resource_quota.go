// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT
/*

 */

package resources

import (
	"text/template"

	k8s_api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"

	tenancyv1alpha1 "github.com/vmware-tanzu-labs/namespace-operator/api/v1alpha1"
)

const resourceResourceQuota = `
---
apiVersion: v1
kind: ResourceQuota
metadata:
  name: tanzu-resource-quota
spec:
{{ $resourceQuotaCpuRequests := defaultResourceQuotaCPURequests .Spec }}
{{ $resourceQuotaMemoryRequests := defaultResourceQuotaMemoryRequests .Spec }}
{{ $resourceQuotaCpuLimits := defaultResourceQuotaCPULimits .Spec }}
{{ $resourceQuotaMemoryLimits := defaultResourceQuotaMemoryLimits .Spec }}
  hard:
    requests.cpu: {{ $resourceQuotaCpuRequests }}
    requests.memory: {{ $resourceQuotaMemoryRequests }}
    limits.cpu: {{ $resourceQuotaCpuLimits }}
    limits.memory: {{ $resourceQuotaMemoryLimits }}
`

// CreateResourceQuota creates the ResourceQuota resource for the parent TanzuNamespace object
func CreateResourceQuota(parent *tenancyv1alpha1.TanzuNamespace) (metav1.Object, error) {

	fmap := template.FuncMap{
		"defaultResourceQuotaCPURequests":    defaultResourceQuotaCPURequests,
		"defaultResourceQuotaMemoryRequests": defaultResourceQuotaMemoryRequests,
		"defaultResourceQuotaCPULimits":      defaultResourceQuotaCPULimits,
		"defaultResourceQuotaMemoryLimits":   defaultResourceQuotaMemoryLimits,
	}

	childContent, err := runTemplate("tanzu-resource-quota", resourceResourceQuota, parent, fmap)
	if err != nil {
		return nil, err
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(childContent), nil, nil)
	if err != nil {
		return nil, err
	}

	resourceObj := obj.(*k8s_api.ResourceQuota)
	resourceObj.Namespace = defaultNamespace(parent.Name, &parent.Spec)

	return resourceObj, nil
}
