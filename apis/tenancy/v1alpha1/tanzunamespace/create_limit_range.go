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

	tenancyv1alpha1 "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha1"
)

const resourceLimitRange = `
---
apiVersion: v1
kind: LimitRange
metadata:
  name: tanzu-limit-range
spec:
{{ $limitRangeCpuLimit := defaultLimitRangeDefaultCPULimit .Spec }}
{{ $limitRangeMemoryLimit := defaultLimitRangeDefaultMemoryLimit .Spec }}
{{ $limitRangeCpuRequest := defaultLimitRangeDefaultCPURequest .Spec }}
{{ $limitRangeMemoryRequest := defaultLimitRangeDefaultMemoryRequest .Spec }}
{{ $limitRangeCpuLimitMax := defaultLimitRangeMaxCPULimit .Spec }}
{{ $limitRangeMemoryLimitMax := defaultLimitRangeMaxMemoryLimit .Spec }}
  limits:
    - default:
        cpu: {{ $limitRangeCpuLimit }}
        memory: {{ $limitRangeMemoryLimit }}
      defaultRequest:
        cpu: {{ $limitRangeCpuRequest }}
        memory: {{ $limitRangeMemoryRequest }}
      max:
        cpu: {{ $limitRangeCpuLimitMax }}
        memory: {{ $limitRangeMemoryLimitMax }}
      type: Container
`

// CreateLimitRange creates the LimitRange resource for the parent TanzuNamespace object
func CreateLimitRange(parent *tenancyv1alpha1.TanzuNamespace) (metav1.Object, error) {

	fmap := template.FuncMap{
		"defaultLimitRangeDefaultCPULimit":      defaultLimitRangeDefaultCPULimit,
		"defaultLimitRangeDefaultMemoryLimit":   defaultLimitRangeDefaultMemoryLimit,
		"defaultLimitRangeDefaultCPURequest":    defaultLimitRangeDefaultCPURequest,
		"defaultLimitRangeDefaultMemoryRequest": defaultLimitRangeDefaultMemoryRequest,
		"defaultLimitRangeMaxCPULimit":          defaultLimitRangeMaxCPULimit,
		"defaultLimitRangeMaxMemoryLimit":       defaultLimitRangeMaxMemoryLimit,
	}

	childContent, err := runTemplate("tanzu-limit-range", resourceLimitRange, parent, fmap)
	if err != nil {
		return nil, err
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(childContent), nil, nil)
	if err != nil {
		return nil, err
	}

	resourceObj := obj.(*k8s_api.LimitRange)
	resourceObj.Namespace = defaultNamespace(parent.Name, &parent.Spec)

	return resourceObj, nil
}
